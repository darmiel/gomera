package main

import "github.com/darmiel/gomera/internal/gomera"

func main() {
	opt := gomera.Parse()
	gomera.New(opt)
}
