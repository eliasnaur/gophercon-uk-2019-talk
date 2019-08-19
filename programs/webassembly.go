package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	err := http.ListenAndServe(":8080", http.FileServer(http.Dir(os.Args[1])))
	log.Fatal(err)
}
