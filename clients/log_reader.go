package clients

import (
	"bufio"
	"errors"
	"io"
)

// logReaderPrefixLen is used when returning a partial line as context in NextLine
const logReaderPrefixLen = 1000

var errLogLineToLong = errors.New("log line too long, discarding")

// logReader is a custom implementation similar to bufio.Scanner, but provides a way to handle lines
// (or tokens) that exceed the buffer size.
type logReader struct {
	bufferedReader *bufio.Reader
	reader         io.ReadCloser // reader provided by the client
}

// newLogReader creates a new logReader to read log lines from an io.ReadCloser
func newLogReader(reader io.ReadCloser) *logReader {
	return &logReader{
		reader:         reader,
		bufferedReader: bufio.NewReader(reader),
	}
}

// NextLine reads and returns the next log line from the reader. An io.EOF error is returned
// if the end of the stream has been reached. This implementation is different from bufio.Scanner as it
// also returns an error if a line is too long to fit into the buffer. In this case, an error is returned
// together with a limited prefix of the line.
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
		if err != nil {
			return nil, err
		}
	}
	return prefix, errLogLineToLong
}
