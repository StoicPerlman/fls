# File Line Seeker

[![Build Status](https://travis-ci.org/StoicPerlman/fls.svg?branch=master)](https://travis-ci.org/StoicPerlman/fls)

## Usage
### Example

```go
import "github.com/stoicperlman/fls"

f, _ := os.OpenFile("test.log", os.O_CREATE|os.O_RDONLY, 0600)
defer f.Close()
file := fls.LineFile(f)

_, err = file.SeekLine(-10, io.SeekEnd)
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
