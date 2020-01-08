package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/gobwas/glob"
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

		// Make a sounds
		fmt.Printf("\a")
	}
	return Request{path, "READ", true}, nil
}

func Exec(bin string, args []string, allowedPatterns []string) (map[string]Request, error) {
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
	_ = cmd.Wait()
	pid := cmd.Process.Pid

	for {
		err := syscall.PtraceGetRegs(pid, &regs)
		if err != nil {
			break
		}

		// https://stackoverflow.com/questions/33431994/extracting-system-call-name-and-arguments-using-ptrace
		if regs.Orig_rax == 0 {
			path, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%d", pid, regs.Rdi))

			if err != nil {
				return nil, err
			}

			for _, pattern := range allowedPatterns {
				// we match the incoming paths against the input whitelist
				// and automatically create approved requests.
				g := glob.MustCompile(pattern)
				matched := g.Match(path)

				if matched {
					matchedReq := Request{path, "READ", true}
					reqs[path] = matchedReq
				}

				if err != nil {
					return nil, err
				}
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

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var allowedPattern arrayFlags

	// TODO add sane defaults like libc etc
	allowedPattern = append(allowedPattern, "")

	flag.Var(&allowedPattern, "y", "Some description for this param.")
	flag.Parse()

	args := flag.Args()

	_, err := Exec(args[0], args[1:], allowedPattern)
	if err != nil {
		fmt.Println(err)
	}
}
