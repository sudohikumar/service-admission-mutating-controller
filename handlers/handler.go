package handlers

import (
	"admission-controller/helpers"
	"admission-controller/router"
	"admission-controller/structures"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"log"
	"net/http"
	"strings"
)

var serviceResource = metav1.GroupVersionResource{
	Resource: "services",
	Version:  "v1",
}

var universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()

func AdmissionHandler(r router.GinRouter) {
	r.Router.POST("/admission", func(c *gin.Context) {
		b, _ := ioutil.ReadAll(c.Request.Body)
		admissionReview, err := admissionHelper(b)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			r, _ := json.Marshal(admissionReview)
			c.Writer.Write(r)
		}
	})
}

func admissionHelper(data []byte) (v1.AdmissionReview, error) {
	var admissionReviewReq v1.AdmissionReview
	_, _, err := universalDeserializer.Decode(data, nil, &admissionReviewReq)
	if err != nil {
		log.Println("Error serializer: ", err.Error())
		return v1.AdmissionReview{}, err
	}
	admissionReviewResponse := v1.AdmissionReview{
		Response: &v1.AdmissionResponse{
			UID: admissionReviewReq.Request.UID,
		},
	}
	err = performOperation(admissionReviewReq.Request)
	if err != nil {
		log.Println("Error Operation: ", err.Error())
		admissionReviewResponse.Response.Allowed = false
		admissionReviewResponse.Response.Result = &metav1.Status{
			Message: err.Error(),
		}
	} else {
		admissionReviewResponse.Response.Allowed = true
		admissionReviewResponse.Response.Patch = createPatch()
		pt := v1.PatchTypeJSONPatch
		admissionReviewResponse.Response.PatchType = &pt
	}
	admissionReviewReq.Response = admissionReviewResponse.Response
	return admissionReviewReq, nil
}

func createPatch() []byte {
	var p []structures.Patch
	patch := structures.Patch{
		Op:    "add",
		Path:  "/metadata/labels",
		Value: map[string]string{
			"mutated-via-controller": "true",
		},
	}
	p = append(p, patch)
	b, _ := json.Marshal(p)
	return b
}

func performOperation(req *v1.AdmissionRequest) error {
	if req.Resource != serviceResource {
		log.Print("resource is not of service type")
		return nil
	}
	// If kube-namespace, return from here and do nothing
	if helpers.IsKubeNamespace(req.Namespace) {
		return nil
	}
	raw := req.Object.Raw
	service := corev1.Service{}
	_, _, err := universalDeserializer.Decode(raw, nil, &service)
	if err != nil {
		log.Print(err)
		return nil
	}
	if strings.Contains(service.Name, "simple") {
		return fmt.Errorf("service name should not contain 'simple' word")
	}
	return nil
}
