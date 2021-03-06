package gphotos

import (
	"io"
	"os"
	"path"
)

// FileUploadItem represents a local file.
type FileUploadItem string

// Open returns a stream.
// Caller should close it finally.
func (m FileUploadItem) Open() (io.ReadSeeker, int64, error) {
	r, err := os.Open(m.path())
	if err != nil {
		return nil, 0, err
	}
	fi, err := r.Stat()
	if err != nil {
		return nil, 0, err
	}
	return r, fi.Size(), nil
}

// Name returns the filename.
func (m FileUploadItem) Name() string {
	f, err := os.Stat(m.path())
	if err != nil {
		return ""
	}
	return path.Base(f.Name())
}

func (m FileUploadItem) path() string {
	return string(m)
}

// Size returns size of the file.
func (m FileUploadItem) Size() int64 {
	f, err := os.Stat(m.path())
	if err != nil {
		return 0
	}
	return f.Size()
}
