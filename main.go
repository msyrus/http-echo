package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	var host, port string
	flag.StringVar(&host, "host", "0.0.0.0", "host name")
	flag.StringVar(&port, "port", "8000", "port number")
	flag.Parse()

	addr := host + ":" + port
	log.Println("Starting echo server on", addr)
	if err := http.ListenAndServe(addr, http.HandlerFunc(EchoHandler)); err != nil {
		log.Println("Server stopped")
		log.Fatal(err)
	}
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	for k, vs := range r.Header {
		for _, v := range vs {
			w.Header().Add("Req-"+k, v)
		}
	}
	host, _ := os.Hostname()
	w.Header().Set("Host", host)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	if r.Body != nil {
		defer r.Body.Close()
		if _, err := io.Copy(w, r.Body); err != nil {
			log.Println(err)
		}
	}
}
