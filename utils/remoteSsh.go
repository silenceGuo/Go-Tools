package utils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

type cli struct {
	user       string
	pwd        string
	addr       string
	sshKeyPath string
	client     *gossh.Client
	session    *gossh.Session
	LastResult *[]string
}

// 远程连接ssh服务器，执行命令，run可以实时获取执行结果，类似ping, 返回*[]string，
// cmd 获取执行结果，返回*[]string
// 不过不能目前共用session，命令执行不连贯
// 不能制定执行命令的工作目录
func publicKeyAuthFunc(kPath string) gossh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		ZapLogger.Fatal("查找密钥的主目录失败", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		ZapLogger.Fatal("ssh 密钥文件读取失败", err)
	}
	// Create the Signer for this private key.
	signer, err := gossh.ParsePrivateKey(key)
	if err != nil {
		ZapLogger.Fatal("ssh 关键签名失败", err)
	}
	return gossh.PublicKeys(signer)
}

// 连接对象
func (c *cli) Connect() (*cli, error) {
	config := &gossh.ClientConfig{}
	config.SetDefaults()
	config.User = c.user
	if c.pwd != "" {
		config.Auth = []gossh.AuthMethod{gossh.Password(c.pwd)}
		ZapLogger.Info("用户名,密码连接:", c.addr)
	}
	if c.sshKeyPath != "" {
		config.Auth = []gossh.AuthMethod{publicKeyAuthFunc(c.sshKeyPath)}
		ZapLogger.Info("密钥连接:", c.addr)
	}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key gossh.PublicKey) error { return nil }
	config.Timeout = 10 * time.Second
	client, err := gossh.Dial("tcp", c.addr, config)
	if nil != err {
		ZapLogger.Fatal("连接失败:", err)
		return c, err
	}
	c.client = client
	return c, nil
}
func (c *cli) RunTerminal(shell string, stdout, stderr io.Writer) error {
	//交互式命令行 liunx 环境
	session, err := c.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = stdout
	session.Stderr = stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}
	// Set up terminal modes
	modes := gossh.TerminalModes{
		gossh.ECHO:          1,     // enable echoing
		gossh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		gossh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}
	session.Run(shell)
	return nil
}

func (c *cli) GetSession() (*gossh.Session, error) {
	if c.client == nil {
		if _, err := c.Connect(); err != nil {
			return nil, err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}
	// 关闭会话
	ZapLogger.Info("新建session")
	return session, nil
}

// 执行shell
func (c *cli) Cmd(shell string) (*[]string, error) {
	// 不可以实时输出长执行命令结果，类似于ping
	session, err := c.GetSession()
	if err != nil {
		return nil, err
	}
	//关闭会话
	defer session.Close()
	ZapLogger.Info("执行命令:", shell)
	buf, err := session.CombinedOutput(shell)
	if err != nil {
		ZapLogger.Error("stderr", err)
	}
	var buflist []string
	if buf != nil {
		buflist = strings.Split(string(buf), "\n")
		for _, line := range buflist {
			line = strings.Trim(line, "")
			if line != "" {
				fmt.Println(line)
			}
		}
	}
	c.LastResult = &buflist
	return c.LastResult, err
}

func (c *cli) Run(shell string) (*[]string, error) {
	// 可以实时输出长执行命令结果，类似于ping,也可以 ls
	session, err := c.GetSession()
	if err != nil {
		return nil, err
	}
	// 关闭会话
	defer session.Close()

	ZapLogger.Info("执行命令:", shell)
	stdout, err := session.StdoutPipe() //未获取错误信息
	if err != nil {
		ZapLogger.Error("stdout", err)
	}

	stderr, err := session.StderrPipe() //获取错误信息
	if err != nil {
		ZapLogger.Error("stderr", err)
	}

	execErr := session.Start(shell)
	if execErr != nil {
		ZapLogger.Error(fmt.Sprintf("执行命令错误:%s,%s", shell, execErr))
		return nil, execErr
	}
	o := readerout(stdout)
	if o != nil {
		c.LastResult = &o
		return c.LastResult, nil
	}
	e := readerout(stderr)
	if e == nil {
		c.LastResult = &e
		return c.LastResult, nil
	}
	c.LastResult = &e
	return c.LastResult, errors.New(e[0])
}

func readerout(read io.Reader) []string {
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
		//ZapLogger.Info(line)
		fmt.Println(line)
		stdouts = append(stdouts, line)
	}
	return stdouts
}

func NewSshClient(user string, pwd string, addr string, sshKeyPath string) *cli {
	return &cli{
		user:       user,
		pwd:        pwd,
		addr:       addr,
		sshKeyPath: sshKeyPath,
		client:     nil,
		session:    nil,
		LastResult: nil,
	}
}
