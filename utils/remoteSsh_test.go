package utils

import (
	"fmt"
	"testing"
)

func TestCli_Run(t *testing.T) {
	newclinet := NewSshClient("root",
		"123456",
		"192.168.254.26:22",
		"")
	newclinet.Run("ping www.baidu.com -w 3")
	fmt.Println(*newclinet.LastResult)
	newclinet.Run("ls -al")
	fmt.Println(*newclinet.LastResult)
}
func TestCli_Cmd(t *testing.T) {
	newclinet := NewSshClient("root",
		"123456",
		"192.168.254.26:22", "")
	//newclinet.Cmd("cd /tmp")
	//fmt.Println(*newclinet.LastResult)
	newclinet.Cmd("pwd")
	fmt.Println(*newclinet.LastResult)
}
func TestCli_GetSession(t *testing.T) {
	newclinet := NewSshClient("root",
		"123456",
		"192.168.254.26:22",
		"")
	session, err := newclinet.GetSession()
	if err != nil {
		fmt.Println(err)
	}
	s2, err := session.Output("cd /tmp && pwd")
	fmt.Println(string(s2))
	fmt.Println(err)
}
