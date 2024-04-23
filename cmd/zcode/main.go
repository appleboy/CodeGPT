package main

import (
	"github.com/carsonfeng/CodeGPT/cmd"

	"github.com/appleboy/graceful"
)

func main() {
	m := graceful.NewManager()
	cmd.Execute(m.ShutdownContext())
}
