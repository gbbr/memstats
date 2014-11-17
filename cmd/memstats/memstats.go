package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
)

var (
	// laddr is the address the HTTP templates will be served on.
	laddr = flag.String("http", ":8080", "HTTP address to listen on")
	// saddr is the address the HTTP page will attempt to connect to via websockets.
	saddr = flag.String("sock", "localhost:6061", "Adress the WebSockets listen on.")
	// web templates
	tpl *template.Template
)

// serveHTTP serves the front-end HTML/JS viewer
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
	tpl, err = template.New("name").Parse(templateString)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", serveHTTP)
	err = http.ListenAndServe(*laddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
