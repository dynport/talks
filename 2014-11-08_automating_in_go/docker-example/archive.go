package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
)

type File struct {
	Name    string
	Content string
}

func createArchive(files []*File) (*bytes.Buffer, error) {
	b := &bytes.Buffer{}
	gz := gzip.NewWriter(b)
	defer gz.Close()
	t := tar.NewWriter(gz)
	defer t.Close()
	for _, f := range files {
		if err := t.WriteHeader(&tar.Header{Name: f.Name, Size: int64(len(f.Content))}); err != nil {
			return nil, err
		}
		if _, err := io.WriteString(t, f.Content); err != nil {
			return nil, err
		}
	}
	return b, nil
}
