package main

import (
	"github.com/silenceGuo/Go-Tools/utils"
	//"/myGoProject/Go-Tools/utils"
)

func main() {

	utils.InitLogger()
	//utils.ZapLogger.Info("aaa")
	//utils.ZapLogger.Info("11","ss")
	//utils.ExecCommand("ping1 www.baidu.com","/")
	//utils.ExecCommandlast("ping www.baidu.com","/")
	//utils.ExecCommand("df1 -h","/")
	//utils.ExecCommand2("df -h","/")
	//fmt.Println(utils.ExecCommand("ping www.baidu.com -c 5","/"))
	//utils.ExecCommand2("ping www.baidu.com -c 5","/")
	//utils.ExecCommand3("ping www.baidu.com","/")
	//utils.GetYAML("C:\\Users\\Administrator\\Desktop\\nginxConf-dev.yml")
	//date :=make(map[string]interface{},2)

	data := make(map[string]interface{})
	//xmx := project.Requests.StartMemory
	data["pinpointid"] = "asd"
	data["xms"] = 123
	data["xmx"] = 1235
	utils.GenFileFromTmp("D:\\MyGoProject\\devops-go\\templates\\catalina.sh.tmp2", data, "C:\\Users\\Administrator\\Desktop\\catalina.sh")
}
