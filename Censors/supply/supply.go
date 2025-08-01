package supply

import (
	"GoNewsAPICommentsCensors/pkg/storage"
	"io"
	"os"
	"strings"
)

func StopList() ([]storage.Stop, error) {
	f, err := os.Open("./pkg/supply/words.txt")
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
		}
	}(f)
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	var sl []storage.Stop
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		str := storage.Stop{
			StopList: trimmedLine,
		}
		sl = append(sl, str)
	}

	return sl, nil
}
