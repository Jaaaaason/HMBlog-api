package logger

import (
	"log"
	"os"

	"github.com/jaaaaason/hmblog/configer"
)

var infoLog *log.Logger
var warnLog *log.Logger
var errLog *log.Logger
var fatalLog *log.Logger

func init() {
	infoLog = log.New(os.Stdout, "[Info] ", log.Ldate|log.Ltime)
	warnLog = log.New(os.Stdout, "[Warn] ", log.Ldate|log.Ltime)
	errLog = log.New(os.Stdout, "[Error] ", log.Ldate|log.Ltime)
	fatalLog = log.New(os.Stdout, "[Fatal] ", log.Ldate|log.Ltime)
}

// SetOutputFile sets the destination file of log's output
func SetOutputFile(filepath string) error {
	file, err := os.OpenFile(configer.Config.LogFile, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	infoLog.SetOutput(file)
	warnLog.SetOutput(file)
	errLog.SetOutput(file)
	fatalLog.SetOutput(file)

	return nil
}

// Info prints the info log with given string
func Info(s string) {
	infoLog.Println(s)
}

// Warn prints the warn log with given string
func Warn(s string) {
	warnLog.Println(s)
}

// Error prints the error log with given string
func Error(s string) {
	errLog.Println(s)
}

// Fatal prints the fatal log with given string
// and then invoke os.Exit(1)
func Fatal(s string) {
	fatalLog.Fatal(s)
}
