package main

import (
	"ethos/log"
	"ethos/syscall"
	"ethos/altEthos"
	"fmt"
//	"io"
	"ethos/kernelTypes"
//	"ethos/defined"
)

//1.0 separate each token by one or more space 
//2.0 separate cmd and args by one or more space, then inside args define a list of enclosing chars, namely " and ', and count them alternatively. If their final count is 0 in means that for each " and ' there is the coresponding closing and the args are valid. Not perfect, but a start.

//Consider creating a type Command with string field to store the command itself and a flag to indicate if it is a built in command or not

func ParseCommand(line string) (cmd string, args CommandList, nCommands uint32){ 
	
	args = make(CommandList)	
	//1.0
	//Process line byte by byte and every time a space is met, populate either cmd or args
	//Count characters in line and number of tokens
	var i uint32
	i = 0
	var k uint32
	k = 0
	var t uint32

	for i<uint32(len(line)) {

		t = i
		var tmp string
		tmp = ""

		for t<uint32(len(line)) && line[t] != ' ' {
			tmp = fmt.Sprintf("%s%c", tmp, line[t])
			t++
		}

		if k == 0 {
			cmd = fmt.Sprintf("%s", tmp)
		} else {
			args[k-1] = fmt.Sprintf("%s", tmp)
		}

		k++
		i = t+1

	}

	nCommands = k-1

	return

}

func IsCmd(cmd string) (b bool){

	switch cmd {
	
	case "ls":
	case "echo":
		b = true
	default:
		b = false

	}
	
	return
}

func main (){

	var cmd_line string
	var i uint32
	var pid syscall.ProcessId
	var newPid syscall.ProcessId
	var logger = log.Initialize("test/log/")

	cmd_line = "echo eeehi"
	//assume cmd_line is the output of ReadCmdLine()
	cmd, args, nCommands := ParseCommand(cmd_line)
	
	if IsCmd(cmd) {

		logger.Println("Command:", cmd)
		logger.Println("Args:")

		for i=0; i<nCommands; i++ {
			logger.Printf("%s ", args[i])
		}

		pid = altEthos.GetPid()
		_, status := altEthos.Fork(0)
		if status != syscall.StatusOk {

			logger.Println(pid,": Fork error")

		}
		newPid = altEthos.GetPid()
		//Parent process
		if newPid == pid {

			//fmt.Fprint(io.Stdout, "Ehi")
			logger.Println(pid,": parent waiting...")

		} else {
			//Child process
			logger.Println(newPid,": child is executing", cmd)

			path := "/programs/"+cmd
			var arg1 String
			arg1 = String(args[0])

			status := altEthos.Exec(path, &arg1)
			if status != syscall.StatusOk {
				logger.Printf("%v: Exec error: %v", newPid, status)
			}
			altEthos.Exit(syscall.StatusOk)
		}

	} else {
	
		logger.Println("Unrecognized command: %s", cmd)

	}


	//var myWriter String
	//myWriter = "porco dio\n"
	//fmt.Printf("%v\n", syscall.Stdout)
	//if status != syscall.StatusOk{
	//	fmt.Printf("Unable to write to stdout.\n")
	//}
	//var r io.Writer
	var t kernelTypes.String
	t = "ear\n\n"
	altEthos.WriteStream(syscall.Stdout, &t)
	return

}
