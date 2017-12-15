package main

import (
	"os"
	"ethos/altEthos"
	"ethos/syscall"
)


func main() {

	var msg String
	var desc String
	var out String

	arg := String(os.Args[1])

	if len(os.Args) > 2 {

		out = "Usage: help command_name\n"

	} else {

		switch(arg) {

		case "cd":

			msg = "Change working directory"
			desc = "dir_name"

		case "date":

			msg = "Get current date and time"
			desc = ""

		case "echo":

			msg = "Print a message to standard output"
			desc = "message"

		case "ls":

			msg = "Print list of directories and files"
			desc = "path"

		default :

			msg = "No manual entry"
			desc = ""

		}

		out = arg + " -- " + msg + "\n"
		if desc != "" {
			out = out + "Usage: " + arg + " " + desc + "\n"
		}

	}

	_ = altEthos.WriteStream(syscall.Stdout, &out)

	return

}
