package logger

import (
	"log"
	"os"
)

var (
	infoLog  *log.Logger
	warnLog  *log.Logger
	errLog   *log.Logger
	fatalLog *log.Logger
)

// Initialize initializes the package logger
func Initialize(file *os.File) {
	infoLog = log.New(file, "[Info] ", log.Ldate|log.Ltime)
	warnLog = log.New(file, "[Warn] ", log.Ldate|log.Ltime)
	errLog = log.New(file, "[Error] ", log.Ldate|log.Ltime)
	fatalLog = log.New(file, "[Fatal] ", log.Ldate|log.Ltime)
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
