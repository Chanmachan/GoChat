package main

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func main() {
	fmt.Print("Hello")
	http.HandleFunc("/", Hello)
	http.ListenAndServe(":8080", nil)
}
