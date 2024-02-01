package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"log/slog"

	consts "github.com/michelemendel/dmtmms/constants"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	InitEnv()
}

// https://betterstack.com/community/guides/logging/logging-in-go/
const (
	LevelTrace  = slog.Level(-8)
	LevelNotice = slog.Level(2)
	LevelFatal  = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace:  "TRACE",
	LevelNotice: "NOTICE",
	LevelFatal:  "FATAL",
}

func StdOutLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, options))
}

func FileLogger() *slog.Logger {
	// file := openfile()
	file := rotate(filepath.Join("log", os.Getenv(consts.ENV_FILE_NAME_KEY)))
	return slog.New(slog.NewTextHandler(file, options))
}

// func openfile() *os.File {
// 	filename := filepath.Join("log", os.Getenv(consts.LOG_FILE_NAME_ENV))
// 	fmt.Println("[LOG]:openfile", "LOG_FILE_NAME_ENV:", filename)
// 	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return f
// }

func rotate(path string) io.WriteCloser {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    10, // In MB before rotating the file
		MaxAge:     30, // In days before deleting the file
		MaxBackups: 5,  // Maximum number of backups to keep track of
	}
}

var options = &slog.HandlerOptions{
	AddSource:   true,
	Level:       slog.LevelDebug,
	ReplaceAttr: MakeReplaceAttr(),
}

func MakeReplaceAttr() func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			level := a.Value.Any().(slog.Level)
			levelLabel, exists := LevelNames[level]
			if !exists {
				levelLabel = level.String()
			}
			a.Value = slog.StringValue(levelLabel)
		}

		if a.Key == slog.SourceKey {
			filename := a.Value.Any().(*slog.Source).File
			lineNumber := a.Value.Any().(*slog.Source).Line
			ps := strings.Split(filename, "/")
			a.Value = slog.StringValue(fmt.Sprintf("%s/%s:%v", ps[len(ps)-2], ps[len(ps)-1], lineNumber))
		}

		return a
	}
}
