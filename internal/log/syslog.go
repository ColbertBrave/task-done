package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cloud-disk/internal/constants"
)

type SystemLog struct {
	logFile    string
	logHandler *os.File
	exit       chan struct{}
	waitGroup  sync.WaitGroup
}

func NewSystemLog(file string) *SystemLog {
	sysLog := new(SystemLog)
	sysLog.logFile = file
	sysLog.exit = make(chan struct{})

	sysLog.init()
	go sysLog.monitor()
	return sysLog
}

func (s *SystemLog) init() {
	if s.logHandler != nil {
		return
	}

	var err error
	s.logHandler, err = os.OpenFile(s.logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("fail to create system logs,", err.(*os.PathError).Err)
		return
	}
	redirectStdOut(s.logHandler)
	redirectStdErr(s.logHandler)
}

func (s *SystemLog) monitor() {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	isFileExist := func(filePath string) bool {
		_, err := os.Stat(filePath)
		if err == nil {
			return true
		}
		return os.IsExist(err)
	}

	for {
		select {
		case _, ok := <-s.exit:
			if ok {
				fmt.Println("system logs is closed at", time.Now().Format(constants.TimeFormat))
			}
			return

		case <-ticker.C:
			if !isFileExist(s.logFile) {
				s.reset()
			}
		}
	}
}

func (s *SystemLog) reset() {
	s.close()
	s.init()
}

func (s *SystemLog) close() {
	if s.logHandler != nil {
		if err := s.logHandler.Close(); err != nil {
			fmt.Println("fail to close system logs handler")
		}
	}
	s.logHandler = nil
	keepStdOut()
	keepStdErr()
}
