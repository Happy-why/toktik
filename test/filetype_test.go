package test

import (
	"fmt"
	"github.com/h2non/filetype"
	"io/ioutil"
	"testing"
)

func TestFileType(t *testing.T) {
	buf, _ := ioutil.ReadFile("E:\\Go\\goproject\\src\\my_project\\toktik\\toktik-video\\storage\\videos\\兰亭序.mp4")

	kind, _ := filetype.Match(buf)
	if kind == filetype.Unknown {
		fmt.Println("Unknown file type")
		return
	}

	fmt.Printf("File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
}
