package log

import (
	"os"
)

var (
	sourceStdOutFd = os.Stdout
	sourceStdErrFd = os.Stderr
)

func redirectStdOut(file *os.File) {
	os.Stdout = file
}

func redirectStdErr(file *os.File) {
	os.Stderr = file
}

func keepStdOut() {
	os.Stdout = sourceStdOutFd
}

func keepStdErr() {
	os.Stderr = sourceStdErrFd
}
