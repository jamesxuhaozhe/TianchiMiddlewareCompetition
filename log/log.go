package log

import (
	"fmt"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.SugaredLogger

var logFilePath string

// InitLogger inits the logger from zap. No log archive is enabled.
func InitLogger() {
	//logFilePath = getLogFilePath()

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout) /**, getLogWriter()**/), level)
	tempLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Logger = tempLogger.Sugar()
}

func getLogFilePath() string {
	serverPort := conf.GetServerPort()
	if serverPort == constants.ClientProcessPort1 {
		return "logfile/client1.log"
	}
	if serverPort == constants.ClientProcessPort2 {
		return "logfile/client2.log"
	}
	if serverPort == constants.BackendProcessPort1 {
		return "logfile/backend.log"
	}
	return "logfile/server.log"
}

func getLogWriter() zapcore.WriteSyncer {
	// create file if not present
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("open file err:%v", err)
	}
	return zapcore.AddSync(file)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	Logger.Debugf(template, args...)
}
func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	Logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

func Fatal(args ...string) {
	Logger.Fatal(args)
}

func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}
