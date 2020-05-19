package main

import (
	"github.com/prometheus/common/log"

	"github.com/ditointernet/go-dito/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Error(err)
	}
}
