package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/clubcleaver/glive/command"
	"github.com/clubcleaver/glive/config"
	"github.com/clubcleaver/glive/help"
	"github.com/clubcleaver/glive/traversal"
	"github.com/fatih/color"
)

func main() {
	hlp := flag.Bool("h", false, "Display Help")
	flag.Parse()
	if *hlp {
		help.PrintHelp()
		os.Exit(0)
	}

	conf, err := config.GetConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	referenceMap, err := traversal.GetPathAndMod(conf.Directory, conf.Skip)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	// INTRO
	help.PrintLogoSmall()
	skipDirs := strings.Join(conf.Skip, ", ")
	fmt.Printf("\n%-35s%v\n", color.BlueString("\tWatching Directory:"), conf.Directory)
	fmt.Printf("%-35s%v\n\n", color.BlueString("\tSkipping Path(s):"), skipDirs)

	// Start Process
	var process command.Cmd
	process, err = command.RunCommand(conf.Command)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	// Capture Interrupt
	var InterruptChan = make(chan os.Signal, 1)
	signal.Notify(InterruptChan, os.Interrupt)

	go func() {
		<-InterruptChan
		color.Red("\nInterrupt Detected")
		err = command.KillProcess(process)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		color.Red("Exiting Glive ...")
		os.Exit(0)
	}()

	var restartChan = make(chan struct{}, 1)
	go func(resChan chan<- struct{}) {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			input := scanner.Text()
			if input == "rs" {
				resChan <- struct{}{} // Signal Restart
			}
		}
	}(restartChan)

	// The LOOP
	for {
		time.Sleep(time.Millisecond * 100)
		newRefMap, err := traversal.GetPathAndMod(conf.Directory, conf.Skip)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}

		for path, modTime := range newRefMap {
			existingModTime, ok := referenceMap[path]
			if !ok {
				color.Red("New File Detected: %v", path)
				referenceMap[path] = modTime // New watch File Added, So Replace Ref and Restart
				process, err = command.RestartCommand(process, conf.Command)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(0)
				}
				break
			}
			if existingModTime != modTime {
				color.Red("Modification detected in file: %v", path)
				referenceMap[path] = modTime // File Update, So Replace Ref and Restart

				process, err = command.RestartCommand(process, conf.Command)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(0)
				}
				break
			}
		}

		// Listen for Restart Signal
		if len(restartChan) > 0 {
			<-restartChan
			color.Red("Restart Requested ...")
			process, err = command.RestartCommand(process, conf.Command)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}
		}
	}
}
