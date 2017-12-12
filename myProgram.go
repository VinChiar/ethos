package main 

import (
	"log"
	"ethos/syscall"
	"ethos/altEthos"
	"strings"
)

func PrintPrompt() {

	thisDirectory = "."

	me := syscall.GetUser()
	path = "/user/" + me 
	prompt := "[" + me + " @ " + path + "]: "
	
	log.Printf("%v", prompt)

}

func is_cmd (cmd string)(result bool) {

	switch cmd {
		case "ls":
			log.Printf("Command ls\n")
			result = true
		case "cd":
			log.Printf("Command cd\n")
			result = true
		default :
			log.Printf("Invalid command\n")
			result = false
		}

	return

}

var path string
var thisDirectory string

func main () {

	PrintPrompt()
	var status syscall.Status
	_, status = altEthos.DirectoryOpen(path)
	if status != syscall.StatusOk {
		log.Fatalf("DirectoryOpen: %v\n", status)
	}

	cmd1 := MyCommands {"echo ciao\nls"}

	status = altEthos.Write(path + "/input", &cmd1)
	if status != syscall.StatusOk {
		log.Fatalf("Write: %v\n", status)
	}

	status = altEthos.Read(path + "/input", &cmd1)
	if status != syscall.StatusOk {
		log.Fatalf("Read: %v\n", status)
	}
	commands := strings.Split(cmd1.Commands, "\n")

	for i, cmd := range commands {

		log.Printf("Command %v: %v\n", i, cmd)

	}

}
