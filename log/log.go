package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const(
	debugLevel = 0
	releaseLevel =1
	warnLevel = 2
	errorLevel =3
	fatalLevel =4
)

const (
	printDebugLevel = 	"[debug  ]"
	printReleaseLevel = "[release]"
	printWranLevel = 	"[warn	 ]"
	printErrorLevel = 	"[error  ]"
	printFatalLevel = 	"[fatal  ]"
)

type Logger struct {
	level int
	baseLogger *log.Logger
	baseFile *os.File
}

const (
	color_red = uint8(iota + 91)
	color_green		//	绿
	color_yellow		//	黄
	color_blue			// 	蓝
	color_magenta 		//	洋红
)


func New(strLevel ,pathname string, flag int) (*Logger, error){
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level= debugLevel
	case "release":
		level = releaseLevel
	case "warn":
		level = warnLevel
	case "error":
		level = errorLevel
	case "fatal":
		level= fatalLevel
	default:
		return nil, errors.New("unknow level: " + strLevel)
	}
	var baseLogger *log.Logger
	var baseFile *os.File
	if pathname != "" {
		now := time.Now()
		filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())
		file, err := os.Create(path.Join(pathname, filename))
		if err != nil {
			return nil, err
		}
		baseLogger = log.New(file, "", flag)
		baseFile = file
	} else {
		baseLogger = log.New(os.Stdout, "", flag)
	}

	logger := new(Logger)
	logger.level = level
	logger.baseLogger = baseLogger
	logger.baseFile = baseFile
	return logger, nil
}

func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}
	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level{
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}
	format = printLevel + format
	str :=fmt.Sprintf(format, a...)
	switch printLevel {
	case printDebugLevel:
		logger.baseLogger.Output(3, green(str))
	case printWranLevel:
		logger.baseLogger.Output(3, blue(str))
	case printFatalLevel:
		logger.baseLogger.Output(3, magenta(str))
	case printErrorLevel:
		logger.baseLogger.Output(3, red(str))
	case printReleaseLevel:
		logger.baseLogger.Output(3, yellow(str))
	default:
		fmt.Printf("[printlevel error]")
		return
	}

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) Debug(format string ,a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Release(format string ,a ...interface{}) {
	logger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func (logger *Logger) Wran(format string ,a ...interface{}) {
	logger.doPrintf(warnLevel, printWranLevel, format, a...)
}

func (logger *Logger) Error(format string ,a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string ,a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

var gLogger ,_ = New("debug", "", log.LstdFlags)

func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}

func Debug(format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func Warn(format string, a ...interface{}) {
	gLogger.doPrintf(warnLevel, printWranLevel, format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func Close() {
	gLogger.Close()
}

func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_green, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_yellow, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, s)
}

func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, s)
}
