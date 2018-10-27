package tempwriter

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func New(prefix string, wipe bool) *TempWriter {
	return &TempWriter{
		prefix: prefix,
		files:  map[string]*os.File{},
		wipe:   wipe,
	}
}

type TempWriter struct {
	prefix string
	files  map[string]*os.File
	wipe   bool

	clientNames []string
}

func (w *TempWriter) NewFile(filename string) (io.WriteCloser, error) {
	dir := path.Join(os.TempDir(), w.prefix)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	f, err := ioutil.TempFile(dir, filename)
	if err != nil {
		return nil, err
	}
	w.files[filename] = f

	return f, nil
}

func (w *TempWriter) Flush() error {
	if w.wipe {
		os.RemoveAll(w.prefix)
	}
	if err := os.MkdirAll(w.prefix, 0777); err != nil {
		panic(err)
	}

	for filename, f := range w.files {
		if err := os.Rename(
			f.Name(), // in tempdir
			filepath.Join(w.prefix, filename), // in target dir
		); err != nil {
			return err
		}
	}
	return nil
}
