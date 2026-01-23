package generate

import (
	"bytes"
	"go/format"
	"os"
	"text/template"
)

// IsExist returns whether a file or directory exists.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func MkdirAll(path string) {
	err := os.MkdirAll(path, 0755)
	MustCheck(err)
}
func WriteToFile(filename, content string) {
	//fmt.Println(content)
	//return
	f, err := os.Create(filename)
	MustCheck(err)
	defer CloseFile(f)
	_, err = f.WriteString(content)
	MustCheck(err)
}

// MustCheck panics when the error is not nil
func MustCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func CloseFile(f *os.File) {
	err := f.Close()
	MustCheck(err)
}

/*
*
needFormat 是否需要格式化，仅限于.go文件
*/
func ReadTmpl(tempName, tempText string, data interface{}, needFormat bool) string {
	codeBuf := new(bytes.Buffer)
	tmpl, err := template.New(tempName).Parse(tempText)
	MustCheck(err)
	err = tmpl.Execute(codeBuf, data)
	MustCheck(err)
	if needFormat {
		codeBytes, err := format.Source(codeBuf.Bytes())
		MustCheck(err)
		return string(codeBytes)
	}
	return string(codeBuf.Bytes())
}
