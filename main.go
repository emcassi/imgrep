package main

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/otiai10/gosseract/v2"
)

// Flags represents command-line flags for GrepImage tool.
type Flags struct {
	IgnoreCase        bool // Ignore case when matching.
	IgnorePunctuation bool // Ignore punctuation when matching.
	Invert            bool // Invert match (display lines that do not match).
	Padding           int  // Padding (chars) for displaying matched text.
}

type FileResult struct {
	Filename string
	Result []string
}

func main() {
	// Collect command-line arguments and flags.
	flags, pattern, files, _, err := collectArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check if files are provided.
	if len(files) == 0 {
		fmt.Println("Error: No files provided")
		return
	}

	

	var wg sync.WaitGroup
	results := make(map[string][]string)

	resultChan := make(chan FileResult)

	// Process each file and perform text extraction and pattern matching.
	for _, file := range files {
		wg.Add(1)
		go grepImage(flags, file, pattern, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		results[result.Filename] = result.Result
	}

	for filename, res := range results {
		fmt.Printf("MATCHES IN FILE: %s\n", filename)
		for _, r := range res {
			fmt.Println("MATCH: ", r)
		}
	}
}

// collectArgs collects and parses command-line arguments and flags.
func collectArgs() (Flags, string, []string, []string, error) {
	flags := Flags{}

	// Define command-line flags.
	ignoreCase := flag.Bool("ic", false, "Ignore case when matching")
	ignorePunctuation := flag.Bool("ip", false, "Ignore punctuation when matching")
	invert := flag.Bool("x", false, "Invert match (display lines that do not match)")
	padding := flag.Int("p", 25, "Padding (chars) for displaying matched text")
	flag.Parse()

	// Parse command-line arguments and flags.
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

	// Define valid image file extensions.
	validExts := map[string]bool{
		".png":  true,
		".jpeg": true,
		".jpg":  true,
		".bmp":  true,
	}

	// Separate files and directories based on their file extensions.
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

func grepImage(flags Flags, file, pattern string, resultChan chan<- FileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	// Initialize the OCR client.
	client := gosseract.NewClient()
	defer client.Close()

	// Extract text from the image file.
	text, err := extractText(client, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Pre-process text based on flags.
	cleanData(&text, &pattern, flags)

	// Perform pattern matching on the text.
	res, err := grep(text, flags, pattern, file)
	if err != nil {
		fmt.Printf("%s: %s\n", file, err)
		return 
	}

	resultChan <- FileResult{
		Filename: file,
		Result: res,
	}
}
// extractText extracts text from an image file using OCR.
func extractText(client *gosseract.Client, filename string) (string, error) {
	client.SetImage(filename)

	text, err := client.Text()
	if err != nil {
		return "", errors.New("file: " + filename + " isn't a valid image file")
	}

	return text, nil
}

// grep performs pattern matching on the given text based on the provided flags.
func grep(text string, flags Flags, pattern string, filename string) ([]string, error) {
	result := []string{}

	rg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	matches := rg.FindAllStringIndex(text, -1)

	if len(matches) == 0 {
		return nil, errors.New("found no matches")
	}

	lastMatch := []int{0, 0}

	for _, match := range matches {
		startFoundIndex := match[0]
		endFoundIndex := match[1]

		if flags.Invert {
			match := text[lastMatch[1]:startFoundIndex]

			if len(result) == 0 {
				result = append(result, match)
			} else {
				result[0] += match
			}
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

			result = append(result, text[startReturnIndex:endReturnIndex])
		}

		lastMatch = match
	}

	if flags.Invert {
		result[0] += text[lastMatch[1]:]
	}

	return result, nil
}

// cleanData pre-processes the text based on the provided flags.
func cleanData(text, pattern *string, flags Flags) {
	var altered []rune
	var addChar bool

	punct := ",.!?:;'=[](){}\\|/~“”’`"

	if flags.IgnoreCase {
		*pattern = "(?i)" + *pattern
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

// containsRune checks if a rune exists in a given string.
func containsRune(list string, char rune) bool {
	for _, c := range list {
		if c == char {
			return true
		}
	}
	return false
}

// containsString checks if a string exists in a given string slice.
func containsString(list []string, s string) bool {
	for _, str := range list {
		if s == str {
			return true
		}
	}
	return false
}
