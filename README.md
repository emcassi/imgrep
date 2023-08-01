# imgrep

## Introduction

imgrep is a command-line tool that allows users to search for a given pattern in images. It uses optical character recognition (OCR) to extract text from images and performs pattern matching on the extracted text.

The tool is designed to be flexible and configurable, supporting various options like case-insensitive matching, ignoring punctuation, inverting the match, and providing a context around the matched text.

## Features

- Perform text pattern matching in images.
- Configurable flags for case-insensitive and punctuation-ignoring matching.
- Invert match to display lines that do not match the pattern.
- Display context around the matched text with padding.

## Installation

To use imgrep, you need to have Go installed on your system.

```bash
git clone https://github.com/emcassi/imgrep.git
cd imgrep
go build -o imgrep
```

## Usage

To search for a pattern in an image, use the following command:

```bash
./imgrep [flags] pattern file.png [file2.png ...]
```

### Flags

- -ic: Ignore case when matching.
- -ip: Ignore punctuation when matching.
- -x: Invert match (display lines that do not match the pattern).
- -p: Padding (characters) for displaying matched text.

### Arguments

- Arg 1 : pattern - accepts regex
- Arg 2+ : file names

Examples:

```bash
./imgrep -ic -ip -p 10 hello image1.png
./imgrep -ic hello imgs/image.png # image in a subfolder
./imgrep -ic hello image1.png image2.png # input for multiple images
./imgrep -p 50 '(?i)john' image1.png # supports go style regex statements
./imgrep -x error image2.png # invert - exclude the pattern in the results
```

## Roadmap

- [x] Add ocr functionality 
- [x] Add argument parsing
- [x] Be able to grep a single image
- [x] Be able to grep multiple images
- [ ] Be able to pass directories of images
- [x] Add concurrency support
- [ ] Further Testing
- [ ] Add documentation

## Contributing

Contributions are welcome! If you want to contribute to imgrep, please follow these steps:

- Read the [Contributing guidelines](CONTRIBUTING.md)
- Fork the repository and create your branch from main.
- Make sure your code follows the Go coding style.
- Run tests to ensure your changes don't break existing functionality.
- Open a pull request with a clear description of your changes.

## License

This project is licensed under the MIT License.

## Acknowledgements

[gosseract](https://github.com/otiai10/gosseract): A Go library for OCR.

## Contact

For any questions or suggestions, feel free to get in touch:

GitHub: [emcassi](http://github.com/emcassi)

Email: <alex.wayne.dev@gmail.com>
