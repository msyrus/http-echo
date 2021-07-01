package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var printBody, printHeaders bool

func main() {
	var host, port string
	flag.StringVar(&host, "host", "0.0.0.0", "host name")
	flag.StringVar(&port, "port", "8000", "port number")
	flag.BoolVar(&printBody, "print-body", false, "print request body to stdout")
	flag.BoolVar(&printHeaders, "print-headers", false, "print request headers to stdout")
	flag.Parse()

	addr := host + ":" + port
	log.Println("Starting echo server on", addr)
	if err := http.ListenAndServe(addr, Logger(http.HandlerFunc(EchoHandler))); err != nil {
		log.Println("Server stopped")
		log.Fatal(err)
	}
}

// Logger is the http request logger middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Proto, r.RemoteAddr, r.Host, r.Method, r.RequestURI)
		if printHeaders {
			for k, vs := range r.Header {
				fmt.Fprintf(os.Stdout, "%s: %s\n", k, strings.Join(vs, ", "))
			}
		}
		if printBody && r.Body != nil && r.ContentLength != 0 {
			defer func() {
				r.Body.Close()
				os.Stdout.WriteString("\n")
			}()
			r.Body = io.NopCloser(io.TeeReader(r.Body, os.Stdout))
		}
		next.ServeHTTP(w, r)
	})
}

// EchoHandler echo request body to response
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	for k, vs := range r.Header {
		for _, v := range vs {
			w.Header().Add("Req-"+k, v)
		}
	}

	w.Header().Add("Req-Host", r.Host)
	w.Header().Set("Req-Remote", r.RemoteAddr)
	host, _ := os.Hostname()
	w.Header().Set("Host", host)
	if cType := r.Header.Get("Content-Type"); cType != "" {
		w.Header().Set("Content-Type", cType)
	}

	if r.Body != nil {
		defer r.Body.Close()
		if _, err := io.Copy(w, r.Body); err != nil {
			log.Println(err)
		}
	}
}
