package utils

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

func SplitFilePath(filePath string) []string {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	filenameList := strings.Split(base, ".")
	filenameList = append(filenameList, dir)
	return filenameList
}

//var s *viper.Viper
func GetYAML(yamlPath string) *viper.Viper {

	v := viper.New()
	pathlist := SplitFilePath(yamlPath)
	//Viper.SetConfigName()
	v.SetConfigName(pathlist[0])
	v.SetConfigType(pathlist[1])
	v.AddConfigPath(pathlist[2])
	//Vipers.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Println("err=", err)
	}
	v.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被更新啦...")
		//if err := viper.Unmarshal(Conf); err != nil {
		//	panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		//}
	})
	//keysALL := v.AllKeys()
	//for _,k := range keysALL {
	//	fmt.Println(v)
	//
	//	fmt.Println(v.Get(k))
	//}
	//viper.Get()
	//return viper

	//viper.Get()
	return v

}

func WriteYaml(yamlPath string, k string, v interface{}) {
	w := viper.New()
	pathlist := SplitFilePath(yamlPath)
	//Viper.SetConfigName()
	w.SetConfigName(pathlist[0])
	w.SetConfigType(pathlist[1])
	w.AddConfigPath(pathlist[2])
	w.SetConfigFile(yamlPath)
	//fmt.Println(pathlist)
	w.Set(k, v)
	ZapLogger.Info("写入键值对：%s-%s", k, v)
	w.WriteConfig()
}
