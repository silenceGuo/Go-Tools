package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ReadFilebufe(filepath string) *[]string {
	//读取文件，返回 切片指针
	filesile :=make([]string,1)
	file,err :=os.Open(filepath)
	defer file.Close()
	if err !=nil{
		fmt.Println("open file err=",err)
	}
	const defaultBufSize  = 4096
	reader := bufio.NewReader(file)
	for {
		str,err:= reader.ReadString('\n')
		if err == io.EOF{
			break
		}
		str = strings.Trim(str,"\r\n")
		str = strings.Trim(str,"\n")
		filesile = append(filesile,str)
	}
	fmt.Println("文件读取结束")
	return &filesile
}
func WriteFile(filepath string,inputstr interface{}) {
    // 文件追加内容，没有新建
	file,err := os.OpenFile(filepath,os.O_RDWR|os.O_APPEND,0666)
	if err !=nil{
		fmt.Println("write file err=",err)
		file,_=os.OpenFile(filepath,os.O_CREATE,0666)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	switch inputstr.(type) {
	  case string:
		  converInput,_ := inputstr.(string)
		  write.WriteString(converInput+"\n")
	  case []string:
		  converInput,_ := inputstr.([]string)
		  for _,input:= range converInput{
			write.WriteString(input+"\n")
	      }
	  default:
		fmt.Println("inputstr only string | []string")
	}
	write.Flush()
}

func PathExists(path string)(bool,error){
	//判断路径是否存在，
	_,err := os.Stat(path)
	if err == nil{
		return true,nil
	}
	if os.IsNotExist(err){
		return false,nil
	}
	return false,err
}

func Copy(src string,dst string)(written int64,err error)  {
	//复制文件
	srcFile,err := os.Open(src)
	defer srcFile.Close()
	if err != nil{
		fmt.Println("open file err=",err)
		return
	}
	reader := bufio.NewReader(srcFile)
	dstfile,err:= os.OpenFile(dst,os.O_WRONLY|os.O_CREATE,0666)
	defer dstfile.Close()
	if err!=nil{
		fmt.Println("open file err=",err)
		return
	}
	writer:= bufio.NewWriter(dstfile)
	return io.Copy(writer,reader)
}