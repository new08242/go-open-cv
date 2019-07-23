package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/gimme_faces", GetFaceFromPictureHandler)

	http.ListenAndServe(":8080", mux)
}