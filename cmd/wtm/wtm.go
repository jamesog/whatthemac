// wtm (What the MAC) looks up a given MAC address against IEEE registrations.
//
// The lists can be downloaded from
// https://regauth.standards.ieee.org/standards-ra-web/pub/view.html#registries
//
// The program looks for the CSV file names as named on the IEEE site.
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	mac "git.jog.li/jamesog/whatthemac"
)

const (
	largeFile  = "oui.csv"
	mediumFile = "mam.csv"
	smallFile  = "oui36.csv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <MAC address>\n", os.Args[0])
		os.Exit(1)
	}

	hwaddr, err := net.ParseMAC(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	macaddr := strings.Replace(hwaddr.String(), ":", "", -1)

	// Read the Large list by default (MA-L)
	mal, err := os.Open(largeFile)
	if err != nil {
		log.Fatalf("Could not open MA-L list: %v", err)
	}
	defer mal.Close()
	entries, err := mac.ReadAssignments(mal)
	if err != nil {
		log.Fatalf("Could not read MA-L list: %v", err)
	}

	// Turn the MAC address into an upper-cased version of the first 6 chars
	macoui := strings.ToUpper(macaddr[0:6])

	if entry, ok := entries[macoui]; ok {
		fmt.Println(entry.Organization)
		os.Exit(0)
	}

	// Not found in MA-L, so try MA-M

	// Turn the MAC address into an upper-cased version of the first 8 chars
	macoui = strings.ToUpper(macaddr[0:8])

	mam, err := os.Open(mediumFile)
	if err != nil {
		log.Fatalf("Could not open MA-M list: %v", err)
	}
	defer mam.Close()
	entries, err = mac.ReadAssignments(mam)
	if err != nil {
		log.Fatalf("Could not read MA-M list: %v", err)
	}
	if entry, ok := entries[macoui]; ok {
		fmt.Println(entry.Organization)
		os.Exit(0)
	}

	// Not found in MA-M, so try MA-S

	// Turn the MAC address into an upper-cased version of the first 8 chars
	macoui = strings.ToUpper(macaddr[0:9])

	mas, err := os.Open(smallFile)
	if err != nil {
		log.Fatalf("Could not open MA-S list: %v", err)
	}
	defer mas.Close()
	entries, err = mac.ReadAssignments(mas)
	if err != nil {
		log.Fatalf("Could not read MA-S list: %v", err)
	}
	if entry, ok := entries[macoui]; ok {
		fmt.Println(entry.Organization)
		os.Exit(0)
	}

	fmt.Printf("%s not found\n", macoui)
	os.Exit(3)
}
