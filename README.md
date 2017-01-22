# File Line Seeker

[![Build Status](https://travis-ci.org/StoicPerlman/fls.svg?branch=master)](https://travis-ci.org/StoicPerlman/fls)

## Usage
### API

https://godoc.org/github.com/StoicPerlman/fls

### Example

```go
import "github.com/stoicperlman/fls"

// use fls.LineFile(file *os.File) *File
// for files opened from os package
f, _ := os.OpenFile("test.log", os.O_CREATE|os.O_RDONLY, 0600)
defer f.Close()
file := fls.LineFile(f)

pos, err := file.SeekLine(-10, io.SeekEnd)

// use os file wrappers to open without os package
// fls.OpenFile(name string, flag int, perm os.FileMode) (*File, error)
f, err := fls.OpenFile("test.log", os.O_CREATE|os.O_WRONLY, 0600)
defer f.Close()

pos, err := file.SeekLine(-10, io.SeekEnd)
```

### Detail
`func (file *File) SeekLine(lines int64, whence int) (int64, error)`
- positive lines will move forward in file
- 0 lines will move to begining of line
- negative lines will move backwards in file

```go
file.SeekLine(-1, io.SeekStart) // return EOF
file.SeekLine(0, io.SeekStart) // return begining line 1
file.SeekLine(1, io.SeekStart) // return begining line 2

file.SeekLine(-1, io.SeekCurrent) // return begining of previous line
file.SeekLine(0, io.SeekCurrent) // return begining of current line
file.SeekLine(1, io.SeekCurrent) // return begining of next line

file.SeekLine(-1, io.SeekEnd) // return begining of second to last line
file.SeekLine(0, io.SeekEnd) // return begining of last line
file.SeekLine(1, io.SeekEnd) // return EOF
```
