package main

import (
	"bufio"
	"fmt"
	"os"

	"mymodule/internal/cmd"

	"github.com/xlab/closer"
)

const path = "tmp"

var cli *cmd.CLI

func main() {
	closer.Bind(cleanup)
	go initCLI()
	closer.Hold()
}

func initCLI() {
	cli = cmd.InitCLI(path)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Print("> ")
		cli.Execute(scanner.Text())
	}
}
func cleanup() {
	fmt.Println("Closing DB...")
	if cli != nil {
		cli.Exit()
	}
	fmt.Println(" Done.")
}
