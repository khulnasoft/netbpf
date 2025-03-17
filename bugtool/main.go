package main

import (
	"fmt"
	"os"

	"github.com/khulnasoft/netbpf/bugtool/cmd"
)

func main() {
	if err := cmd.BugtoolRootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
