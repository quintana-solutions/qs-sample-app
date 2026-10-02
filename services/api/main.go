// api: servicio HTTP mínimo de ejemplo para el pipeline de la plataforma.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
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

	// /db: consulta la base propia de la app (database.postgres en quintana.yaml). La plataforma inyecta DATABASE_URL.
	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		url := os.Getenv("DATABASE_URL")
		if url == "" {
			http.Error(w, "sin base de datos (database.postgres no está en quintana.yaml)", http.StatusServiceUnavailable)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		conn, err := pgx.Connect(ctx, url)
		if err != nil {
			http.Error(w, "no se pudo conectar a la base: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer conn.Close(ctx)
		var db, version string
		if err := conn.QueryRow(ctx, "select current_database(), version()").Scan(&db, &version); err != nil {
			http.Error(w, "consulta fallida: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "base %q: %s\n", db, version)
	})

	log.Printf("api escuchando en :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
