package logger

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

type Logger struct {
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
}

func NewLogger() *Logger {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	l := new(Logger)

	l.InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return l
}

func (l *Logger) Fatalf(format string, v ...any) {
	log.Printf(format, v)
	l.ErrorLogger.Fatalf(format, v)
}

func (l *Logger) EPrintf(format string, v ...any) {
	log.Printf(format, v)
	l.ErrorLogger.Printf(format, v)
}

func (l *Logger) IPrintf(format string, v ...any) {
	log.Printf(format, v)
	l.InfoLogger.Printf(format, v)
}

func (l *Logger) IPrintln(v ...any) {
	log.Println(v)
	l.InfoLogger.Println(v)
}
