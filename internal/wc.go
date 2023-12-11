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

	option := getOption()

	fileName := getFileName()

	isPipe := getIsPipe()

	fileContent := getFileContent(isPipe, fileName)

	fileNameSuffix := generateFileNameSuffix(fileName)
	println("fileNameSuffix =", fileNameSuffix)

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

func getOption() string {
	var option string
	flag.StringVar(&option, "option", "default", "Choose an option: b, l, w, c")
	flag.Parse()

	return option
}

func getFileName() string {
	if flag.NArg() > 0 {
		return flag.Arg(0)
	}

	return ""
}

func getIsPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("Error checking if stdin is a pipe:", err)
		os.Exit(1)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func getFileContent(isPipe bool, fileName string) []byte {
	var file []byte
	var err error

	if isPipe {
		file, err = io.ReadAll(os.Stdin)
	} else {
		file, err = os.ReadFile(fileName)
	}

	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
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
	lines := strings.Split(string(fileContent), "\n")
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

	count := 0

	for len(fileContent) > 0 {

		_, size := utf8.DecodeRune(fileContent)
		count++

		fileContent = fileContent[size:]
	}

	return count
}

func printFileStatistics(fileContent []byte, fileNameSuffix string) {
	fmt.Printf("Byte count%s: %d\n", fileNameSuffix, countBytes(fileContent))
	fmt.Printf("Line count%s: %d\n", fileNameSuffix, countLines(fileContent))
	fmt.Printf("Word count%s: %d\n", fileNameSuffix, countWords(fileContent))
	fmt.Printf("Character count%s: %d\n", fileNameSuffix, countCharacters(fileContent))
}