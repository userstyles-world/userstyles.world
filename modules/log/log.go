package log

import (
	"io"
	"log"
	"os"
	"path"

	"userstyles.world/modules/config"
)

var (
	Info *log.Logger
	Warn *log.Logger
)

func setOutput(f *os.File) io.Writer {
	if config.Production {
		return io.MultiWriter(f)
	}

	return io.MultiWriter(os.Stdout, f)
}

func Initialize() {
	f, err := os.OpenFile(path.Join(config.DataDir, "server.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatalf("Failed to open log file, err: %v\n", err)
	}

	// Configure output.
	mw := setOutput(f)

	// Initialize loggers.
	Info = log.New(mw, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(mw, "WARN ", log.Ldate|log.Ltime|log.Lshortfile)
}
