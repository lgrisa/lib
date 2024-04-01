package utils

import (
	"bytes"
	"fmt"
	"github.com/disgoorg/log"
	"github.com/pkg/errors"
	"io"
	"os/exec"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
)

func RunCommand(name string, arg ...string) error {

	cmd := exec.Command(name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()

	if err != nil {
		return errors.Wrapf(err, "执行命令失败，name: %s, arg: %v, stdout: %s, stderr: %s", name, arg, stdout.String(), stderr.String())
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	if len(outStr) > 0 {
		fmt.Println(outStr)
	}

	if len(errStr) > 0 {
		return errors.Errorf(errStr)
	}

	return nil
}

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
		log.Errorf("Error starting command: %s......", err.Error())
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
		log.Errorf("Error waiting for command execution: %s......", err.Error())
		return err
	}

	return nil
}

func RunCommandGetOutPut(cmd string) ([]byte, error) {
	log.Tracef("RunCommandGetOutPut: %s", cmd)

	pwdCmd := exec.Command("sh", "-c", cmd)
	return pwdCmd.CombinedOutput()
}
