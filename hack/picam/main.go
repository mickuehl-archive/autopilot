/*
 * Copyright (C) 2018 Josh A. Beam
 * All rights reserved.
 *
 * See https://github.com/joshb/pi-camera-go
 */

package main

import (
	"flag"
	"fmt"

	"shadow-racer/autopilot/v1/test/picam/picam"
)

func main() {
	address := flag.String("address", "0.0.0.0:10042", "The address (including port) to bind to")
	useHTTPS := flag.Bool("https", false, "Use HTTPS")
	flag.Parse()

	s, err := picam.New(*useHTTPS)
	if err != nil {
		fmt.Println("Unable to create server:", err)
		return
	}

	if err := s.Start(*address); err != nil {
		fmt.Println("Unable to start server:", err)
	}
}
