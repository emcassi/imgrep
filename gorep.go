package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

type Flags struct {
	IgnoreCase        bool
	IgnorePunctuation bool
	Padding           int
}

func main() {

	flags, pattern, files, _, err := collectArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(files) == 0 {
		fmt.Println("Error: No files provided")
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	res, err := grepImage(client, flags, pattern, files[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}

func newFlags() Flags {
	return Flags{
		IgnoreCase:        false,
		IgnorePunctuation: false,
		Padding:           25,
	}
}

func collectArgs() (Flags, string, []string, []string, error) {
	flags := newFlags()
	var pattern string
	var files []string
	var dirs []string

	args := os.Args[1:]
	for i, arg := range args {

		if arg[0] == '-' {
			if len(files) > 0 {
				return flags, pattern, files, dirs, errors.New("all arguments must come before your files and pattern")
			}
			switch arg[1:] {
			case "ic":
				flags.IgnoreCase = true
			case "ip":
				flags.IgnorePunctuation = true
			case "p":
				break
			default:
				return flags, pattern, files, dirs, errors.New("invalid argument: " + arg)
			}
		} else {
			if i > 0 && args[i-1] == "-p" {
				padding, err := strconv.Atoi(arg)
				if err != nil {
					return flags, pattern, files, dirs, errors.New("padding value must be an integer. you entered: " + arg)
				}

				flags.Padding = padding
				continue
			}
			if pattern == "" {
				pattern = arg
				continue
			}
			switch filepath.Ext(arg) {
			case "":
				if containsString(dirs, arg) {
					continue
				}
				dirs = append(dirs, arg)
			case ".png", ".jpeg", ".jpg", ".bmp":
				if containsString(files, arg) {
					continue
				}
				files = append(files, arg)
			default:
				break
			}
		}
	}

	return flags, pattern, files, dirs, nil
}

func grepImage(client *gosseract.Client, flags Flags, pattern string, filename string) (string, error) {
	var result string
	client.SetImage(filename)

	text, err := client.Text()
	if err != nil {
		return "", errors.New("file: " + filename + " isn't a valid image file") 
	} else {
		cleanData(&text, flags)

		rg, err := regexp.Compile(pattern)
		if err != nil {
			return "", err
		}

		matches := rg.FindAllStringIndex(text, -1)

		for _, match := range matches {
			startFoundIndex := match[0]
			endFoundIndex := match[1]

			var startReturnIndex int
			var endReturnIndex int

			if startFoundIndex-flags.Padding < 0 {
				startReturnIndex = 0
			} else {
				startReturnIndex = startFoundIndex - flags.Padding
			}

			if endFoundIndex+flags.Padding > len(text)-1 {
				endReturnIndex = len(text) - 1
			} else {
				endReturnIndex = endFoundIndex + flags.Padding
			}

			fmt.Println("MATCH " + filename + ": \n" + text[startReturnIndex:endReturnIndex] + "\n")

		}
	}

	return result, nil
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
		if c == '\n' {
			c = ' '
		}

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
