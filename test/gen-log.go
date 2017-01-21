package main
import "os"
import "strconv"
import "strings"
import "os/exec"
import "log"

func main() {
    f, err := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
	    panic(err)
	}

	defer f.Close()

	for i := 0; i <= getNumLinesToGen(); i++ {

		lineNum := strconv.Itoa(i)

		if i > 0 {
			lineNum = "\n" + lineNum
		}

		if _, err = f.WriteString(lineNum); err != nil {
		    panic(err)
		}
	}
}

func getUuid() string {
	out, err := exec.Command("uuidgen").Output()
    if err != nil {
        log.Fatal(err)
    }
    var outStr string = string(out[:])
    return strings.TrimSpace(outStr)
}

func getNumLinesToGen() int {
	args := os.Args[1:]

	if len(args) > 0 {
		i, err := strconv.Atoi(args[0])
		if err != nil {
		    panic(err)
		}
		return i
	} else {
		return 9999
	}
}