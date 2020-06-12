package runtime

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

type zipper struct {
	srcDir string
}

func newZipper(srcDir string) *zipper {
	return &zipper{
		srcDir: srcDir,
	}
}

func (z *zipper) zipp() ([]byte, error) {
	files := make(map[string][]byte)

	// go through the dir and read all the files
	err := filepath.Walk(z.srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == z.srcDir {
			return nil
		}

		// skip directory
		if info.IsDir() {
			return filepath.SkipDir
		}
		f, e := os.Open(path)
		if e != nil {
			return e
		}
		defer f.Close()
		contents, e := ioutil.ReadAll(f)
		if e != nil {
			return e
		}
		// split returns dir, file name
		// we just need the file name
		_, path = filepath.Split(path)
		files[path] = contents
		return nil
	})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	f, err := w.Create("_entry.py")
	if err != nil {
		return nil, err
	}

	_, err = f.Write(z.entryFileContent())
	if err != nil {
		return nil, err
	}

	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			return nil, err
		}
		_, err = f.Write(content)
		if err != nil {
			return nil, err
		}
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (z *zipper) entryFileContent() []byte {
	return []byte(`from deta.lib.handler import handle
from deta.lib.debug import debug

try:
	import main
except:
	pass
	
@debug
def handler(event, context):
	import main
	return handle(event, main)`)
}