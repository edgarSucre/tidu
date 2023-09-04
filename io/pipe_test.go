package io_test

import (
	"bufio"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	mio "github.com/edgarsucre/tidu/io"

	"github.com/stretchr/testify/assert"
)

func Test_Pipe_Reader_Blocks(t *testing.T) {
	pr, pw := io.Pipe()
	sr := strings.NewReader("test")

	go func() {
		defer pw.Close()

		time.Sleep(time.Millisecond * 200)
		io.Copy(pw, sr)
	}()

	buf := make([]byte, 4)

	start := time.Now()
	pr.Read(buf)

	end := time.Since(start).Truncate(time.Millisecond)
	assert.Equal(t, 200*time.Millisecond, end)
}

func Test_Pipe_Reader_With_Error(t *testing.T) {
	pr, pw := io.Pipe()
	br := new(mio.BadReader)

	go func() {
		defer pw.Close()

		_, err := io.Copy(pw, br)
		if err != nil {
			pw.CloseWithError(err)
		}
	}()

	_, err := io.Copy(io.Discard, pr)
	assert.ErrorContains(t, err, "reader error")
}

func Test_Pipe_Example(t *testing.T) {
	mw := new(mio.MetricsWriter)
	signal := make(chan error)

	go mio.ListenAndServe(mw, signal)

	err := <-signal
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get("http://localhost:8000/")
	assert.NoError(t, err)

	r := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	var lines int
	for {
		l, _, err := r.ReadLine()
		if err != nil {
			break
		}
		_ = l

		lines++
	}

	assert.Equal(t, lines, mw.Lines())
}
