package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
)

func main() {
	var mode string
	var parameter string

	// Ключи для командной строки
	flag.StringVar(&mode, "mode", "", "mode of the check [urldecode, urlencode]")
	flag.StringVar(&parameter, "p", "", "incoming value (parameter)")

	flag.Parse()

	if mode == "urldecode" {
		query, err := url.QueryUnescape(parameter)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Result: ", query)
	} else if mode == "urlencode" {
		query := url.QueryEscape(parameter)
		fmt.Println("Result: ", query)
	} else {
		fmt.Println("ERROR: unrecognized mode, valid = [urldecode, urlencode]")
	}
}
