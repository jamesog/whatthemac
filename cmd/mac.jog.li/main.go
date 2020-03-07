// +build !appengine

//go:generate go-bindata -prefix "../../" -o lists.go ../../csv

package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	mac "git.jog.li/jamesog/whatthemac"
)

var (
	mal, mam, mas mac.Assignments
)

func init() {
	malcsv, err := Asset("csv/oui.csv")
	if err != nil {
		log.Fatalf("Could not load MA-L list: %v", err)
	}
	mamcsv, err := Asset("csv/mam.csv")
	if err != nil {
		log.Fatalf("Could not load MA-M list: %v", err)
	}
	mascsv, err := Asset("csv/oui36.csv")
	if err != nil {
		log.Fatalf("Could not load MA-S list: %v", err)
	}

	mal, err = mac.ReadAssignments(bytes.NewReader(malcsv))
	if err != nil {
		log.Fatalf("Couldn't load MA-L data: %v", err)
	}

	mam, err = mac.ReadAssignments(bytes.NewReader(mamcsv))
	if err != nil {
		log.Fatalf("Couldn't load MA-M data: %v", err)
	}

	mas, err = mac.ReadAssignments(bytes.NewReader(mascsv))
	if err != nil {
		log.Fatalf("Couldn't load MA-S data: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	s := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	http.HandleFunc("/", handle)
	log.Fatal(s.ListenAndServe())
}

func handle(w http.ResponseWriter, r *http.Request) {
	var addr string
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "What the MAC?\n")
		return
	default:
		addr = strings.TrimLeft(r.URL.Path, "/")
	}

	hwaddr, err := net.ParseMAC(addr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Couldn't parse MAC address %s: %v\n", addr, err)
		return
	}

	macaddr := strings.Replace(hwaddr.String(), ":", "", -1)

	oui := strings.ToUpper(macaddr[0:6])
	if entry, ok := mal[oui]; ok {
		fmt.Fprintf(w, "%s\n", entry.Organization)
		return
	}
	oui = strings.ToUpper(macaddr[0:8])
	if entry, ok := mam[oui]; ok {
		fmt.Fprintf(w, "%s\n", entry.Organization)
		return
	}
	oui = strings.ToUpper(macaddr[0:9])
	if entry, ok := mas[oui]; ok {
		fmt.Fprintf(w, "%s\n", entry.Organization)
		return
	}

	fmt.Fprintf(w, "Unknown manufacturer\n")
}
