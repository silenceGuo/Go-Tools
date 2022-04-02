package utils

import "testing"

func TestReadFilebufe(t *testing.T) {
	ReadFilebufe("C:\\Users\\Administrator\\Desktop\\startServerdev19.yml")
}

func TestWriteFile(t *testing.T) {
	WriteFile("C:\\Users\\Administrator\\Desktop\\startServerdev1.yml", "1")
	WriteFile("C:\\Users\\Administrator\\Desktop\\startServerdev1.yml", 1)
	WriteFile("C:\\Users\\Administrator\\Desktop\\startServerdev1.yml", []string{"1", "sd"})
}
