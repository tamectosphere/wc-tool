package wc

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func Process(args []string) {

	option, text, fileName := getFlagOptions()

	isPipe := getIsPipe()

	if !isPipe {
		validateTextAndFileName(text, fileName)
	}

	fileContent := getFileContent(isPipe, text, fileName)

	fileNameSuffix := generateFileNameSuffix(fileName)

	switch option {
	case "b":
		fmt.Printf("Byte count%s: %d\n", fileNameSuffix, countBytes(fileContent))
	case "l":
		fmt.Printf("Line count%s: %d\n", fileNameSuffix, countLines(fileContent))
	case "w":
		fmt.Printf("Word count%s: %d\n", fileNameSuffix, countWords(fileContent))
	case "c":
		fmt.Printf("Character count%s: %d\n", fileNameSuffix, countCharacters(fileContent))
	case "default":
		printFileStatistics(fileContent, fileNameSuffix)
	default:
		log.Fatal("Invalid option")
	}
}

func getFlagOptions() (string, string, string) {
	var option string
	var text string
	var fileName string
	flag.StringVar(&option, "option", "default", "Choose an option: b, l, w, c")
	flag.StringVar(&text, "text", "", "Input text to process")
	flag.StringVar(&fileName, "file", "", "Input file name to process")
	flag.Parse()

	return option, text, fileName
}

func getText() string {
	var text string
	flag.StringVar(&text, "text", "", "Input text to process")
	flag.Parse()

	return text
}

func getFileName() string {
	var fileName string
	flag.StringVar(&fileName, "file", "", "Input file name to process")
	flag.Parse()

	return fileName
}

func validateTextAndFileName(text string, fileName string) {
	if text == "" && fileName == "" {
		log.Fatal("both 'text' and 'fileName' options are empty, at least one must be provided")
	}

	if text != "" && fileName != "" {
		log.Fatal("both 'text' and 'fileName' options are provided, but only one should be specified")
	}

	return
}

func getIsPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func getFileContent(isPipe bool, text string, fileName string) []byte {
	var file []byte
	var err error

	if isPipe {
		file, err = io.ReadAll(os.Stdin)
	} else if text != "" {
		file = []byte(text)
	} else {
		file, err = os.ReadFile(fileName)
	}

	if err != nil {
		log.Fatal(err)
	}

	return file
}

func generateFileNameSuffix(fileName string) string {
	if fileName != "" {
		return " in " + fileName
	}
	return ""
}

func countBytes(fileContent []byte) int {
	return len(fileContent)
}

func countLines(fileContent []byte) int {
	lines := strings.Split(getContentString(fileContent), "\n")
	return len(lines)
}

func countWords(fileContent []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(fileContent))
	scanner.Split(bufio.ScanWords)

	wordCount := 0
	for scanner.Scan() {
		wordCount++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return wordCount
}

func countCharacters(fileContent []byte) int {

	if len(fileContent) >= 3 && fileContent[0] == 0xef && fileContent[1] == 0xbb && fileContent[2] == 0xbf {
		fileContent = fileContent[3:]
	}

	return utf8.RuneCountInString(getContentString(fileContent))
}

func getContentString(fileContent []byte) string {
	return strings.TrimSpace(string(fileContent))
}

func printFileStatistics(fileContent []byte, fileNameSuffix string) {
	fmt.Printf("Byte count%s: %d\n", fileNameSuffix, countBytes(fileContent))
	fmt.Printf("Line count%s: %d\n", fileNameSuffix, countLines(fileContent))
	fmt.Printf("Word count%s: %d\n", fileNameSuffix, countWords(fileContent))
	fmt.Printf("Character count%s: %d\n", fileNameSuffix, countCharacters(fileContent))
}
