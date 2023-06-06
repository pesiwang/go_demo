package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
)

func InitLogger(fileName string) {

	core := zapcore.NewCore(
		getEncoder(),
		getWriteSyncer(fileName),
		zap.InfoLevel,
	)
	logger = zap.New(core, zap.AddCallerSkip(1), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getWriteSyncer(fileName string) zapcore.WriteSyncer {
	if fileName == "/dev/stdout" {
		return zapcore.AddSync(os.Stdout)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    2048,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func innerGoLogger() *zap.Logger { return logger }

func getLogger() *zap.Logger {
	return innerGoLogger()
}

func Debug(msg string, fields ...zap.Field)  { getLogger().Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)   { getLogger().Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)   { getLogger().Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field)  { getLogger().Error(msg, fields...) }
func DPanic(msg string, fields ...zap.Field) { getLogger().DPanic(msg, fields...) }
func Panic(msg string, fields ...zap.Field)  { getLogger().Panic(msg, fields...) }
func Fatal(msg string, fields ...zap.Field)  { getLogger().Fatal(msg, fields...) }
func Sync()                                  { getLogger().Sync(); logger.Sync() }

func S() *zap.SugaredLogger { return getLogger().Sugar() }

func Debugw(msg string, keysAndValues ...interface{}) { S().Debugw(msg, keysAndValues...) }
func Infow(msg string, keysAndValues ...interface{})  { S().Infow(msg, keysAndValues...) }
func Warnw(msg string, keysAndValues ...interface{})  { S().Warnw(msg, keysAndValues...) }
func Errorw(msg string, keysAndValues ...interface{}) { S().Errorw(msg, keysAndValues...) }
func Panicw(msg string, keysAndValues ...interface{}) { S().Panicw(msg, keysAndValues...) }
func Fatalw(msg string, keysAndValues ...interface{}) { S().Fatalw(msg, keysAndValues...) }

func Debugf(msg string, keysAndValues ...interface{}) { S().Debugf(msg, keysAndValues...) }
func Infof(msg string, keysAndValues ...interface{})  { S().Infof(msg, keysAndValues...) }
func Warnf(msg string, keysAndValues ...interface{})  { S().Warnf(msg, keysAndValues...) }
func Errorf(msg string, keysAndValues ...interface{}) { S().Errorf(msg, keysAndValues...) }
func Panicf(msg string, keysAndValues ...interface{}) { S().Panicf(msg, keysAndValues...) }
func Fatalf(msg string, keysAndValues ...interface{}) { S().Fatalf(msg, keysAndValues...) }
