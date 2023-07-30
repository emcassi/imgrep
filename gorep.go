package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

type Flags struct {
	IgnoreCase        bool
	IgnorePunctuation bool
}

func main() {

	flags, files, dirs, err := collectArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(flags)
	fmt.Println(len(files))
	fmt.Println(len(dirs))

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage("imgs/test2.png")

	text, err := client.Text()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		cleanData(&text, flags)
		fmt.Println(text)
	}
}

func newFlags() Flags {
	return Flags{
		IgnoreCase:        false,
		IgnorePunctuation: false,
	}
}

func collectArgs() (Flags, []string, []string, error) {
	flags := newFlags()
	var files []string
	var dirs []string

	for _, arg := range os.Args[1:] {

		if arg[0] == '-' {
			if len(files) > 0 {
				return flags, files, dirs, errors.New("all arguments must come before your files")
			}
			switch arg[1:] {
			case "ic":
				flags.IgnoreCase = true
			case "ip":
				flags.IgnorePunctuation = true
			default:
				return flags, files, dirs, errors.New("invalid argument: " + arg)
			}
		} else {
			if filepath.Ext(arg) == "" {
				if containsString(dirs, arg) { continue }
				dirs = append(dirs, arg)
			} else {
				if containsString(files, arg) { continue }
				files = append(files, arg)
			}
		}
	}

	return flags, files, dirs, nil
}

func cleanData(text *string, flags Flags) {

	var altered []rune
	var addChar bool

	punct := ",.!?:;'=[](){}\\|/~“”’`"

	if flags.IgnoreCase {
		*text = strings.ToLower(*text)
	}

	for _, c := range *text {
		addChar = true
		//		if c == '\n' {
		//c = ' '
		//}

		if flags.IgnorePunctuation {
			if containsRune(punct, c) {
				addChar = false
			}
		}

		if addChar {
			altered = append(altered, c)
		}
	}

	*text = string(altered)
}

func containsRune(list string, char rune) bool {
	for _, c := range list {
		if c == char {
			return true
		}
	}
	return false
}

func containsString(list []string, s string) bool {
	for _, str := range list {
		if s == str {
			return true
		}
	}
	return false
}
