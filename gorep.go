package main

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

type Flags struct {
	IgnoreCase        bool
	IgnorePunctuation bool
	Padding           int
	Invert            bool
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

func collectArgs() (Flags, string, []string, []string, error) {
	flags := Flags{}

	ignoreCase := flag.Bool("ic", false, "Ignore case when matching")
	ignorePunctuation := flag.Bool("ip", false, "Ignore punctuation when matching")
	invert := flag.Bool("x", false, "Invert match (display lines that do not match)")
	padding := flag.Int("p", 25, "Padding (chars) for displaying matched text")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		return flags, "", nil, nil, errors.New("pattern and at least one file/directory must be provided")
	}

	flags.IgnoreCase = *ignoreCase
	flags.IgnorePunctuation = *ignorePunctuation
	flags.Invert = *invert
	flags.Padding = *padding

	pattern := args[0]
	filesAndDirs := args[1:]

	var files []string
	var dirs []string

	validExts := map[string]bool{
		".png":  true,
		".jpeg": true,
		".jpg":  true,
		".bmp":  true,
	}

	for _, arg := range filesAndDirs {
		ext := filepath.Ext(arg)
		if ext == "" {
			if containsString(dirs, arg) {
				continue
			}
			dirs = append(dirs, arg)
		} else if validExts[ext] {
			if containsString(files, arg) {
				continue
			}
			files = append(files, arg)
		} else {
			return flags, "", nil, nil, fmt.Errorf("invalid file format for %s", arg)
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

			if flags.Invert {
				text = text[:startFoundIndex] + text[endFoundIndex:]
			} else {
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

				result += fmt.Sprintf("MATCH %s:\n%s\n\n", filename, text[startReturnIndex:endReturnIndex])
			}
		}
	}

	if flags.Invert {
		result = text
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
