package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
	case pluginName + ":updates":
		n, err := getLinuxAvailableUpdates()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("%d updates available\n", n)
	case pluginName + ":exec":
		args := ""
		for i, s := range flag.Args() {
			if i > 0 {
				args += s
			}
		}
		fmt.Printf("exec(): %s\n", args)
		// out, err := exec.Command("bash", "-c", args).Output()
		out, err := exec.Command(args).Output()
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

func getLinuxAvailableUpdates() (n int, err error) {
	if runtime.GOOS == "darwin" {
		return
	}
	// for non-Debian (Ubuntu) see cosandr/go-check-updates
	// apt-get -q -y --ignore-hold --allow-change-held-packages --allow-unauthenticated -s dist-upgrade
	cmd := exec.Command("apt-get", "-q", "-y", "--ignore-hold", "--allow-change-held-packages", "--allow-unauthenticated", "-s", "dist-upgrade")
	/*
		OUTPUT:
			Reading package lists...
			Building dependency tree...
			Reading state information...
			Calculating upgrade...
			The following packages will be upgraded:
				apport python3-apport python3-problem-report
			3 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.
			Inst python3-problem-report [2.20.11-0ubuntu82.1] (2.20.11-0ubuntu82.2 Ubuntu:22.04/jammy-updates [all])
			Inst python3-apport [2.20.11-0ubuntu82.1] (2.20.11-0ubuntu82.2 Ubuntu:22.04/jammy-updates [all])
			Inst apport [2.20.11-0ubuntu82.1] (2.20.11-0ubuntu82.2 Ubuntu:22.04/jammy-updates [all])
			Conf python3-problem-report (2.20.11-0ubuntu82.2 Ubuntu:22.04/jammy-updates [all])
			Conf python3-apport (2.20.11-0ubuntu82.2 Ubuntu:22.04/jammy-updates [all])
			Conf apport (2.20.11-0ubuntu82.2 Ubuntu:22.04/jammy-updates [all])
		OR:
			Reading package lists...
			Building dependency tree...
			Reading state information...
			Calculating upgrade...
			0 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.
		OR:
			Reading package lists...
			Building dependency tree...
			Reading state information...
			Calculating upgrade...
			The following packages have been kept back:
			  qemu-guest-agent
			0 upgraded, 0 newly installed, 0 to remove and 1 not upgraded.

		Alternatively:
		/usr/lib/update-notifier/apt-check which only returns two numbers: 0;0 = first is number of upgrades, second number is security related upgrades (already in first number)

	*/
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		n = 0
		return
	}
	l := strings.Split(out.String(), "\n")
	upgrade := true
	for i := range l {
		if strings.Contains(l[i], "The following packages will be upgraded:") {
			upgrade = true
		}
		if strings.Contains(l[i], "The following packages have been kept back:") {
			upgrade = false
		}
		// list of updatable packages can span several lines and only (?) those lines always (?) start with two spaces
		if upgrade && strings.HasPrefix(l[i], "  ") {
			// zlog.Trace().Msgf("found system updates: '%s'", l[i])
			u := strings.Split(strings.TrimSpace(l[i]), " ")
			n += len(u)
		}
	}
	return
}
