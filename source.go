package logtee

import (
	"bufio"
	"github.com/hpcloud/tail"
	"io"
	"os"
)

type Source <-chan string

func fileSource(f *os.File) Source {
	reader := bufio.NewReader(f)
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

func StdinSource() Source {
	return fileSource(os.Stdin)
}

func FileSource(filename string, follow bool) (Source, error) {
	if follow {
		t, err := tail.TailFile(filename, tail.Config{
			Follow: true,
			Location: &tail.SeekInfo{
				Whence: os.SEEK_END,
			},
			Logger: tail.DiscardingLogger,
		})
		if err != nil {
			return nil, err
		}
		s := make(chan string)
		go func() {
			for line := range t.Lines {
				s <- line.Text
			}
		}()
		return Source(s), nil
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		return fileSource(f), nil
	}
}
