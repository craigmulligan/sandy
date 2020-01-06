package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type RequestPermission func(string) (string, error)

type Request struct {
	path    string
	syscall string
	allowed bool
}

func Exec(bin string, args []string, requestPermission RequestPermission) {
	var regs syscall.PtraceRegs
	cmd := exec.Command(bin, args...)
	cmd.Stderr = nil
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}
	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		// fmt.Printf("Wait returned: %v\n", err)
	}
	pid := cmd.Process.Pid
	exit := true

	for {
		if exit {
			err = syscall.PtraceGetRegs(pid, &regs)
			if err != nil {
				break
			}

			// https://stackoverflow.com/questions/33431994/extracting-system-call-name-and-arguments-using-ptrace
			if regs.Orig_rax == 0 {
				path, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%d", pid, regs.Rdi))

				if err != nil {
					log.Fatal(err)
				}

				_, err = requestPermission(path)
				if err != nil {
					log.Fatal(err)
				}
			}
			// TODO: print syscall parameters
		}
		err = syscall.PtraceSyscall(pid, 0)
		if err != nil {
			panic(err)
		}
		_, err = syscall.Wait4(pid, nil, 0, nil)
		if err != nil {
			panic(err)
		}
		exit = !exit
	}
}

func askPermission(path string) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(fmt.Sprintf("Wanting to use %s, type y to allow and n to deny.", path))
	for scanner.Scan() {
		if scanner.Text() == "y" {
			break
		}
		if scanner.Text() == "n" {
			fmt.Println("Tried to access a file it is not allowed to.")
			return "", errors.New("Not Allowed")
		}

		fmt.Println("Sorry I didn't understand")
	}
	return path, nil
}

func main() {
	Exec(os.Args[1], os.Args[2:], askPermission)
}
