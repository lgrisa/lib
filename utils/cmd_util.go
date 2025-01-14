package utils

import (
	"fmt"
	"github.com/lgrisa/lib/utils/logutil"
	"github.com/pkg/errors"
	"io"
	"os/exec"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
)

func asyncLog(reader io.ReadCloser) error {
	bucket := make([]byte, 1024)
	buffer := make([]byte, 100)
	for {
		num, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "closed") {
				err = nil
			}
			return err
		}
		if num > 0 {
			line := ""
			bucket = append(bucket, buffer[:num]...)
			tmp := string(bucket)
			if strings.Contains(tmp, "\n") {
				ts := strings.Split(tmp, "\n")
				if len(ts) > 1 {
					line = strings.Join(ts[:len(ts)-1], "\n")
					bucket = []byte(ts[len(ts)-1]) //不够整行的以后再处理
				} else {
					line = ts[0]
					bucket = bucket[:0]
				}
				fmt.Printf("%s\n", line)
			}

		}
	}
}

func RunSyncCommand(shellCmd string) error {
	cmd := exec.Command("sh", "-c", shellCmd)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		logutil.LogErrorF("Error starting command: %s......", err.Error())
		return err
	}

	wg := sync.WaitGroup{}
	loadErrorRef := &atomic.Value{}

	wg.Add(2)
	go func() {
		defer func() {
			if errRecover := recover(); errRecover != nil {
				loadErrorRef.Store(errRecover)
				debug.PrintStack()
			}

			wg.Done()
		}()

		if err := asyncLog(stdout); err != nil {
			loadErrorRef.Store(err)
		}
	}()

	go func() {
		defer func() {
			if errRecover := recover(); errRecover != nil {
				loadErrorRef.Store(errRecover)
				debug.PrintStack()
			}

			wg.Done()
		}()

		if err := asyncLog(stderr); err != nil {
			loadErrorRef.Store(err)
		}
	}()

	wg.Wait()

	if loadError := loadErrorRef.Load(); loadError != nil {
		return errors.Errorf("loadErrorRef.Load() error: %v", loadError)
	}

	if err := cmd.Wait(); err != nil {
		logutil.LogErrorF("Error waiting for command execution: %s......", err.Error())
		return err
	}

	return nil
}

func RunCommandGetOutPut(cmd string) ([]byte, error) {
	logutil.LogTraceF("RunCommandGetOutPut: %s", cmd)

	pwdCmd := exec.Command("sh", "-c", cmd)
	return pwdCmd.CombinedOutput()
}
