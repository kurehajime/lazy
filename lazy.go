// lazy
package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	UTF8      = 65001
	SHIFT_JIS = 932
)

func main() {
	var text string
	var err error
	text, ok := readPipe()
	if ok == true {
		encode := getConsoleLang()
		switch encode {
		case SHIFT_JIS:
			text, err = utf2sjis(text)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	} else {
		text, ok = readFileByArg()
		if ok == false {
			os.Exit(1)
		}
	}

	lazyPrint(text)
}

func readPipe() (string, bool) {
	stats, _ := os.Stdin.Stat()

	if stats != nil && (stats.Mode()&os.ModeCharDevice) == 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err.Error())
			return "", false
		}

		return string(bytes), true
	} else {
		return "", false
	}
}
func readFileByArg() (string, bool) {

	if len(os.Args) < 2 {

		return "", false
	}
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}

	return string(content), true
}
func lazyPrint(text string) {

	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		runes := []rune(lines[i])
		for j := 0; j < len(runes); j += 1 {
			os.Stdout.Write([]byte(string(runes[j])))
			time.Sleep(5 * time.Millisecond)
			if string(runes[j]) == " " {
				time.Sleep(20 * time.Millisecond)
			}
		}
		os.Stdout.Write([]byte("\n"))

		time.Sleep(30 * time.Millisecond)
	}
}

//for Windows
func getConsoleLang() int {
	if runtime.GOOS != "windows" {
		return UTF8 //default:utf-8
	}
	out, err := exec.Command("chcp").Output()
	if err != nil {
		fmt.Println(err)
		return UTF8 //default:utf-8
	}
	rtn, err := strconv.Atoi(string(regexp.MustCompile(`(\d+)`).FindSubmatch([]byte(string(out)))[0]))
	if err != nil {
		fmt.Println(err)
		return UTF8 //default:utf-8
	}
	return rtn
}

func utf2sjis(text string) (string, error) {
	str, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(text)), japanese.ShiftJIS.NewDecoder()))
	if err == nil {
		return string(str), nil
	} else {
		return "", err
	}
}
