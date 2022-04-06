package utils

import (
	"bytes"
	"os"
	"text/template"
)

var tmpPath *string

func FileCreate(content bytes.Buffer, name string) error {
	file, err := os.Create(name)
	if err != nil {
		ZapLogger.Error("错误信息:", err)

		return err
	}
	_, err = file.WriteString(content.String())
	if err != nil {
		ZapLogger.Error("错误信息:", err)
		return err
	}
	file.Close()
	return nil
}

func GenFileFromTmp(tmpfile string, data map[string]interface{}, outFile string) {
	t1, err := template.ParseFiles(tmpfile)
	ZapLogger.Info("模板文件:", tmpfile)
	//data = make(map[string]interface{})
	if err != nil {
		ZapLogger.Error("错误信息:", err)
		return
		//panic(err)
	}
	var b1 bytes.Buffer
	//将上面获得的模板数据进行相关的数据绑定（即将m中数据输入到字节流中）
	err = t1.Execute(&b1, data)
	err = FileCreate(b1, outFile)
	if err != nil {
		ZapLogger.Error("生成模板文件错误:", err)
	}
	ZapLogger.Info("生成模板文件:", outFile)

}

func GenTmpFromString(tmpstr string, data map[string]interface{}) string {
	//通过传入的字符串模板和map 返回格式化好的字符串
	//tmp "{{.env}}"
	//data[env]=test
	tmp1, err := template.New("envNameTmp").Parse(tmpstr)
	if err != nil {
		panic(err)
	}
	//data["envName"] = global.GoDevOpsFlagConfig.EnvName
	var b1 bytes.Buffer
	err = tmp1.Execute(&b1, data)
	if err != nil {
		panic(err)
	}
	s := b1.String()
	return s
	//fmt.Println(s)
}
