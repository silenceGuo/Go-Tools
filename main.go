package main

import (
	"github.com/silenceGuo/Go-Tools/utils"
)

func main() {

	//utils.InitLogger()
	utils.ZapLogger.Info("aaa")
	//utils.ZapLogger.Info("11","ss")
	//utils.ExecCommand("ping1 www.baidu.com","/")
	//utils.ExecCommandlast("ping www.baidu.com","/")
	//utils.ExecCommand("df1 -h","/")
	//utils.ExecCommand2("df -h","/")
	//fmt.Println(utils.ExecCommand("ping www.baidu.com -c 5","/"))
	//utils.ExecCommand2("ping www.baidu.com -c 5","/")
	//utils.ExecCommand3("ping www.baidu.com","/")
	//y:=utils.GetYAML("C:\\Users\\Administrator\\Desktop\\nginxConf-dev.yml")
	//y:=utils.GetYAML("C:\\Users\\Administrator\\Desktop\\1.yml")
	//date :=make(map[string]interface{},2)
	//utils.YmalToMap(y)
	//utils.Testssh()
	c := utils.NewSshClient(
		"root",
		"123456",
		"192.168.254.26:22",
		"",
	)
	//c := utils.NewSshClient("root","123456","192.168.254.11:22")
	//res,err:= c.Run("ping www.baidu.com")
	//_,err:= c.Cmd("ping www.baidu.com -w 5")
	//c.Cmd("cd /tmp ")
	//fmt.Println(c.LastResult)
	c.Run("pwd")
	//fmt.Println(c)
	//c.Cmd("ls -al")
	//fmt.Println(c.LastResult)
	//c.RunTerminal("top",os.Stdout,os.Stdin)
	//c.RunTerminal("ls -al",os.Stdout,os.Stdin)
	//fmt.Println(res)
	//fmt.Println(res1)
	//fmt.Println(c.LastResult)
	//if err !=nil{
	//	utils.ZapLogger.Error("ss",err)
	//}
	//fmt.Println("33",res)
	//data := make(map[string]interface{})
	//xmx := project.Requests.StartMemory
	//data["pinpointid"] = "asd"
	//data["xms"] = 123
	//data["xmx"] = 1235
	//utils.GenFileFromTmp("D:\\MyGoProject\\devops-go\\templates\\catalina.sh.tmp2", data, "C:\\Users\\Administrator\\Desktop\\catalina.sh")
}
