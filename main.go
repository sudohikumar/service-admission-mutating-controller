package main

import (
	"admission-controller/handlers"
	"admission-controller/router"
	"path/filepath"
)

const (
	tlsDir      = `/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

func main() {
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	ginRouter := router.CreateRouter()
	ginRouter.AddHandler(handlers.AdmissionHandler)
	ginRouter.Run(certPath, keyPath)
}