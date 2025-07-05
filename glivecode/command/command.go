package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/clubcleaver/glive/help"
	"github.com/fatih/color"
)

type Cmd = *exec.Cmd

func RunCommand(cmd string) (process Cmd, err error) {
	name, args := splitArgs(cmd)

	process = exec.Command(name, args...)
	process.Stdout = os.Stdout
	process.Stdin = os.Stdin
	process.Stderr = os.Stderr

	process.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Starting the process again
	err = process.Start()
	if err != nil {
		err = fmt.Errorf("process error, Ref: Could not start: %v, Error: %v", cmd, err.Error())
		return
	}
	// help.PrintLogoSmall()
	fmt.Printf("%-35s%v\n%-35s%v\n%-35s%v\n%v", color.BlueString("\tStarted New Process: "), cmd, color.BlueString("\tTime Started: "), time.Now().Format(time.DateTime), color.BlueString("\tNew Process ID: "), process.Process.Pid, color.WhiteString("\t------------------------------------------------------------------\n\n"))
	return
}

func RestartCommand(process Cmd, cmd string) (newProcess Cmd, err error) {

	// Attempt to kill Existing Process
	err = KillProcess(process)
	if err != nil {
		return
	}

	// Start New Process
	help.PrintLogoSmall()
	newProcess, err = RunCommand(cmd)
	if err != nil {
		return
	}
	return

}

// Split Arguments
func splitArgs(input string) (string, []string) {

	parts := strings.Fields(input)
	name := parts[0]
	args := parts[1:]
	return name, args

}

func KillProcess(process Cmd) error {
	if process.Process != nil {
		color.Red("Killing Existing Process: %v\n", process.Process.Pid)
		// color.White("------------------------------------------------------------------")
	}
	if err := syscall.Kill(-process.Process.Pid, syscall.SIGINT); err != nil {
		return fmt.Errorf("process error, Ref:  Killing process group: %v, Error: %v", process.Process, err.Error())
	}

	err := process.Wait() // May Return Error even though process eventually ends.

	if err != nil {
		fmt.Println(err.Error())
	}

	if process.ProcessState != nil {
		if process.ProcessState.Exited() {
			fmt.Println("Process Exited ...")
		} else if process.ProcessState.ExitCode() != 0 {
			fmt.Println("Exited with Code: ", process.ProcessState.ExitCode())
		}
	}

	return nil
}
