package help

import (
	"fmt"

	"github.com/fatih/color"
)

// Print Help
func PrintHelp() {
	PrintLogo()
	color.Red("\n\n\tVersion 1.04\n")
	color.Red("\tAuthor: Rana Nadir\t\tBuilt in Go.\n\n")
	color.Blue("\n\n\tServer Reloader for Multi-Threaded Applications, But can work with any platform.\n\n")
	fmt.Println("\tCONFIG: ")
	color.Green("\tglive.json file format: \n\n\t{\n\t\t\"watch\": \"<directory>\", \n\t\t\"command\": \"<command to run>\", \n\t\t\"skip\": \"<[\"path(s)/to\", \"/skip\"]>\"\n\t}")
	fmt.Println("\n\tUSAGE: ")
	color.Blue("\tIn same Directory as glive.json file, Run: \n\n")
	color.Green("\tglive\n\n")
	color.Green("\twhile running, type 'rs' & 'enter' to restart manually.\n\n")
	fmt.Println("\tYou can start any Process: with Multiple Threads, watch for changes in any Directroy, and Skip Multiple Paths. Enjoy!!")
}

func PrintLogoSmall() {
	fmt.Println(`
        ╔═╗┬  ┬┬  ┬┌─┐
        ║ ╦│  │└┐┌┘├┤
        ╚═╝┴─┘┴ └┘ └─┘`)
}

func PrintLogo() {
	fmt.Println(`
	    ▄████  ██▓     ██▓ ██▒   █▓ █████
	   ██▒ ▀█  ██▒    ▓██▒ ██░   █▒ █   ▀
	  ▒██░▄▄▄ ▒██░    ▒██▒ ▓██  █▒  ███
	  ░▓█  ██ ▒██░    ░██░  ▒██ █░░▒▓█  ▄
	  ░▒▓███▀▒░██████ ░██░   ▒▀█░   ▒████▒
	  ░▒   ▒ ░ ▒░▓  ░ ▓     ░ ▐░  ░░ ▒░ ░
	  ░ ░   ░  ░ ░    ▒ ░     ░░     ░
	        ░    ░    ░        ░     ░  ░`)
}
