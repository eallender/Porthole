package main

import (
	"flag"
	"fmt"
)

const version = "0.1.0-dev"

func showUsage() {
	fmt.Println("Usage: porthole [options] <command>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --version   Show version information")
	fmt.Println("  --help      Show this help message")
}

func main() {
	var (
		showVersion = flag.Bool("version", false, "show version information")
		showHelp	= flag.Bool("help", false, "show help information")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("porthole version %s\n", version)
		return
	}

	if *showHelp {
		showUsage()
		return
	}

	fmt.Println("Porthole - Docker Network Validation Tool")
	fmt.Println("Coming soon: network traffic validation!")

	args := flag.Args()
	if len(args) == 0 {
		showUsage()
		return
	}	
}