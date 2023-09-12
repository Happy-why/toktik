package test

import (
	"fmt"
	"github.com/h2non/filetype"
	"io/ioutil"
	"testing"
)

func TestFileType(t *testing.T) {
	buf, _ := ioutil.ReadFile("E:\\Go\\goproject\\src\\my_project\\toktik\\toktik-video\\storage\\videos\\兰亭序.mp4")

	if filetype.IsVideo(buf) {
		fmt.Println("File is an image")
	} else {
		fmt.Println("Not an image")
	}
}
