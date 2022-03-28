package logger

import (
	"awesomeProjectRentaTeam/internal/config"
	"awesomeProjectRentaTeam/pkg/erx"
	"fmt"
	logrus "github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"sync"
)

var instance *logrus.Logger
var once sync.Once

func GetLogrusInstance() *logrus.Logger {
	once.Do(func() {
		instance = Logrus()
	})
	return instance
}

type PlainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf(entry.Time.Format(f.TimestampFormat))
	return []byte(fmt.Sprintf("[%s] [%s] %s\n", timestamp, f.LevelDesc[entry.Level], entry.Message)), nil
}

func Logrus() *logrus.Logger {
	logger := logrus.New()
	cfg := config.GetConfigInstance()
	path := cfg.Logger.Path

	if path != "" {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Println(erx.New(err))
		}
		defer func() {
			err = file.Close()
			if err != nil {
				log.Println(erx.New(err))
			}
		}()
		mw := io.MultiWriter(os.Stdout, file)
		if err == nil {
			log.SetOutput(mw)
		}
	}
	plainFormatter := new(PlainFormatter)
	plainFormatter.TimestampFormat = "2006-01-02 15:04:05 MST"
	plainFormatter.LevelDesc = []string{"PANIC", "FATAL", "ERROR", "WARNING", "INFO", "DEBUG"}
	logger.SetFormatter(plainFormatter)

	return logger
}
