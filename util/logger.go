package util

import (
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger // Just about anything
	Info    *log.Logger // Important information
	Warning *log.Logger // Be concerned
	Error   *log.Logger // Critical problem
)

func init() {
	file, err := os.OpenFile("git-sync.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Trace = log.New(io.MultiWriter(file, os.Stdout),
		"TRACE: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Info = log.New(io.MultiWriter(file, os.Stdout),
		"INFO: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Warning = log.New(io.MultiWriter(file, os.Stderr),
		"WARNING: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Error = log.New(io.MultiWriter(file, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)
}
