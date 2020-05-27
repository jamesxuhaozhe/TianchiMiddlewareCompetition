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

func InitLogger() {

	//logFilePath = getLogFilePath()
	// define the place where you want to store the log file
	//writeSyncer := getLogWriter()

	// define encoding
	//encoder := getEncoder()

	// init Logger
	//core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	/*	rawJSON := []byte(`{
		  "level": "debug",
		  "encoding": "json",
		  "outputPaths": ["stdout", "logfile/backend.log"],
		  "errorOutputPaths": ["stderr"],
		  "initialFields": {"foo": "bar"},
		  "encoderConfig": {
		    "messageKey": "message",
		    "levelKey": "level",
		    "levelEncoder": "lowercase"
		  }
		}`)

		var cfg zap.Config
		if err := json.Unmarshal(rawJSON, &cfg); err != nil {
			panic(err)
		}
		tempLogger, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		Logger = tempLogger.Sugar()*/
	tempLogger, _ := zap.NewDevelopment()
	Logger = tempLogger.Sugar()
	zap.AddCaller()
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

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
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
