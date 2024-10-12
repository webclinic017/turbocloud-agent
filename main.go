package main

import (
	"fmt"
	"log"
	"net/http"
)

func use(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = r

	for _, mw := range middlewares {
		s = mw(s)
	}

	return s
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Before %s", r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func acceptHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Accept: %v", r.Header.Get("Accept"))
		next.ServeHTTP(w, r)
	})
}

func main() {

	databaseInit()

	mux := http.NewServeMux()
	//mux.HandleFunc("POST /proxy/{id}", handleProxyPost)
	mux.HandleFunc("POST /proxy", handleProxyPost)
	mux.HandleFunc("GET /proxy", handleProxyGet)
	mux.HandleFunc("DELETE /proxy/{id}", handleProxyDelete)

	wrapped := use(mux, loggingMiddleware, acceptHeaderMiddleware)

	//http.HandleFunc("/hey", handler)
	//create_service()

	PORT := "5445"
	fmt.Println("Starting an agent on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, wrapped))
}
