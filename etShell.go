package main

import (
	"ethos/syscall"
	"ethos/altEthos"
	"fmt"
	"strings"
)

//Global variables
var shellName string
var path string
var home String
var currDirectory String

// Print the prompt line with the user name and the current working directory
func PrintPrompt() {

	var prompt String

	me := syscall.GetUser()
	path = "/user/" + me
	prompt = String("[" + String(me) + " @ " + currDirectory + "]: ")

	altEthos.WriteStream(syscall.Stdout, &prompt)

}

// Separates each token by one or more space 
func ParseCommand(line String) (cmd string, args CommandList, nArg uint32){

	args = make(CommandList)
	//Process line byte by byte and every time a space is met, populate either cmd or args
	//Count characters in line and number of tokens
	var i uint32
	i = 0
	var k uint32
	k = 0
	var t uint32

	nArg = 0
	f := 0

	for i<4 {
		args[i] = ""
		i++
	}

	i = 0

	for i<uint32(len(line)) && k<5 && line[i] != '\n' {

		t = i
		var tmp string
		tmp = ""

		for line[t] == ' ' {
			t++
		}

		for t<uint32(len(line)) && line[t] != ' ' && line[t] != '\n' {
			tmp = fmt.Sprintf("%s%c", tmp, line[t])
			t++
		}

		if k == 0 {

			cmd = fmt.Sprintf("%s", tmp)

		} else {

			if tmp == "<" || tmp == ">" {
				f = 1
			}
			if tmp != "" {
				if f == 0 {
					nArg++
				}
			}
			args[k-1] = fmt.Sprintf("%s", tmp)
		}

		k++
		i = t

	}


	return

}

// Converts the CommandList dictionary into a []String
func GetArgsArray(cList CommandList) (cmdArray []String){

	//var i uint32
	var nArgs uint32
	nArgs = uint32(len(cList))

	cmdArray = make([]String, nArgs)

	for i, arg := range cList {

		cmdArray[i] = String(arg)

	}

	return

}

// Check the number of arguments and call the exec according to it
func WrapExec(cmd string, args []String, nArg uint32) (status syscall.Status){


	path := "/programs/"+cmd

	if nArg == 0 {

		status = altEthos.Exec(path)

	} else if nArg == 1 {

		status = altEthos.Exec(path, &args[0])

	} else if nArg == 2 {

		status = altEthos.Exec(path, &args[0], &args[1])

	} else if nArg == 3 {

		status = altEthos.Exec(path, &args[0], &args[1], &args[2])

	} else if nArg == 4 {

		status = altEthos.Exec(path, &args[0], &args[1], &args[2], &args[3])

	}

	return

}

// Manage a command containint a redirect
func Redirect(cmd string, args []String, nArg uint32)(status syscall.Status) {


	status = altEthos.Close(syscall.Stdout)
	if status != syscall.StatusOk {

		shellStatus := String("Close failed\n")
		altEthos.WriteStream(syscall.Stdout, &shellStatus)
		return

	}

	fd, status := altEthos.DirectoryOpen(string(args[nArg+1]))
	if status != syscall.StatusOk {
		shellStatus := String("DirectoryOpen failed\n")
		altEthos.WriteStream(syscall.Stderr, &shellStatus)
		return

	}


	//status = altEthos.MoveFd(fd, syscall.Stdout)

	//if status != syscall.StatusOk {

	//	shellStatus := String("MoveFd failed\n")
	//	altEthos.WriteStream(fdd, &shellStatus)
	//	return

	//}


	WrapExec(cmd, args, nArg)

	status = altEthos.MoveFd(syscall.Stdout, fd)
	if status != syscall.StatusOk {
		shellStatus := String("MoveFd failed\n")
		altEthos.WriteStream(syscall.Stderr, &shellStatus)
		return

	}

	status = altEthos.Close(fd)
	if status != syscall.StatusOk {

		shellStatus := String("Close failed\n")
		altEthos.WriteStream(syscall.Stdout, &shellStatus)
		return

	}

	return

}

// Check if the command is in the list of the available commands
func IsCmd(cmd string)(b bool) {

	var i uint32

	progPath := "/programs/"

	files, status := altEthos.SubFiles(progPath)
	if status != syscall.StatusOk {

		shellStatus := String("Subfiles failed\n")
		altEthos.WriteStream(syscall.Stdout, &shellStatus)

	}

	b = false

	for i=0; i<uint32(len(files)); i++ {

		if files[i] == cmd {
			b = true
		}

	}

	return
}

// Manage the change of directory both virtual and actual directory
func HandleCd(firstArg String){

	folder := firstArg

	if firstArg == "" {

		folder = "/"

	}

	status := altEthos.Chdir(string(folder))

	if status != syscall.StatusOk {

		shellStatus := String("cd: not a directory\n")
		altEthos.WriteStream(syscall.Stdout, &shellStatus)

	} else {

		if firstArg == "" {

			currDirectory = home

		} else if firstArg == ".." {

			var temps []string
			var i uint32

			temps = strings.Split(string(currDirectory), "/")
			currDirectory = ""

			for i=0; i<uint32(len(temps)-2); i++ {
				currDirectory = String(currDirectory + String(temps[i]) + "/")
			}

		} else {

			currDirectory = String(currDirectory + firstArg + "/")

		}

	}

}

func main (){

	var pid syscall.ProcessId
	var newPid syscall.ProcessId
	var shellStatus String
	var cmd_line String

	shellName = "etShell"
	home = "home/"
	currDirectory = home

	// print the prompt line at the first shell call
	PrintPrompt()

	// read first command from standard input
	status := altEthos.ReadStream(syscall.Stdin, &cmd_line)
	if status != syscall.StatusOk {

		shellStatus = String("Read failed\n")
		altEthos.WriteStream(syscall.Stdout, &shellStatus)

	}

	// loop until exit is called
	for cmd_line != "exit\n" {

		// extract the arguments 
		cmd, args, nArg := ParseCommand(cmd_line)

		//if IsCmd(cmd) {
		if true {

			var argArr []String

			// put the arguments in an array
			argArr = GetArgsArray(args)

			// if the command is a cd, follow some additional steps
			if cmd == "cd" {

				// manage cd command
				HandleCd(argArr[0])

			}

			// get parent pid
			pid = altEthos.GetPid()

			// fork the program
			_, status = altEthos.Fork(0)
			if status != syscall.StatusOk {

				shellStatus = String("Fork failed\n")
				altEthos.WriteStream(syscall.Stdout, &shellStatus)

			}

			// get pid after forking
			newPid = altEthos.GetPid()

			//Parent process
			if newPid == pid {

				// parent waits
				time := altEthos.GetTime()
				_ = altEthos.Beep(time+6*100000000)


			} else {
			//Child process
				// the child ignore the cd command
				if cmd != "cd" {

					// check if the command contains a redirect
					if argArr[nArg] == ">" || argArr[nArg] == "<" {

						// execute the routine to manage the redirect
						_ = Redirect(cmd, argArr, nArg)

					} else {

						// manage the number of arguments
						status := WrapExec(cmd, argArr, nArg)
						if status != syscall.StatusOk {

							shellStatus = String("[" + shellName + "] Command not found: " + cmd + "\n")
							altEthos.WriteStream(syscall.Stdout, &shellStatus)

						}

					}

				}

				// the child exits when the command is terminated
				altEthos.Exit(syscall.StatusOk)

			}

		} else {

			shellStatus = String("[" + shellName + "] Command not found: " + cmd + "\n")
			altEthos.WriteStream(syscall.Stdout, &shellStatus)

		}


		// print the prompt line for the next command
		PrintPrompt()

		// read next command from standard input
		status := altEthos.ReadStream(syscall.Stdin, &cmd_line)
		if status != syscall.StatusOk {

			shellStatus = String("Read failed")
			altEthos.WriteStream(syscall.Stdout, &shellStatus)

		}

	}


	return

}
