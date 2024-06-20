package cmd

import (
	"encoding/hex"
	"fmt"
	"mymodule/internal/object"
	"mymodule/internal/storage"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/shlex"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type CLI struct {
	Storage *storage.Storage
}

func InitCLI(path string) *CLI {
	storage, err := storage.InitStorage(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	cli := CLI{
		Storage: storage,
	}
	return &cli
}

func (cli *CLI) Execute(command string) {
	commandSplit, err := shlex.Split(command)
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(commandSplit) == 0 {
		return
	}
	cmd := commandSplit[0]
	args := commandSplit[1:]
	switch cmd {
	case "help":
		fmt.Printf("Available command:\n")
		fmt.Printf("  %-8s - show help\n", "help")
		fmt.Printf("  %-8s - create new commit\n", "commit")
		fmt.Printf("  %-8s - show info about branches\n", "branch")
		fmt.Printf("  %-8s - switch branches\n", "checkout")
		fmt.Printf("  %-8s - show differences between versions\n", "diff")
		fmt.Printf("  %-8s - show info about objects\n", "show")
		fmt.Printf("  %-8s - exit program\n", "exit")

		return
	case "diff":
		cli.diff(args)
		return
	case "commit":
		cli.commit(args)
		return
	case "branch":
		cli.branch(args)
		return
	case "show":
		cli.show(args)
		return
	case "checkout":
		cli.checkout(args)
		return
	case "exit":
		cli.Exit()
		os.Exit(0)
		return
	default:
		fmt.Printf("Unknown command %s. Type \"help\" for help.\n", cmd)
		return
	}
}
func (cli *CLI) commit(args []string) {
	if len(args) == 0 {
		fmt.Printf("Wrong usage of commit. Type \"commit -h\" for help.\n")
		return
	}
	var author string = ""
	var description string = ""

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			fmt.Printf("usage: commit\n")
			fmt.Printf("   or: commit -d <description>\n")
			fmt.Printf("   or: commit -a <author> -d <description>\n")
			fmt.Printf("\n")
			fmt.Printf("Available options\n")
			fmt.Printf("  %-16s    show help (this message)\n", "-h --help")
			fmt.Printf("  %-16s    set commit's author \n", "-a --author")
			fmt.Printf("  %-16s    default - current username\n", "")
			fmt.Printf("  %-16s    set commit's description\n", "-d --description")
			fmt.Printf("  %-16s    default - \"\"\n", "")

			return
		case "-a", "--author":
			if i+1 >= len(args) {
				fmt.Printf("Wrong usage of argument %s. Type \"commit -h\" for help.\n", arg)
				return
			}
			author = args[i+1]
			i++
		case "-d", "--description":
			if i+1 >= len(args) {
				fmt.Printf("Wrong usage of argument %s. Type \"commit -h\" for help.\n", arg)
				return
			}
			description = args[i+1]
			i++
		default:
			fmt.Printf("Unknown argument %s. Type \"commit -h\" for help.\n", arg)
			return
		}
	}
	if author == "" {
		user, err := user.Current()
		if err != nil {
			fmt.Printf("User is not specified. Type \"commit -h\" for help.\n")
		}
		author = user.Username
	}
	err := cli.Storage.CreateCommit(author, description)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Commit created")
}
func (cli *CLI) branch(args []string) {
	if len(args) == 0 {
		branches := cli.Storage.GetBranches()
		for _, branch := range branches {
			fmt.Printf("%s", branch)
			if branch == cli.Storage.Branch {
				fmt.Printf(" (current)")
			}
			fmt.Printf("\n")
		}
		return
	}
	var branch string = ""
	var count uint64 = 5
	var verbose bool = false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			fmt.Printf("usage: branch\n")
			fmt.Printf("   or: branch <branch>\n")
			fmt.Printf("   or: branch <branch> -c <count>\n")
			fmt.Printf("   or: branch <branch> -a -v\n")
			fmt.Printf("\n")
			fmt.Printf("Available options\n")
			fmt.Printf("  %-12s    show help (this message)\n", "-h --help")
			fmt.Printf("  %-12s    enable verbose output\n", "-v --verbose")
			fmt.Printf("  %-12s    show all commits\n", "-a --all")
			fmt.Printf("  %-12s    set commit's limit\n", "-c --count")
			fmt.Printf("  %-12s    default: \"5\"\n", "")
			return
		case "-a", "--all":
			count = 0
		case "-c", "--count":
			if i+1 >= len(args) {
				fmt.Printf("Wrong usage of argument %s. Type \"branch -h\" for help.\n", arg)
				return
			}
			c, err := strconv.ParseUint(args[i+1], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
			}
			count = c
			i++
		case "-v", "--verbose":
			verbose = true
		default:
			if branch == "" {
				branch = arg
			} else {
				fmt.Printf("Unknown argument %s. Type \"branch -h\" for help.\n", arg)
				return
			}
		}
	}
	if branch == "" {
		branch = cli.Storage.Branch
	}
	commits, err := cli.Storage.GetCommits(branch, count)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Branch: %s\n", branch)
	fmt.Printf("\n---------------------------------------------------------------------------\n\n")
	for i, commitData := range commits {
		if verbose {
			fmt.Printf("Commit:        %x\n", commitData.Hash)
			fmt.Printf("Description:   %s\n", commitData.Commit.Description)
			fmt.Printf("Author:        %s\n", commitData.Commit.Author)
			fmt.Printf("Time:          %s\n", time.Unix(commitData.Commit.Time, 0).Format("02.01.2006 15:04:05"))
			fmt.Printf("Origin:        %x\n", commitData.Commit.Origin)
			fmt.Printf("Tree:          %x\n", commitData.Commit.Tree)
			if i != len(commits)-1 {
				fmt.Printf("\n---------------------------------------------------------------------------\n\n")
			}
		} else {
			fmt.Printf("%x\n", commitData.Hash)
		}
	}
}
func (cli *CLI) diff(args []string) {
	var verbose bool = false
	hashes := make([][]byte, 0)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			fmt.Printf("usage: diff\n")
			fmt.Printf("   or: diff -v\n")
			fmt.Printf("   or: diff <commit>\n")
			fmt.Printf("   or: diff <commit1> <commit2>\n")
			fmt.Printf("\n")
			fmt.Printf("Available options\n")
			fmt.Printf("  %-12s    show help (this message)\n", "-h --help")
			fmt.Printf("  %-12s    enable verbose output\n", "-v --verbose")
			return
		case "-v", "--verbose":
			verbose = true
		default:
			if len(hashes) < 2 {
				hash, err := hex.DecodeString(arg)
				if err != nil {
					fmt.Printf("%s is not hash", arg)
					return
				}
				hashes = append(hashes, hash)
			} else {
				fmt.Printf("Unknown argument %s. Type \"diff -h\" for help.\n", arg)
				return
			}
		}
	}
	var changes []*object.FileChange
	var err error
	switch len(hashes) {
	case 0:
		changes, err = cli.Storage.Diffs()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	case 1:
		changes, err = cli.Storage.DiffsWithCommit(hashes[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	case 2:
		changes, err = cli.Storage.DiffsBetweenCommits(hashes[0], hashes[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	for i, c := range changes {
		if verbose {
			fmt.Printf("Filename:  %s\n", c.FileName)
			fmt.Printf("Changes:\n")
			fmt.Println(diffmatchpatch.New().DiffPrettyText(c.Changes))
			if i != len(changes)-1 {
				fmt.Printf("\n---------------------------------------------------\n\n")
			}
		} else {
			delete := 0
			insert := 0
			for _, c := range c.Changes {
				switch c.Type {
				case diffmatchpatch.DiffInsert:
					insert++
				case diffmatchpatch.DiffDelete:
					delete++
				}
			}
			fmt.Printf(
				"%s %s%s\n",
				c.FileName,
				color.RedString("%s", strings.Repeat("-", delete)),
				color.GreenString("%s", strings.Repeat("+", insert)),
			)
		}
	}
}
func (cli *CLI) checkout(args []string) {
	if len(args) == 0 {
		fmt.Printf("Wrong usage of commit. Type \"chechout -h\" for help.\n")
		return
	}

	var b bool = false
	var branch string = ""
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			fmt.Printf("usage: chechout <branch>\n")
			fmt.Printf("   or: chechout -b <branch>\n")
			fmt.Printf("\n")
			fmt.Printf("Available options\n")
			fmt.Printf("  %-9s    show help (this message)\n", "-h --help")
			fmt.Printf("  %-9s    Create branch and switch\n", "-b")
			return
		case "-b":
			b = true
			return
		default:
			if branch == "" {
				branch = arg
			} else {
				fmt.Printf("Unknown argument %s. Type \"chechout -h\" for help.\n", arg)
				return
			}
		}
	}
	if branch == "" {
		fmt.Printf("Branch name is not specified. Type \"chechout -h\" for help.\n")
	}
	if b {
		err := cli.Storage.CreateBranch(branch)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	err := cli.Storage.ChangeBranch(branch)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Current branch is %s.\n", cli.Storage.Branch)
}
func (cli *CLI) show(args []string) {
	if len(args) == 0 {
		fmt.Printf("Wrong usage of show. Type \"show -h\" for help.\n")
		return
	}

	var hash []byte = nil
	var err error

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			fmt.Printf("usage: show <hash>\n")
			fmt.Printf("\n")
			fmt.Printf("Available options\n")
			fmt.Printf("  %-9s    show help (this message)\n", "-h --help")
			return
		default:
			if hash == nil {
				hash, err = hex.DecodeString(arg)
				if err != nil {
					fmt.Printf("%s is not hash", arg)
					return
				}
			} else {
				fmt.Printf("Unknown argument %s. Type \"show -h\" for help.\n", arg)
				return
			}
		}
	}

	if hash == nil {
		fmt.Printf("Wrong usage of show. Type \"show -h\" for help.\n")
		return
	}
	obj, err := cli.Storage.GetObject(hash)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Hash:          %x\n", hash)
	fmt.Printf("Type:          %s\n\n", object.TypeToString(obj.Type))
	switch obj.Type {
	case object.TypeBlob:
		blob, err := obj.ParseBlob()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("---------- Content ----------\n")
		fmt.Println(string(blob.Data))
		fmt.Printf("------------ End ------------\n")

	case object.TypeTree:
		tree, err := obj.ParseTree()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, c := range tree.Children {
			fmt.Printf("%-10s%x      %s\n", object.TypeToString(c.Type), c.Hash, c.Name)
		}
	case object.TypeCommit:
		commit, err := obj.ParseCommit()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Description:   %s\n", commit.Description)
		fmt.Printf("Author:        %s\n", commit.Author)
		fmt.Printf("Time:          %s\n", time.Unix(commit.Time, 0).Format("02.01.2006 15:04:05"))
		fmt.Printf("Origin:        %x\n", commit.Origin)
		fmt.Printf("Tree:          %x\n", commit.Tree)
	}
}

func (cli *CLI) Exit() {
	fmt.Println("Closing database...")
	cli.Storage.CloseStorage()
	fmt.Println("Exit")
}
