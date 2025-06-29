package main

import (
	"pgtoch/cmd"
	"pgtoch/internal/log"
)

func main() {
	log.InitLogger()
	cmd.Execute()
}
