package main

import (
	"net/http"
	"zap_sample/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var sugarLogger *zap.SugaredLogger

func test_log() {
	name := "test"
	logger.Debugf("this is debug msg: %s", name)
	logger.Infof("this is info msg: %s", name)
	logger.Warnf("this is warn msg: %s", name)
	logger.Errorf("this is error msg: %s", name)
}

func main() {
	// InitLogger()
	// defer sugarLogger.Sync()
	// simpleHttpGet("www.sogo.com")
	// simpleHttpGet("http://www.sogo.com")

	// logger.InitLogger("./new_test.log")
	logger.InitLogger("/dev/stdout")

	name := "wolf"
	logger.Debugf("this is debug msg: %s", name)
	logger.Infof("this is info msg: %s", name)
	logger.Warnf("this is warn msg: %s", name)
	logger.Errorf("this is error msg: %s", name)
	// logger.Fatalf("this is fatal msg: %s", name)

	test_log()

	// for i := 0; i < 10000; i++ {
	// 	logger.Infof("test log roator, hidddkfjdkfjkhidddkfjdkfjkdjfkdjfkdjfkdjfkjdfkjdkfjdkjfkdjfkjdfkdjfkjdkfjdkjfkjdjfkdjfkdjfkdjfkjdfkjdkfjdkjfkdjfkjdfkdjfkjdkfjdkjfkj")
	// }

}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}
