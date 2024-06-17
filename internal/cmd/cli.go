package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"mymodule/internal/storage"
	"os"

	"github.com/google/shlex"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type CLI struct {
	Storage *storage.Storage
}

func InitCLI(path string) *CLI {
	storage := storage.InitStorage(path)
	cli := CLI{
		Storage: &storage,
	}
	return &cli
}

func (cli *CLI) Exit() {
	fmt.Println("Closing database...")
	cli.Storage.CloseStorage()
	fmt.Println("Exit")
}

func (cli *CLI) Execute(command string) {
	commandSplit, err := shlex.Split(command)
	if err != nil {
		log.Fatal(err)
	}
	cmd := commandSplit[0]
	args := commandSplit[1:]
	switch cmd {
	case "help":
		fmt.Println("help")
		return
	case "diff":
		changes, err := cli.Storage.FindDiffs()
		if err != nil {
			fmt.Println(err)
		}
		for _, change := range changes {
			fmt.Printf("fileName: %s\n%s\n", change.FileName, diffmatchpatch.New().DiffPrettyText(change.Changes))
		}
		return
	case "commit":
		cli.commit(args)
	case "branch":
		cli.branch(args)
	case "checkout":
	case "exit":
		cli.Exit()
		os.Exit(0)
		return
	default:
	}
}

func (cli *CLI) commit(args []string) {
	if len(args) == 0 {
		fmt.Printf("Wrong usage of commit. Type \"commit help\" for help")
		return
	}
	option := args[0]
	args = args[1:]
	switch option {
	case "help":
		fmt.Printf("Command \"commit\" used for operating with commits.\n")
		fmt.Printf("Usage: commit <options> \n")
		fmt.Printf("Available options: \n\n")
		fmt.Printf("  %-6s - Show usage help (this msg)\n", "help")
		fmt.Printf("  %-6s - Usage commit help\n", "")
		return

	case "show":
		if len(args) == 0 {
			fmt.Printf("Wrong usage of commit. Type \"commit help\" for help")
			return
		}
		hash, err := hex.DecodeString(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		commit, err := cli.Storage.GetCommit(hash)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Commit", commit)
		fmt.Printf("Description")
		fmt.Printf("Author")
		fmt.Printf("Timr")
	case "branch":
		cli.branch(args)
		// Origin      []byte		//Reference to prev commit
		// Tree        []byte
		// Author      []byte
		// Time        int64
		// Description []byte
	}

}

func (cli *CLI) branch(args []string) {
	// branch
	//
	// if (len(args) == 0){
	// 	branches, current := cli.Storage.GetBranches()
	// 	for _, b  := range branches {
	// 		fmt.Printf("%s", b)
	// 		if (b == current) {
	// 			fmt.Printf(" (current)")
	// 		}
	// 		fmt.Printf("\n")
	// 	}
	// 	return
	// }
	// branchName = args[0]
	// commitsCount := 10
	// args = args[1:]
	// for _, arg := range args {
	// 	switch (arg) {
	// 	case "-a":
	// 		commitsCount = 0
	// 	case "-h":

	// 	}
	// }
}
