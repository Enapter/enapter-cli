package enaptercli

import (
	"io"
	"sync"
)

type onceWriter struct {
	once sync.Once
	w    io.Writer
}

func (w *onceWriter) Write(p []byte) (int, error) {
	n := len(p)
	var err error
	w.once.Do(func() {
		n, err = w.w.Write(p)
	})
	return n, err
}
