package clients

import (
	"bufio"
	"errors"
	"io"
)

const logReaderPrefixLen = 1000

type logReader struct {
	bufferedReader *bufio.Reader
	reader         io.ReadCloser // reader provided by the client
}

func newLogReader(reader io.ReadCloser) *logReader {
	return &logReader{
		reader:         reader,
		bufferedReader: bufio.NewReader(reader),
	}
}

func (r *logReader) NextLine() ([]byte, error) {
	line, isPrefix, err := r.bufferedReader.ReadLine()
	if !isPrefix || err != nil {
		return line, err
	}
	prefix := make([]byte, logReaderPrefixLen)
	for i := 0; isPrefix; i++ {
		// this loop is entered if a log line is too long to fit into the buffer. We discard it by
		// iterating until isPrefix becomes false. We only log the first few bytes of the line to help with
		// identification.
		if i == 0 {
			prefixLen := logReaderPrefixLen
			if len(line) < prefixLen {
				prefixLen = len(line)
			}
			copy(prefix, line[:prefixLen])
		}
		line, isPrefix, err = r.bufferedReader.ReadLine()
	}
	return prefix, errors.New("log line too long, discarding")
}
