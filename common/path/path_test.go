package path

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pborman/uuid"
)

func TestCopyFile(t *testing.T) {
	dir, err := ioutil.TempDir("", uuid.NewUUID().String())
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	file1 := filepath.Join(dir, uuid.NewUUID().String())
	file2 := filepath.Join(dir, "file2")

	f, err := os.Create(file1)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100000; i++ {
		f.Write([]byte(uuid.NewUUID().String() + "\n"))
	}

	f.Close()

	if err := CopyFile(file1, file2, BufferSize); err != nil {
		t.Fatal(err)
	}

	file1Data, err := ioutil.ReadFile(file1)
	if err != nil {
		t.Fatal(err)
	}

	file2Data, err := ioutil.ReadFile(file2)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(file1Data, file2Data) {
		t.Fatal("files not equal")
	}
}
