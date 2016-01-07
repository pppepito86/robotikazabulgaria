package ws

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Getwd() string {
	pwd, _ := os.Getwd()
	return pwd + string(filepath.Separator) + "work_dir"
}

func GetFilePath(relativePath ...string) string {
	var buffer bytes.Buffer
	buffer.WriteString(Getwd())
	for i := 0; i < len(relativePath)-1; i++ {
		buffer.WriteString(string(filepath.Separator))
		buffer.WriteString(relativePath[i])
	}
	os.MkdirAll(buffer.String(), os.ModePerm)
	return buffer.String() + string(filepath.Separator) + relativePath[len(relativePath)-1]
}
func ReadFile(relativePath ...string) []byte {
	fileAsByte, err := ioutil.ReadFile(GetFilePath(relativePath...))
	if err != nil {
		return []byte{}
	}
	return fileAsByte
}
