package main

import (
	"github.com/thibaultmg/clingua/cmd/clingua/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
