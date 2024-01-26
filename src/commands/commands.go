package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	columnize "github.com/ryanuber/columnize"
)

const (
	pluginName = "test-plugin"
	helpHeader = `Usage: dokku ` + pluginName + `[:COMMAND]

Runs commands that interact with the app's repo

Additional commands:`

	helpContent = `
	` + pluginName + `:test, prints test message
`
)

func main() {
	flag.Usage = usage
	flag.Parse()

	cmd := flag.Arg(0)
	switch cmd {
	case pluginName + ":exec":
		args := ""
		for _, s := range flag.Args() {
			args += s
		}
		fmt.Printf("exec(): %s\n", flag.Args())
		out, err := exec.Command("bash", "-c", args).Output()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", out)
	case pluginName + ":help":
		usage()
	case "help":
		fmt.Print(helpContent)
	case pluginName + ":test":
		fmt.Println("triggered: " + pluginName + " from: commands")
	default:
		dokkuNotImplementExitCode, err := strconv.Atoi(os.Getenv("DOKKU_NOT_IMPLEMENTED_EXIT"))
		if err != nil {
			fmt.Println("failed to retrieve DOKKU_NOT_IMPLEMENTED_EXIT environment variable")
			dokkuNotImplementExitCode = 10
		}
		os.Exit(dokkuNotImplementExitCode)
	}
}

func usage() {
	config := columnize.DefaultConfig()
	config.Delim = ","
	config.Prefix = "\t"
	config.Empty = ""
	content := strings.Split(helpContent, "\n")[1:]
	fmt.Println(helpHeader)
	fmt.Println(columnize.Format(content, config))
}
