package io

import (
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
)

func ListenAndServe(w io.Writer, signal chan error) {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		signal <- err
	}

	http.HandleFunc("/", pipeHandler(w))

	signal <- nil
	log.Fatal(http.Serve(l, nil))
}

func pipeHandler(metricsWriter io.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pr, pw := io.Pipe()

		cmd := exec.Command("cat", "fruits.txt")
		cmd.Stdout = pw
		cmd.Stderr = pw

		go func() {
			mw := io.MultiWriter(metricsWriter, w)
			io.Copy(mw, pr)
		}()

		if err := cmd.Run(); err != nil {
			pw.CloseWithError(err)
		}

		pw.Close()
	}
}
