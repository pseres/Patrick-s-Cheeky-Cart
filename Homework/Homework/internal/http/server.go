package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	MatchingProductIds []uint32 `json:"matchingProductIds"`
}

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
	mux.HandleFunc("/", handler)

	return mux
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")
	if search == "" {
		status := http.StatusBadRequest
		msg := "Missing 'search' query parameter"
		http.Error(w, msg, status)
		return
	}

	if search == "foo" {
		status := http.StatusNotFound
		msg := fmt.Sprintf("No products found matching the search criteria: '%s'", search)
		http.Error(w, msg, status)
		return
	}

	resp := Response{
		MatchingProductIds: []uint32{1, 2, 3},
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
}
