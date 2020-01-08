package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Request struct {
	path    string
	syscall string
	allowed bool
}

func requestPermission(path string) (Request, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(fmt.Sprintf("Wanting to READ %s [y/n]", path))
	for scanner.Scan() {
		input := strings.ToLower(scanner.Text())
		if input == "y" {
			break
		}
		if scanner.Text() == "n" {
			return Request{path, "READ", false}, nil
		}

		fmt.Println("Sorry I didn't understand")
	}
	return Request{path, "READ", true}, nil
}

func Exec(bin string, args []string) (map[string]Request, error) {
	var regs syscall.PtraceRegs
	reqs := make(map[string]Request)
	cmd := exec.Command(bin, args...)
	cmd.Stderr = nil
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace:    true,
		Pdeathsig: syscall.SIGKILL,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		// fmt.Printf("Wait returned: %v\n", err)
	}
	pid := cmd.Process.Pid

	for {
		err = syscall.PtraceGetRegs(pid, &regs)
		if err != nil {
			break
		}

		// https://stackoverflow.com/questions/33431994/extracting-system-call-name-and-arguments-using-ptrace
		if regs.Orig_rax == 0 {
			path, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%d", pid, regs.Rdi))

			if err != nil {
				return nil, err
			}

			_, ok := reqs[path]

			if !ok {
				req, err := requestPermission(path)

				if !req.allowed {
					return nil, errors.New(fmt.Sprintf("Blocked %s on %s", req.syscall, req.path))
				}

				reqs[req.path] = req

				if err != nil {
					return nil, err
				}
			}
		}

		err = syscall.PtraceSyscall(pid, 0)
		if err != nil {
			return nil, err
		}

		_, err = syscall.Wait4(pid, nil, 0, nil)
		if err != nil {
			return nil, err
		}
	}
	return reqs, nil
}

func main() {
	_, err := Exec(os.Args[1], os.Args[2:])
	if err != nil {
		fmt.Println(err)
	}
}
