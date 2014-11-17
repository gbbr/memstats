package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

var (
	laddr = flag.String("http", ":8080", "HTTP address to listen on")
	saddr = flag.String("sock", "localhost:6061", "Adress the WebSockets listen on.")
)

func serveHTTP(w http.ResponseWriter, req *http.Request) {
	if err := tpl.ExecuteTemplate(w, "main", *saddr); err != nil {
		fmt.Fprintf(w, "Error parsing template: %s", err)
	}
}

func main() {
	flag.Parse()
	hst, _, err := net.SplitHostPort(*saddr)
	if len(hst) == 0 || err != nil {
		log.Fatal("sockaddr must be host[:port]. ERR: %s", err)
	}
	http.HandleFunc("/", serveHTTP)
	err = http.ListenAndServe(*laddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
