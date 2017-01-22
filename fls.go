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
