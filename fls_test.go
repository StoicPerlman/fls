package fls

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

type T struct {
	*testing.T
}

const TestFileEOFPos = 48889

func init() {
	f, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for i := 0; i <= 9999; i++ {

		lineNum := strconv.Itoa(i)

		if i > 0 {
			lineNum = "\n" + lineNum
		}

		if _, err = f.WriteString(lineNum); err != nil {
			panic(err)
		}
	}
}

func TestSeekLineStart(t *testing.T) {
	myT := &T{t}

	f, err := os.OpenFile("test.log", os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	file := LineFile(f)

	_, err = file.SeekLine(0, io.SeekStart)
	line := GetLine(file)

	myT.Ok(line, 0, err, false)

	_, err = file.SeekLine(1, io.SeekStart)
	line = GetLine(file)

	myT.Ok(line, 1, err, false)

	_, err = file.SeekLine(100, io.SeekStart)
	line = GetLine(file)

	myT.Ok(line, 100, err, false)

	_, err = file.SeekLine(10000, io.SeekStart)
	line = GetLine(file)

	myT.Ok(line, 9999, err, true)

	_, err = file.SeekLine(-1, io.SeekStart)
	line = GetLine(file)

	myT.Ok(line, 0, err, true)
}

func TestSeekLineEnd(t *testing.T) {
	myT := &T{t}

	f, err := os.OpenFile("test.log", os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	file := LineFile(f)

	pos, _ := file.Seek(0, io.SeekEnd)

	// Test const TestFileEOFPos = 48889
	if pos != TestFileEOFPos {
		t.Error("\nEOF hit at unknown position")
	}

	_, err = file.SeekLine(1, io.SeekEnd)
	line := GetLine(file)

	myT.Ok(line, 9999, err, true)

	_, err = file.SeekLine(0, io.SeekEnd)
	line = GetLine(file)

	myT.Ok(line, 9999, err, false)

	_, err = file.SeekLine(-1, io.SeekEnd)
	line = GetLine(file)

	myT.Ok(line, 9998, err, false)

	_, err = file.SeekLine(-100, io.SeekEnd)
	line = GetLine(file)

	myT.Ok(line, 9899, err, false)
}

// Test helper functions
func GetLine(file *File) int {
	pos, _ := file.Seek(0, io.SeekCurrent)

	if pos == TestFileEOFPos {
		return 9999
	}

	buf := make([]byte, 100)
	file.Read(buf)

	// strips null chars from end of buffer on EOF
	buf = bytes.Trim(buf, "\x00")

	s := string(buf)
	sp := strings.Split(s, "\n")

	// resets pos to pos before read
	file.Seek(pos, io.SeekStart)

	line, _ := strconv.Atoi(sp[0])
	return line
}

func (t *T) Ok(got int, expected int, err error, expectEOF bool) {
	if got != expected {
		t.Error("\nExpected line: ", expected, "\ngot: ", got)
	}

	if err != nil && err != io.EOF {
		t.Error("\nError: ", err)
	} else if expectEOF && err != io.EOF {
		t.Error("\nExpected to hit EOF")
	} else if !expectEOF && err == io.EOF {
		t.Error("\nDid not expect to hit EOF")
	}
}
