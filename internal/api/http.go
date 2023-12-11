package api

import (
	"context"
	"fmt"
	"net/http"
)

func StartHttpServer(cancel context.CancelFunc) {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Invoker is running...")
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
	if err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
		cancel()
	}
}
