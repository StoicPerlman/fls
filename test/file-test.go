package main

import "fmt"
import "github.com/StoicPerlman/fls"
import "os"
import "strings"

type SeekInfo struct {
	Offset int64
	Whence int // os.SEEK_*
}

const BufferLength = 32 * 1024

func main() {
	f, err := os.Open("test.log")
	file := fls.LineFile(f)

	check(err)

	pos, _ := file.SeekLine(0,0)
	printLn(file, pos)

	pos, _ = file.SeekLine(1,0)
	printLn(file, pos)

	pos, _ = file.SeekLine(3,0)
	printLn(file, pos)

	pos, _ = file.SeekLine(3,1)
	printLn(file, pos)

	pos, _ = file.SeekLine(3,1)
	printLn(file, pos)

	pos, _ = file.SeekLine(-3,1)
	printLn(file, pos)

	pos, _ = file.SeekLine(100,1)
	printLn(file, pos)

	pos, _ = file.SeekLine(0,2)
	printLn(file, pos)

	pos, _ = file.SeekLine(-1,2)
	printLn(file, pos)

	pos, _ = file.SeekLine(-1,1)
	printLn(file, pos)

	pos, _ = file.SeekLine(-1,1)
	printLn(file, pos)

	pos, _ = file.SeekLine(-100,1)
	printLn(file, pos)

	f.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printLn(file *fls.File, pos int64) {
	file.Seek(pos, 0)

	buf := make([]byte, 100)
	file.Read(buf)

	s := string(buf)
	sp := strings.Split(s, "\n")
	fmt.Println(sp[0])

	file.Seek(pos, 0)
}
