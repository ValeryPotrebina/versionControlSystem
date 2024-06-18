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
	fmt.Print("> ")
	for scanner.Scan() {
		cli.Execute(scanner.Text())
		fmt.Print("> ")
	}
}
func cleanup() {
	fmt.Println("Closing DB...")
	if cli != nil {
		cli.Exit()
	}
	fmt.Println(" Done.")
}
