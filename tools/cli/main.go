package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

func main() {
	var ip *int = pflag.Int("ip", 1234, "help message for flagname")
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.Parse()
	fmt.Println(*ip)
}
