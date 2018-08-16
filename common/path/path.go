// Package path implements utility routines for manipulating slash-separated paths.
package path

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	fileMode = 0644
	dirMode  = 0755

	jsonIdent  = "    "
	jsonPrefix = ""
)

var (
	// BufferSize is the size of the buffer. Default: one megabyte.
	BufferSize = 1048576
)

// IsDirWriteable checks if dir is writable by writing and removing a file
// to dir. It returns nil if dir is writable.
func IsDirWriteable(dir string) error {
	f := filepath.Join(dir, ".dummyFile")
	if err := ioutil.WriteFile(f, []byte(""), fileMode); err != nil {
		return err
	}
	return os.Remove(f)
}

// IsFile if true then object is file.
func IsFile(name string) bool {
	stat, err := os.Stat(name)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

// ReadDir returns the file names in the given directory.
func ReadDir(dirpath string) ([]string, error) {
	info, err := os.Stat(dirpath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s - is not a directory", dirpath)
	}

	dir, err := os.Open(dirpath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	return names, nil
}

// MakeDir is similar to os.MkdirAll. It creates directories
// with 0755 permission if any directory does not exists.
func MakeDir(dir string) error {
	if err := os.MkdirAll(dir, dirMode); err != nil {
		return err
	}
	return IsDirWriteable(dir)
}

// ReadFile reads the file named by filename and returns the contents.
func ReadFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

// WriteFile writes data to a file named by filename.
func WriteFile(file string, data []byte) error {
	return ioutil.WriteFile(file, data, fileMode)
}

func openSrc(name string) (*os.File, error) {
	sourceFileStat, err := os.Stat(name)
	if err != nil {
		return nil, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", name)
	}

	source, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func createDst(name string) (*os.File, error) {
	_, err := os.Stat(name)
	if err == nil {
		return nil, fmt.Errorf("file %s already exists", name)
	}

	return os.Create(name)
}

// CopyFile copies source file to destination.
func CopyFile(src, dst string, bufferSize int) error {
	source, err := openSrc(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := createDst(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, bufferSize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// ReadJSONFile reads and parses a JSON file filling a given data instance.
func ReadJSONFile(name string, data interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(data)
}

// WriteJSONFile converts a given data instance to JSON and writes it to file.
func WriteJSONFile(name string, data interface{}) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent(jsonPrefix, jsonIdent)
	return enc.Encode(data)
}
