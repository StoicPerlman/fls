package fls

import (
	"io"
	"math"
	"os"
)

type File struct {
	*os.File
}

const BufferLength = 32 * 1024

const (
	SeekStart int = io.SeekStart     // seek relative to the origin of the file
	SeekCurrent int = io.SeekCurrent // seek relative to the current offset
	SeekEnd int = io.SeekEnd         // seek relative to the end

	O_RDONLY int = os.O_RDONLY // open the file read-only.
	O_WRONLY int = os.O_WRONLY // open the file write-only.
	O_RDWR   int = os.O_RDWR   // open the file read-write.
	O_APPEND int = os.O_APPEND // append data to the file when writing.
	O_CREATE int = os.O_CREAT  // create a new file if none exists.
	O_EXCL   int = os.O_EXCL   // used with O_CREATE, file must not exist
	O_SYNC   int = os.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = os.O_TRUNC  // if possible, truncate file when opened.
)

var EOF = io.EOF

func LineFile(file *os.File) *File {
	return &File{file}
}

// seeks through file by line
// positive lines will move forward in file
// 0 lines will move to begining of line
// negative lines will move backwards in file

// -1 lines, 0 whence will return EOF
// 0 lines, 0 whence will return begining line 1
// 1 lines, 0 whence will return begining line 2

// -1 lines, 1 whence will return begining of previous line
// 0 lines, 1 whence will return begining of current line
// 1 lines, 1 whence will return begining of next line

// -1 lines, 2 whence will return begining of second to last line
// 0 lines, 2 whence will return begining of last line
// 1 lines, 2 whence will return EOF
func (file *File) SeekLine(lines int64, whence int) (int64, error) {

	// return error on bad whence
	if whence < 0 || whence > 2 {
		return file.Seek(0, whence)
	}

	position, err := file.Seek(0, whence)

	buf := make([]byte, BufferLength)
	bufLen := 0
	lineSep := byte('\n')
	seekBack := lines < 1
	lines = int64(math.Abs(float64(lines)))
	matchCount := int64(0)

	// seekBack ignores first match
	// allows 0 to go to begining of current line
	if seekBack {
		matchCount = -1
	}

	leftPosition := position
	offset := int64(BufferLength * -1)

	for b := 1; ; b++ {
		if err != nil {
			break
		}

		if seekBack {

			// on seekBack 2nd buffer onward needs to seek
			// past what was just read plus another buffer size
			if b == 2 {
				offset *= 2
			}

			// if next seekBack will pass beginning of file
			// buffer is 0 to unread position
			if position+int64(offset) <= 0 {
				buf = make([]byte, leftPosition)
				position, err = file.Seek(0, io.SeekStart)
				leftPosition = 0
			} else {
				position, err = file.Seek(offset, io.SeekCurrent)
				leftPosition = leftPosition - BufferLength
			}
		}
		if err != nil {
			break
		}

		bufLen, err = file.Read(buf)
		if err != nil {
			break
		} else if seekBack && leftPosition == 0 {
			err = io.EOF
		}

		for i := 0; i < bufLen; i++ {
			iToCheck := i
			if seekBack {
				iToCheck = bufLen - i - 1
			}
			byteToCheck := buf[iToCheck]

			if byteToCheck == lineSep {
				matchCount++
			}

			if matchCount == lines {
				if seekBack {
					return file.Seek(int64(i)*-1, io.SeekCurrent)
				}
				return file.Seek(int64(bufLen*-1+i+1), io.SeekCurrent)
			}
		}
	}

	if err == io.EOF && !seekBack {
		position, _ = file.Seek(0, io.SeekEnd)
	} else if err == io.EOF && seekBack {
		position, _ = file.Seek(0, io.SeekStart)

		// no EOF err on SeekLine(0,0)
		if lines == 0 {
			err = nil
		}
	}

	return position, err
}

// os file wrappers
func Create(name string) (*File, error) {
	f, err := os.Create(name)
	file := LineFile(f)
	return file, err
}

func NewFile(fd uintptr, name string) *File {
	f := os.NewFile(fd, name)
	file := LineFile(f)
	return file
}

func Open(name string) (*File, error) {
	f, err := os.Open(name)
	file := LineFile(f)
	return file, err
}

func OpenFile(name string, flag int, perm os.FileMode) (*File, error) {
	f, err := os.OpenFile(name, flag, perm)
	file := LineFile(f)
	return file, err
}

func Pipe() (r *File, w *File, err error) {
	f1, f2, err := os.Pipe()
	file1 := LineFile(f1)
	file2 := LineFile(f2)
	return file1, file2, err
}
