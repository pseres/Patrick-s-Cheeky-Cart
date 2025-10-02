package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
)

var server *http.Server

func StartServer(port int) {
	server = &http.Server{Addr: ":" + strconv.Itoa(port)}
	server.Handler = buildMux()

	go func() {
		log.Printf("HTTP server started at localhost%v\n", server.Addr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to start HTTP server on port %v [%v]\n", server.Addr, err)
			os.Exit(-1)
		}

		log.Println("HTTP server stopped")
	}()
}

func StopServer() {
	if server != nil {
		server.Shutdown(context.TODO())
	}
}

func buildMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received GET /products")
		w.Write([]byte("List of products"))
	})

	return mux
}
