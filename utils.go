package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

//ExecShell is basically equivalent to "echo "$stdinPipe" | command arg"
func ExecShell(command string, arg string, stdinPipe string) (string, int) {
	var exitCode int

	cmd := exec.Command(command, arg)
	stdin, err := cmd.StdinPipe()
	CheckError(err, true)

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, stdinPipe)
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		// if failed because of exit status, don't log to console
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.Sys().(syscall.WaitStatus).ExitStatus()
		} else {
			fmt.Printf("err is %v\n", err)
		}
	} else {
		exitCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	}

	return string(output), exitCode
}

// Opens text files and returns its content as string
// Also is able to return os.Stdin if argument "stdin" is passed
func ReadTextFile(filename string) string {
	var fileContent string

	if filename == "stdin" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fileContent += scanner.Text() + "\n"
		}

		fileContent = strings.TrimSuffix(fileContent, "\n")
	} else {
		data, err := ioutil.ReadFile(filename)
		CheckError(err, true, "Error reading file")
		fileContent = string(data)
	}

	return fileContent
}

// Check error and do the needful.
// message is printed in the error message as "%s: %s",  message, err.
// If message defaults to "An error ocurred"
// fatal determines if execution should stop or not.
// log.Fatal() if true, log.Default().Print() otherwise
func CheckError(err error, fatal bool, argsMessage ...string) {
	var message string

	if len(argsMessage) == 0 {
		message = "An error ocurred"
	} else {
		message = strings.Join(argsMessage, "")
	}
	if err != nil {
		if fatal {
			log.Fatalf("%s: %v", message, err)
		} else {
			log.Default().Printf("%s: %v", message, err)
		}
	}
}

func ReadLine(prompt string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(prompt)

	input, err := reader.ReadString('\n')
	CheckError(err, false, "Error reading from stdin")

	return input
}
