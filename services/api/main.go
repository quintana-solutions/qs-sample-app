// api: servicio HTTP mínimo de ejemplo para el pipeline de la plataforma.
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

	// /whoami: en una ruta auth: sso, Envoy ya validó el login y pone el ID token de Zitadel en x-id-token. La app solo
	// lee sus claims (no verifica la firma: confía en el Gateway, el único camino para llegar a ella).
	http.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("x-id-token")
		parts := strings.Split(tok, ".")
		if len(parts) != 3 {
			fmt.Fprintln(w, "anónimo (esta ruta no pide login)")
			return
		}
		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		var claims struct {
			Sub   string `json:"sub"`
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		if err != nil || json.Unmarshal(payload, &claims) != nil {
			http.Error(w, "x-id-token ilegible", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "hola %s <%s> (sub %s)\n", claims.Name, claims.Email, claims.Sub)
	})

	log.Printf("api escuchando en :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
