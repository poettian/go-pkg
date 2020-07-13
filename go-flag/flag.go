package main

import (
	"flag"
	"fmt"
	"os"
)

var name string

func init()  {
	flag.StringVar(&name, "name", "everyone", "The greeting object.")
	flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "question")
		flag.PrintDefaults()
	}
}

func main()  {
	flag.Parse()

	fmt.Printf("Hello, %s!\n", name)
}