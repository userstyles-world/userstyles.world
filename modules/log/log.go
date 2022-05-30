package log

import (
	"io"
	"log"
	"os"

	"userstyles.world/modules/config"
)

var (
	Info *log.Logger
	Warn *log.Logger

	// Database logger will emit output from our DBMS.
	Database *log.Logger
)

func setOutput(f *os.File) io.Writer {
	if config.Production {
		return io.MultiWriter(f)
	}

	return io.MultiWriter(os.Stdout, f)
}

func Initialize() {
	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(config.LogFile, flags, 0o666)
	if err != nil {
		log.Fatalf("Failed to open %v: %s\n", config.LogFile, err)
	}

	// Configure output.
	mw := setOutput(f)

	// Initialize loggers.
	Info = log.New(mw, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(mw, "WARN ", log.Ldate|log.Ltime|log.Lshortfile)
	Database = log.New(mw, "DBMS ", log.Ldate|log.Ltime)
}
