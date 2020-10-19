package main

import (
	"log"
	"net/http"
)
func main() {
	srv := &http.Server{Addr: ":10000", Handler: http.HandlerFunc(handle)}

	// Start the server with TLS, since we are running
	// HTTP/2 it must be run with TLS.
	// Exactly how you would run an HTTP/1.1 server with TLS connection.
	log.Printf("Serving on https://0.0.0.0:10000")
	log.Fatal(srv.ListenAndServeTLS("/Users/ijaehyeon/keys/server.crt",
		"/Users/ijaehyeon/keys/server.pem"))

	//log.Fatal(srv.ListenAndServe())
}

func default_handle(w http.ResponseWriter, r *http.Request) {
	// Log the request protocol
	log.Printf("Got connection: %s", r.Proto)
	// Send a message back to the client
	w.Write([]byte("Hello"))
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got connection: %s", r.Proto)
	// Handle req,

	if r.URL.Path == "/2nd" {
		log.Println("Handling 2nd")
		w.Write([]byte("Hello Again!(2nd)"))
		return
	} else {
		log.Println("Handling Default")
		w.Write([]byte("Hello Again!(Default)"))
		return
	}
	log.Println("Handling 1st")
	// Server push는 resp body가 작성되기 전에 수행되어야 한다.
	// http 커넥션이 push를 지원하는지 체크하기 위해서, resp writer에
	// type assertion을 해본다.
	pusher, ok := w.(http.Pusher)

	if !ok {
		log.Println("Can't push to client")
	} else {
		if err := pusher.Push("/2nd", nil) ; err != nil {
			log.Printf("Failed push : %v", err)
		}
	}
	// Resp body
	w.Write([]byte("Hello"))
}
