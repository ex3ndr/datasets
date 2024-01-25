package utils

import (
	"io"

	"github.com/schollz/progressbar/v3"
)

type ProgressReader struct {
	r   io.Reader
	bar *progressbar.ProgressBar
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.bar.Add(n)
	return n, err
}
