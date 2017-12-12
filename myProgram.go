package main 

import (
	"log"
	"ethos/syscall"
	"ethos/kernelTypes"
	"ethos/altEthos"
)

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

func main () {

	var myReader kernelTypes.String

	me := syscall.GetUser()

	path = "[" + me + "|" + "~]:"

	statusR := altEthos.ReadStream(syscall.Stdin, &myReader)
	if statusR != syscall.StatusOk {

	}
	statusW := altEthos.WriteStream(syscall.Stdout, &myReader)
	if statusW != syscall.StatusOk {}

	for {
		//if is_cmd(cmd) {
		//execute cmd
		//}
	}

}
