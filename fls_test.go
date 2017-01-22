package fls

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

type T struct {
	*testing.T
}

// known EOF position for test.log
const TestFileEOFPos = 588889

func init() {
	f, err := os.OpenFile("test.log", O_CREATE|O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for i := 0; i <= 99999; i++ {

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

	f, err := os.OpenFile("test.log", O_CREATE|O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	file := LineFile(f)

	_, err = file.SeekLine(-1, SeekStart)
	line := GetLine(file)
	myT.Ok(line, 0, err, true)

	_, err = file.SeekLine(0, SeekStart)
	line = GetLine(file)
	myT.Ok(line, 0, err, false)

	_, err = file.SeekLine(1, SeekStart)
	line = GetLine(file)
	myT.Ok(line, 1, err, false)

	_, err = file.SeekLine(100, SeekStart)
	line = GetLine(file)
	myT.Ok(line, 100, err, false)

	// bigger than buffer
	_, err = file.SeekLine(10000, SeekStart)
	line = GetLine(file)
	myT.Ok(line, 10000, err, false)

	_, err = file.SeekLine(50000, SeekStart)
	line = GetLine(file)
	myT.Ok(line, 50000, err, false)

	_, err = file.SeekLine(90000, SeekStart)
	line = GetLine(file)
	myT.Ok(line, 90000, err, false)

	_, err = file.SeekLine(100000, SeekStart)
	line = GetLine(file)
	myT.Ok(line, 99999, err, true)
}

func TestSeekLineEnd(t *testing.T) {
	myT := &T{t}

	f, err := os.OpenFile("test.log", os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	file := LineFile(f)

	pos, _ := file.Seek(0, SeekEnd)

	// Test const TestFileEOFPos = 588889
	if pos != TestFileEOFPos {
		t.Error("\nEOF hit at unknown position: ", pos)
	}

	_, err = file.SeekLine(1, SeekEnd)
	line := GetLine(file)
	myT.Ok(line, 99999, err, true)

	_, err = file.SeekLine(0, SeekEnd)
	line = GetLine(file)
	myT.Ok(line, 99999, err, false)

	_, err = file.SeekLine(-1, SeekEnd)
	line = GetLine(file)
	myT.Ok(line, 99998, err, false)

	_, err = file.SeekLine(-100, SeekEnd)
	line = GetLine(file)
	myT.Ok(line, 99899, err, false)

	// bigger than buffer
	_, err = file.SeekLine(-10000, SeekEnd)
	line = GetLine(file)
	myT.Ok(line, 89999, err, false)

	_, err = file.SeekLine(-50000, SeekEnd)
	line = GetLine(file)
	myT.Ok(line, 49999, err, false)

	_, err = file.SeekLine(-90000, SeekEnd)
	line = GetLine(file)
	myT.Ok(line, 9999, err, false)

	_, err = file.SeekLine(-100000, SeekEnd)
	line = GetLine(file)
	myT.Ok(line, 0, err, true)
}

func TestSeekLineCurrent(t *testing.T) {
	myT := &T{t}

	f, err := os.OpenFile("test.log", os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	file := LineFile(f)

	_, err = file.SeekLine(-1, SeekCurrent)
	line := GetLine(file)
	myT.Ok(line, 0, err, true)

	_, err = file.SeekLine(0, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 0, err, false)

	_, err = file.SeekLine(1, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 1, err, false)

	_, err = file.SeekLine(100, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 101, err, false)

	// bigger than buffer
	_, err = file.SeekLine(10000, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 10101, err, false)

	_, err = file.SeekLine(50000, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 60101, err, false)

	_, err = file.SeekLine(50000, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 99999, err, true)

	_, err = file.SeekLine(0, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 99999, err, false)

	_, err = file.SeekLine(-1, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 99998, err, false)

	_, err = file.SeekLine(-100, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 99898, err, false)

	// bigger than buffer
	_, err = file.SeekLine(-10000, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 89898, err, false)

	_, err = file.SeekLine(-50000, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 39898, err, false)

	_, err = file.SeekLine(-50000, SeekCurrent)
	line = GetLine(file)
	myT.Ok(line, 0, err, true)
}

// os file wrapper tests
func TestCreate(t *testing.T) {
	myT := &T{t}
	file, err := Create("test-create.log")
	defer file.Close()
	myT.CheckError(err)

	fi, err := file.Stat()
	myT.CheckError(err)

	if fi.Name() != "test-create.log" {
		t.Error("\nUnable to get file stats: ", fi)
	}
}

func TestNewFile(t *testing.T) {
	myT := &T{t}
	file := NewFile(uintptr(syscall.Stdin), "/dev/stdin")
	defer file.Close()

	fi, err := file.Stat()
	myT.CheckError(err)

	if fi.Name() != "stdin" {
		t.Error("\nUnable to get file stats: ", fi)
	}
}

func TestOpen(t *testing.T) {
	myT := &T{t}

	f, err := os.Create("test-open.log")
	myT.CheckError(err)
	f.Close()

	file, err := Open("test-open.log")
	myT.CheckError(err)
	defer file.Close()

	fi, err := file.Stat()
	myT.CheckError(err)

	if fi.Name() != "test-open.log" {
		t.Error("\nUnable to get file stats: ", fi)
	}
}

func TestOpenFile(t *testing.T) {
	myT := &T{t}

	file, err := OpenFile("test-open-file.log", O_CREATE|O_WRONLY, 0600)
	myT.CheckError(err)
	defer file.Close()

	fi, err := file.Stat()
	myT.CheckError(err)

	if fi.Name() != "test-open-file.log" {
		t.Error("\nUnable to get file stats: ", fi)
	}
}

func TestPipe(t *testing.T) {
	myT := &T{t}

	file1, file2, err := Pipe()
	myT.CheckError(err)
	defer file1.Close()
	defer file2.Close()

	fi1, err := file1.Stat()
	myT.CheckError(err)
	fi2, err := file2.Stat()
	myT.CheckError(err)

	if fi1.Name() != "|0" {
		t.Error("\nUnable to get file stats: ", fi1)
	}

	if fi2.Name() != "|1" {
		t.Error("\nUnable to get file stats: ", fi2)
	}
}

// Test helper functions
func GetLine(file *File) int {
	pos, _ := file.Seek(0, SeekCurrent)

	if pos == TestFileEOFPos {
		return 99999
	}

	buf := make([]byte, 100)
	file.Read(buf)

	// strips null chars from end of buffer on EOF
	buf = bytes.Trim(buf, "\x00")

	s := string(buf)
	sp := strings.Split(s, "\n")

	// resets pos to pos before read
	file.Seek(pos, SeekStart)

	line, _ := strconv.Atoi(sp[0])
	return line
}

func (t *T) Ok(got int, expected int, err error, expectEOF bool) {
	if got != expected {
		t.Error("\nExpected line: ", expected, "\ngot: ", got)
	}

	if err != nil && err != EOF {
		t.Error("\nError: ", err)
	} else if expectEOF && err != EOF {
		t.Error("\nExpected to hit EOF")
	} else if !expectEOF && err == EOF {
		t.Error("\nDid not expect to hit EOF")
	}
}

func (t *T) CheckError(err error) {
	if err != nil {
		t.Error("\nUnexpexted error: ", err)
	}
}
