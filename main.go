package main

import (
	//"MyGoProject/Go-Tools/utils"
	"fmt"
	"github.com/silenceGuo/Go-Tools/utils"
)

func main() {
	f := *utils.ReadFilebufe("C:\\Users\\Administrator\\Desktop\\startServerdev1.yml")
	for _, v := range f {
		fmt.Printf("%s", string(v))
	}
	var s []string
	s = append(s, "1", "12easdad")
	utils.WriteFile("C:\\Users\\Administrator\\Desktop\\1.TXT", 3)
	fmt.Println(utils.PathExists("C:\\Users\\Administrator\\Desktop\\1.TXT1"))
}
