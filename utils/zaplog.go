package utils

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var ZapLogger *zap.SugaredLogger

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	ZapLogger = logger.Sugar()
	defer ZapLogger.Sync()
	//sugarLogger.Info("test")
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "D:\\MyGoProject\\logs\\test.log", // ⽇志⽂件路径
		MaxSize:    1,                                 // 1M=1024KB=1024000byte
		MaxBackups: 5,                                 // 最多保留5个备份
		MaxAge:     30,                                // days
		Compress:   false,                             // 是否压缩 disabled by default
	}
	//return zapcore.AddSync(lumberJackLogger)//只输出到文件
	//输出到控制台和文件
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}
