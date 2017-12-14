package main

import (
//	"ethos/log"
	"ethos/syscall"
	"ethos/altEthos"
	"fmt"
//	"ethos/kernelTypes"
	"strings"
)

//Global variables
var shellName string
var path string
var home String
var currDirectory String

func PrintPrompt() {

	var prompt String

	me := syscall.GetUser()
	path = "/user/" + me
	prompt = String("[" + String(me) + " @ " + currDirectory + "]: ")

	altEthos.WriteStream(syscall.Stdout, &prompt)

}

//1.0 separate each token by one or more space 
//2.0 separate cmd and args by one or more space, then inside args define a list of enclosing chars, namely " and ', and count them alternatively. If their final count is 0 in means that for each " and ' there is the coresponding closing and the args are valid. Not perfect, but a start.

//Consider creating a type Command with string field to store the command itself and a flag to indicate if it is a built in command or not

func ParseCommand(line String) (cmd string, args CommandList, nArg uint32){ 

	args = make(CommandList)
	//1.0
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
		//if k == 0 {
		tmp = ""
		//} else {
		//	tmp = " "
		//}

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


	//shellStatus = String("cmd[" + String(cmd) + "]\n")
	//altEthos.WriteStream(syscall.Stdout, &shellStatus)

	//shellStatus = String("ParseCommand: 1:" + args[0] + "2:" + args[1] + ",3:" + args[2] + ",4:" + args[3]+"e")
	//altEthos.WriteStream(syscall.Stdout, &shellStatus)

	return

}

//This function converts the CommandList dictionary into a []String
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

func WrapExec(cmd string, args []String, nArg uint32) (status syscall.Status){


	path := "/programs/"+cmd

	if nArg == 0 {

		//shellStatus := String("Path:" + String(path)+"\n")
		//altEthos.WriteStream(syscall.Stdout, &shellStatus)
		status = altEthos.Exec(path)

	} else if nArg == 1 {

		//shellStatus := String("Path:"+String(path)+" 1:" + args[0]+"\n")
		//altEthos.WriteStream(syscall.Stdout, &shellStatus)
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

func Redirect(cmd string, args []String, nArg uint32)(status syscall.Status) {


	//shellStatus := String("Folder:" + args[nArg+1] + "\n")
	//altEthos.WriteStream(syscall.Stdout, &shellStatus)
	fd, status := altEthos.DirectoryOpen(string(args[nArg+1]))
	
	if status != syscall.StatusOk {

		shellStatus := String("DirectoryOpen failed\n")
		altEthos.WriteStream(syscall.Stdout, &shellStatus)
		return

	}
	
	status = altEthos.Close(syscall.Stdout)
	if status != syscall.StatusOk {

		shellStatus := String("Close failed\n")
		altEthos.WriteStream(fd, &shellStatus)
		return

	}
	
	fdd := fd

	status = altEthos.Close(fd)
	if status != syscall.StatusOk {

		shellStatus := String("Close failed\n")
		altEthos.WriteStream(fd, &shellStatus)
		return

	}

	status = altEthos.MoveFd(fdd, 1)

	if status != syscall.StatusOk {

		shellStatus := String("MoveFd failed\n")
		altEthos.WriteStream(fd, &shellStatus)
		return

	}



	shellStatus := String("Cmd:" + String(cmd) + "\n")
	altEthos.WriteStream(syscall.Stdout, &shellStatus)

	shellStatus = String("Argument:"+args[0]+")\n")
	altEthos.WriteStream(syscall.Stdout, &shellStatus)

	WrapExec(cmd, args, nArg)

	status = altEthos.MoveFd(syscall.Stdout, fd)
	shellStatus = String("All good\n")
	altEthos.WriteStream(syscall.Stdout, &shellStatus)

	return

}

func IsCmd(cmd string)(b bool) {
	
	//shellStatus := String("cmd [" + String(cmd) + "]")
	//altEthos.WriteStream(syscall.Stdout, &shellStatus)

	switch(cmd) {
	case "cd":
		b = true
	case "echo":
		b = true
	case "ps":
		b = true
	case "ls":
		b = true
	default:
		b = false

	}

	return
}

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

	//var cmd_line string
	//var i uint32
	var pid syscall.ProcessId
	var newPid syscall.ProcessId
	//var logger = log.Initialize("test/log/")

	shellName = "pizza"
	home = "home/"
	currDirectory = home
	var shellStatus String
	//t = "----------------------------\n\n"
	//altEthos.WriteStream(syscall.Stdout, &t)

	var cmd_line String


	PrintPrompt()

	status := altEthos.ReadStream(syscall.Stdin, &cmd_line)
	if status != syscall.StatusOk {

		shellStatus = String("Read failed\n")
		altEthos.WriteStream(syscall.Stdout, &shellStatus)

	}

	for cmd_line != "exit\n" {

		//cmd_line = String("echo ciao")

		//shellStatus = String("1\n")
		//altEthos.WriteStream(syscall.Stdout, &shellStatus)

		//shellStatus = String("cmd_line[" + String(cmd_line) + "]\n")
		//altEthos.WriteStream(syscall.Stdout, &shellStatus)

		cmd, args, nArg := ParseCommand(cmd_line)

		//shellStatus = String("2\n")
		//altEthos.WriteStream(syscall.Stdout, &shellStatus)

		//if IsCmd(cmd) {
		if true {

		//logger.Println("Command:", cmd)
			//logger.Println("Args:")

			//for i=0; i<nCommands; i++ {
			//	logger.Printf("%s ", args[i])
			//}

			var argArr []String
			argArr = GetArgsArray(args)

			if cmd == "cd" {

				HandleCd(argArr[0])

			}

			pid = altEthos.GetPid()
			_, status = altEthos.Fork(0)
			if status != syscall.StatusOk {

				shellStatus = String("Fork failed\n")
				altEthos.WriteStream(syscall.Stdout, &shellStatus)

			}

			newPid = altEthos.GetPid()

			//Parent process
			if newPid == pid {


				time := altEthos.GetTime()
				_ = altEthos.Beep(time+6*100000000)


			} else {
			//Child process
				if cmd != "cd" {

					if argArr[nArg] == ">" || argArr[nArg] == "<" {

						//shellStatus := String("****" + argArr[nArg] + "\n")
						//altEthos.WriteStream(syscall.Stdout, &shellStatus)

						_ = Redirect(cmd, argArr, nArg)

					} else {

						//shellStatus = String("1:" + argArr[0] + "2:" + argArr[1] + ",3:" + argArr[2] + ",4:" + argArr[3]+"e")
						//altEthos.WriteStream(syscall.Stdout, &shellStatus)

						status := WrapExec(cmd, argArr, nArg)
						//status := altEthos.Exec(path, &argArr[0], &argArr[1], &argArr[2], &argArr[3])
						if status != syscall.StatusOk {

							shellStatus = String("[" + shellName + "] Command not found: " + cmd + "\n")
							altEthos.WriteStream(syscall.Stdout, &shellStatus)

						}

					}

				}

				altEthos.Exit(syscall.StatusOk)

			}

		} else {

			shellStatus = String("[" + shellName + "] Command not found: " + cmd + "\n")
			altEthos.WriteStream(syscall.Stdout, &shellStatus)

		}


		PrintPrompt()
		status := altEthos.ReadStream(syscall.Stdin, &cmd_line)
		if status != syscall.StatusOk {

			shellStatus = String("Read failed")
			altEthos.WriteStream(syscall.Stdout, &shellStatus)

		}

	}


	return

}
