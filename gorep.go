package main

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage("test.png")

	text, err := client.Text()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		removeChars(&text)
		fmt.Println(text)
	}
}

func removeChars(text *string) {

	var altered []rune

	for _, c := range *text {
		if c == '\n' {
			c = ' '
		}
		altered = append(altered, c)
	}

	*text = string(altered)
}
