package main

import (
	"fmt"
	"net/http"
)

func main() {
	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: http.HandlerFunc(endpoint),
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func endpoint(http.ResponseWriter, *http.Request) {
	fmt.Println("hit")
}
