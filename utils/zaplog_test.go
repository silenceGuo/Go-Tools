package utils

import "testing"

func TestInitLogger(t *testing.T) {
	//InitLogger()
	ZapLogger.Info("test log")
}
