/**
 * main process
 * Usage: http-mock[ -v]
 * Rosbit Xu
 */
package main

import (
	"os"
	"fmt"
)

// variables set via go build -ldflags
var (
	buildTime string
	osInfo    string
	goInfo    string
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		showVersion()
		return
	}

	if err := CheckGlobalConf(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(3)
	}
	DumpConf()

	if err := StartService(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(4)
	}
	os.Exit(0)
}

func showInfo(prompt, info string) {
	if info != "" {
		fmt.Printf("%10s: %s\n", prompt, info)
	}
}

func showVersion() {
	showInfo("name",       os.Args[0])
	showInfo("build time", buildTime)
	showInfo("os name",    osInfo)
	showInfo("compiler",   goInfo)
}
