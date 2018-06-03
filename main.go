package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type registration struct {
	List         string // MA-L, MA-M, MA-S
	OUI          string
	Organization string
}

type entries map[string]registration

func readEntries() entries {
	e := make(entries)

	f, err := os.Open("oui.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		e[rec[1]] = registration{
			List:         rec[0],
			OUI:          rec[1],
			Organization: rec[2],
		}
	}

	return e
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <MAC>\n", os.Args[0])
		os.Exit(1)
	}

	mac := os.Args[1]
	if mac == "" {
		fmt.Fprintf(os.Stderr, "MAC not well formatted\n")
		os.Exit(2)
	}

	e := readEntries()

	macoui := strings.ToUpper(strings.Join(strings.Split(mac, ":")[0:3], ""))
	fmt.Printf("Looking for %s\n", macoui)
	en, ok := e[macoui]
	if !ok {
		fmt.Printf("%s not found\n", macoui)
		os.Exit(3)
	}
	fmt.Printf("%s %s\n", en.OUI, en.Organization)

}
