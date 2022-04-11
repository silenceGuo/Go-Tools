package utils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"strings"
)

type Cli struct {
	user       string
	pwd        string
	addr       string
	sshKeyPath string
	client     *gossh.Client
	session    *gossh.Session
	LastResult string
}

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
func (c *Cli) Connect() (*Cli, error) {
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
	client, err := gossh.Dial("tcp", c.addr, config)
	if nil != err {
		ZapLogger.Fatal("连接失败:", err)
		return c, err
	}
	c.client = client
	return c, nil
}

// 执行shell
func (c Cli) Run(shell string) (string, error) {
	// 不可以实时输出长执行命令结果，类似于ping
	if c.client == nil {
		if _, err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	// 关闭会话
	defer session.Close()
	buf, err := session.CombinedOutput(shell)
	//session.Stderr
	c.LastResult = string(buf)
	return c.LastResult, err
}
func (c Cli) Cmd(shell string) ([]string, error) {
	// 可以实时输出长执行命令结果，类似于ping,也可以 ls
	if c.client == nil {
		if _, err := c.Connect(); err != nil {
			return []string{}, err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return []string{}, err
	}
	// 关闭会话
	defer session.Close()
	stdout, err := session.StdoutPipe() //未获取错误信息
	if err != nil {
		ZapLogger.Error("stdout", err)
	}

	stderr, err := session.StderrPipe() //获取错误信息
	if err != nil {
		ZapLogger.Error("stderr", err)
	}

	execErr := session.Start(shell)
	//session.Start(shell)
	if execErr != nil {
		ZapLogger.Error(fmt.Sprintf("执行命令错误:%s,%s", shell, execErr))
		return []string{}, execErr
	}
	o := readerout(stdout)
	if o != nil {
		return o, nil
	}
	e := readerout(stderr)
	return e, errors.New(e[0])
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

func NewSshClient(user string, pwd string, addr string, sshKeyPath string) *Cli {
	return &Cli{
		user:       user,
		pwd:        pwd,
		addr:       addr,
		sshKeyPath: sshKeyPath,
		client:     nil,
		session:    nil,
		LastResult: "",
	}
}

func Testssh() {
	cli := Cli{
		addr: "192.168.254.11:22",
		user: "root",
		//pwd:  "123456",
		sshKeyPath: "C:\\Users\\Administrator\\Desktop\\id_rsa_2048-dev",
	}
	// 建立连接对象
	fmt.Println("ss", cli.sshKeyPath)
	c, _ := cli.Connect()
	// 退出时关闭连接
	defer c.client.Close()
	res, _ := c.Run("ls ")
	res1, _ := c.Run("pwd")
	fmt.Println(res)
	fmt.Println(res1)

}
