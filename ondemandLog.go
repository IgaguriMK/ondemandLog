package ondemandLog

import (
	"io"
	"log"
	"os"
	"time"
)

const (
	TimeFormat = "2006-01-02_15-04-05"
)

func TimeLog(format string, flag int, lazy bool) {
	name := time.Now().Format(format)

	if lazy {
		LazyLog(name, flag, false)
	} else {
		StrictLog(name, flag, false)
	}
}

func StrictLog(name string, flag int, append bool) {
	var file *os.File
	var err error
	if append {
		file, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(name)
	}
	if err != nil {
		panic(err.Error())
	}

	logw := io.MultiWriter(file, os.Stderr)
	log.SetOutput(logw)
	log.SetFlags(flag)
}

func LazyLog(name string, flag int, append bool) {
	w := NewLazyWriter(name, append)
	logw := io.MultiWriter(w, os.Stderr)
	log.SetOutput(logw)
	log.SetFlags(flag)
}

type LazyFileWriter struct {
	name   string
	append bool
	file   *os.File
}

func NewLazyWriter(name string, append bool) *LazyFileWriter {
	return &LazyFileWriter{
		name:   name,
		append: append,
	}
}

func (lw *LazyFileWriter) Write(p []byte) (n int, err error) {
	n = 0

	if lw.file == nil {
		if lw.append {
			lw.file, err = os.OpenFile(lw.name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
		} else {
			lw.file, err = os.Create(lw.name)
			if err != nil {
				return
			}
		}
	}

	n, err = lw.file.Write(p)
	return
}

func (lw *LazyFileWriter) Close() error {
	if lw.file != nil {
		return lw.file.Close()
	}

	return nil
}
