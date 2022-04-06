package utils

import (
	"bufio"
	//"devops-go/global"
	"fmt"
	//"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

func readout(read io.ReadCloser) []string {
	var stdouts []string
	reader := bufio.NewReader(read)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			//ZapLogger.Error("cmd StdoutPipe error1:", err)
			break
		}
		if err != nil {
			ZapLogger.Error("read StdoutPipe error:", err)
			break
		}
		line = strings.Trim(line, "\n")
		ZapLogger.Info(line)
		stdouts = append(stdouts, line)
	}
	return stdouts
}

func ExecCommand(cmd_sub string, workespace string) ([]string, bool) {
	//获取 错误，正确输出，返回切片
	var stdouts []string
	if workespace != "" {
		path, _ := os.Getwd()
		if path != workespace {
			ZapLogger.Info("切换到工作目录:", workespace)
			err := os.Chdir(workespace)
			if err != nil {
				ZapLogger.Error("切换工作目录错误:", err)
				return stdouts, false
			}
		}
	}
	cmd := exec.Command("/bin/bash", "-c", cmd_sub)
	stdout, err := cmd.StdoutPipe() //未获取错误信息
	if err != nil {
		ZapLogger.Error(err)
	}
	stderr, err := cmd.StderrPipe() //获取错误信息
	if err != nil {
		ZapLogger.Error(err)
	}
	execErr := cmd.Start()
	if execErr != nil {
		ZapLogger.Error(fmt.Sprintf("执行命令错误:%s,%s", cmd_sub, execErr))
		return []string{execErr.Error()}, false
	}
	o := readout(stdout)
	if o != nil {
		return o, true
	}
	e := readout(stderr)
	return e, false

}
