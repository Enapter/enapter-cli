package enaptercli_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/app/enaptercli"
)

var errExitTimeout = errors.New("exit timed out")

type testApp struct {
	app    *cli.App
	outBuf *lineBuffer
	errBuf *bytes.Buffer
	errCh  chan error
	cancel func()
}

func startTestApp(args ...string) *testApp {
	outBuf := newLineBuffer()
	errBuf := &bytes.Buffer{}

	app := enaptercli.NewApp()
	app.HideVersion = true
	app.Writer = outBuf
	app.ErrWriter = errBuf
	app.ExitErrHandler = func(*cli.Context, error) {}

	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		errCh <- app.RunContext(ctx, args)
	}()

	return &testApp{
		app:    app,
		outBuf: outBuf,
		errBuf: errBuf,
		errCh:  errCh,
		cancel: cancel,
	}
}

func (a *testApp) Stop() {
	a.cancel()
}

func (a *testApp) Wait() error {
	t := time.NewTimer(5 * time.Second)
	select {
	case err := <-a.errCh:
		return err
	case <-t.C:
		return errExitTimeout
	}
}

func (a *testApp) Stdout() *lineBuffer {
	return a.outBuf
}

type lineBuffer struct {
	mu   sync.Mutex
	cond *sync.Cond
	buf  *bytes.Buffer
}

func newLineBuffer() *lineBuffer {
	lb := &lineBuffer{buf: &bytes.Buffer{}}
	lb.cond = sync.NewCond(&lb.mu)
	return lb
}

func (lb *lineBuffer) Write(b []byte) (int, error) {
	lb.mu.Lock()
	defer lb.cond.Broadcast()
	defer lb.mu.Unlock()
	return lb.buf.Write(b)
}

func (lb *lineBuffer) Read(b []byte) (int, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.buf.Read(b)
}

func (lb *lineBuffer) ReadLine() (string, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	var sb strings.Builder
	for {
		b, err := lb.buf.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				lb.cond.Wait()
				continue
			}
			return "", err
		}

		sb.WriteByte(b)
		if b == '\n' {
			break
		}
	}
	return sb.String(), nil
}
