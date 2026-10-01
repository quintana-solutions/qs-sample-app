// worker: proceso en segundo plano de ejemplo (sin puerto); escribe un log cada intervalo.
package main

import (
	"log"
	"os"
	"time"
)

func main() {
	every, err := time.ParseDuration(os.Getenv("INTERVAL"))
	if err != nil || every <= 0 {
		every = 30 * time.Second
	}
	for i := 1; ; i++ {
		log.Printf("worker: tarea %d procesada", i)
		time.Sleep(every)
	}
}
