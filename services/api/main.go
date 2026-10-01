// api: servicio HTTP mínimo de ejemplo para el pipeline de la plataforma.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	version := os.Getenv("APP_VERSION")

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host, _ := os.Hostname()
		fmt.Fprintf(w, "hola desde qs-sample-app api, desplegado solo con un git push (pod %s, versión %q)\n", host, version)
	})

	log.Printf("api escuchando en :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
