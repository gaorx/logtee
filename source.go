package logtee

import (
	"bufio"
	"errors"
	"github.com/hpcloud/tail"
	"io"
	"os"
)

type Source <-chan string

func StdinSource() Source {
	reader := bufio.NewReader(os.Stdin)
	s := make(chan string)
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				close(s)
				return
			}
			if err != nil {
				// TODO
				continue
			}
			s <- line
		}
	}()
	return Source(s)
}

func FileSource(filename string, follow bool) (Source, error) {
	if follow {
		t, err := tail.TailFile(filename, tail.Config{Follow: true})
		if err != nil {
			return nil, err
		}
		s := make(chan string)
		for line := range t.Lines {
			s <- line.Text
		}
		return Source(s), nil
	} else {
		return nil, errors.New("Not support")
	}
}
