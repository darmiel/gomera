package main

import "gomera/internal/gomera"

func main() {
	opt := gomera.Parse()
	gomera.New(opt)
}
