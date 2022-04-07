package utils

import (
	"bufio"

	"github.com/prometheus/tsdb/fileutil"
	"io"
	"os"
	"sort"
	"strings"
)

func DirExist(pathstr string) (bool, error) {
	_, err := os.Stat(pathstr)
	if err == nil {
		//ZapLogger.Info()
		ZapLogger.Info("目录存在:", pathstr)
		//fmt.Println("目录存在:",pathstr)
		return true, err
	}
	ZapLogger.Error("目录不存在:", err)
	return false, err
}
func IsHave(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	//index的取值：0 ~ (len(str_array)-1)
	return index < len(str_array) && str_array[index] == target
}

func CreateDir(InitFartherDir string) error {
	return os.MkdirAll(InitFartherDir, 0o755)
}

// go my file tools...111
func ReadFilebufe(filepath string) (*[]string, error) {
	//读取文件，返回 切片指针
	filesile := make([]string, 1)
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		ZapLogger.Error("open file err=", err)
		return nil, err
	}
	const defaultBufSize = 4096
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		str = strings.Trim(str, "\r\n")
		str = strings.Trim(str, "\n")
		filesile = append(filesile, str)
	}
	ZapLogger.Info("文件读取结束")
	return &filesile, nil
}
func WriteFile(filepath string, inputstr interface{}) {
	// 文件追加内容，没有新建
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		ZapLogger.Error("write file err=", err)
		file, _ = os.OpenFile(filepath, os.O_CREATE, 0666)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	switch inputstr.(type) {
	case string:
		converInput, _ := inputstr.(string)
		write.WriteString(converInput + "\n")
	case []string:
		converInput, _ := inputstr.([]string)
		for _, input := range converInput {
			write.WriteString(input + "\n")
		}
	default:
		ZapLogger.Error("inputstr only string | []string")
	}
	write.Flush()
}

func PathExists(path string) (bool, error) {
	//判断路径是否存在，
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Copy(src string, dst string) (written int64, err error) {
	//复制文件
	srcFile, err := os.Open(src)
	defer srcFile.Close()
	if err != nil {
		ZapLogger.Error("open file err=", err)
		return 0, err
	}
	reader := bufio.NewReader(srcFile)
	dstfile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0666)
	defer dstfile.Close()
	if err != nil {
		ZapLogger.Error("open file err=", err)
		return 0, err
	}
	writer := bufio.NewWriter(dstfile)
	return io.Copy(writer, reader)
}
func CopyDirs(src, dst string) error {
	if b, err := DirExist(src); !b {
		return err
	}
	err := fileutil.CopyDirs(src, dst)
	if err != nil {
		ZapLogger.Error("复制目录错误:", err)
		return err
	}
	ZapLogger.Info("复制目录:%s至:%s", src, dst)
	return nil
}
