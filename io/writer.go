package io

import (
	"bufio"
	"bytes"
)

type MetricsWriter struct {
	data  []byte
	lines int
}

func (mr *MetricsWriter) Write(b []byte) (int, error) {
	mr.data = append(mr.data, b...)
	return len(b), nil
}

func (mr MetricsWriter) Lines() int {
	r := bytes.NewReader(mr.data)
	br := bufio.NewReader(r)

	var lines int

	for {
		_, _, err := br.ReadLine()
		if err != nil {
			break
		}

		lines++
	}

	return lines
}
