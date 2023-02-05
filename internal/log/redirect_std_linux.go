package log

import (
	"fmt"
	"os"
	"syscall"
)

var (
	invalidHanlder = -1
	sourceStdOutFd = invalidHanlder
	sourceStdErrFd = invalidHanlder
)

func redirectStdOut(file *os.File) {
	if file == nil {
		return
	}

	if sourceStdOutFd == invalidHanlder {
		if fd, err := syscall.Dup(int(os.Stdout.Fd())); err != nil {
			sourceStdOutFd = fd
		}
	}

	err := syscall.Dup3(int(file.Fd()), 1, 0)
	if err != nil {
		fmt.Println("fail to redirect stdout,", err)
	}
}

func redirectStdErr(file *os.File) {
	if file == nil {
		return
	}

	if sourceStdErrFd == invalidHanlder {
		if fd, err := syscall.Dup(int(os.Stderr.Fd())); err != nil {
			sourceStdErrFd = fd
		}
	}

	err := syscall.Dup3(int(file.Fd()), 2, 0)
	if err != nil {
		fmt.Println("fail to redirect stderr,", err)
	}
}

func keepStdOut() {
	err := syscall.Dup3(sourceStdOutFd, 1, 0)
	if err != nil {
		fmt.Println("fail to un-redirects stdout,", err)
	}
}

func keepStdErr() {
	err := syscall.Dup3(sourceStdErrFd, 2, 0)
	if err != nil {
		fmt.Println("fail to un-redirects stderr,", err)
	}
}
