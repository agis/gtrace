package main

//go:generate go run gen/syscall_gen.go

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

var pid int

func main() {
	flag.IntVar(&pid, "p", 0, "PID of the process to trace")
	flag.Parse()

	if pid == 0 {
		fmt.Println("Provide a valid process ID with -p")
		os.Exit(1)
	}

	err := unix.PtraceAttach(pid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Attached to process %d...\n", pid)

	// PtraceAttach sends a SIGSTOP to the child; ensure it has properly
	// stopped
	s := new(unix.WaitStatus)
	unix.Wait4(pid, s, 0, new(unix.Rusage))

	err = unix.PtraceSetOptions(pid, syscall.PTRACE_O_TRACESYSGOOD)
	if err != nil {
		panic(err)
	}

	for {
		regs := new(unix.PtraceRegs)

		exitCode := waitSyscall(pid)
		if exitCode != 0 {
			fmt.Println("Process exited with", exitCode)
			os.Exit(0)
		}

		err = unix.PtraceGetRegs(pid, regs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s = ", Syscalls[regs.Orig_rax])

		exitCode = waitSyscall(pid)
		if exitCode != 0 {
			fmt.Println("Process exited with", exitCode)
			os.Exit(0)
		}

		err = unix.PtraceGetRegs(pid, regs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d\n", regs.Rax)
	}
}

func waitSyscall(pid int) int {
	s := new(unix.WaitStatus)
	for {
		err := unix.PtraceSyscall(pid, 0)
		if err != nil {
			panic(err)
		}

		unix.Wait4(pid, s, 0, new(unix.Rusage))
		if s.Stopped() && (s.StopSignal()&0x80 > 0) {
			return 0
		} else if s.Exited() {
			fmt.Println("process exited")
			return 1
		}
	}
}
