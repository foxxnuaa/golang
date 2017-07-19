package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func testPipeSingle() {
	cmd := exec.Command("echo", "-n", "My first command from golang")

	/*获取命令的输出*/
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:Cannot not obtain the stdout pipe for command No.0:%s\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error:The command No.0 can not be startup:%s\n", err)
		return
	}

	var outputBuf bytes.Buffer
	for {
		tempOutput := make([]byte, 5)
		n, err := stdout.Read(tempOutput)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Error:Can not read data from the pipe:%s\n", err)
				return
			}
		}

		if n > 0 {
			outputBuf.Write(tempOutput[:n])
		}
	}

	fmt.Printf("%s\n", outputBuf.String())
}

func testPipeDual() {
	fmt.Println("Run command `ps aux | grep apipe`: ")

	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "apipe")

	stdout1, err := cmd1.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:Cannot not obtain the stdout pipe for command No.0:%s\n", err)
		return
	}

	if err := cmd1.Start(); err != nil {
		fmt.Printf("Error:The command No.0 can not be startup:%s\n", err)
		return
	}

	outputBuf1 := bufio.NewReader(stdout1)
	stdin2, err := cmd2.StdinPipe()
	if err != nil {
		fmt.Printf("Error:Cannot not obtain the stdout pipe for command No.0:%s\n", err)
		return
	}

	outputBuf1.WriteTo(stdin2)

	var outputBuf2 bytes.Buffer
	cmd2.Stdout = &outputBuf2
	if err := cmd2.Start(); err != nil {
		fmt.Printf("Error:The command can not be startup:%s\n", err)
		return
	}

	err = stdin2.Close()
	if err != nil {
		fmt.Printf("Error:Can not close the stdio pipe:%s\n", err)
		return
	}

	if err := cmd2.Wait(); err != nil {
		fmt.Printf("Error:Can not wait for the command:%s\n", err)
		return
	}

	fmt.Printf("%s\n", outputBuf2.Bytes())
}
