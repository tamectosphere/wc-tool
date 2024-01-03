package wc

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

func Process(args []string) map[string]interface{} {

	option, text, fileName, ignorePipe := getFlagOptions()

	isPipe := getIsPipe(ignorePipe)

	if !isPipe {
		validateTextAndFileName(text, fileName)
	}

	fileContent := getFileContent(isPipe, text, fileName)

	fileNameSuffix := generateFileNameSuffix(fileName)

	result := make(map[string]interface{})

	switch option {
	case "b":
		byteCount := countBytes(fileContent)
		fmt.Printf("Byte count%s: %d\n", fileNameSuffix, byteCount)
		result["byteCount"] = byteCount
	case "l":
		lineCount := countLines(fileContent)
		fmt.Printf("Line count%s: %d\n", fileNameSuffix, countLines(fileContent))
		result["lineCount"] = lineCount
	case "w":
		wordCount := countWords(fileContent)
		fmt.Printf("Word count%s: %d\n", fileNameSuffix, countWords(fileContent))
		result["wordCount"] = wordCount
	case "c":
		characterCount := countCharacters(fileContent)
		fmt.Printf("Character count%s: %d\n", fileNameSuffix, countCharacters(fileContent))
		result["characterCount"] = characterCount
	case "default":
		result["fileStatistics"] = printFileStatistics(fileContent, fileNameSuffix)
	default:
		log.Fatal("Invalid option")
	}

	return result
}

func getFlagOptions() (string, string, string, bool) {
	var option string
	var text string
	var fileName string
	var ignorePipe bool

	flag.StringVar(&option, "option", "default", "Choose an option: b, l, w, c")
	flag.StringVar(&text, "text", "", "Input text to process")
	flag.StringVar(&fileName, "file", "", "Input file name to process")
	flag.BoolVar(&ignorePipe, "ignore-pipe", false, "Ignore pipe condition (stdin not treated as a pipe)")
	flag.Parse()

	return option, text, fileName, ignorePipe
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

func getIsPipe(ignorePipe bool) bool {
	if ignorePipe {
		return false
	}

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
	lines := strings.Split(getContentString(fileContent, true), "\n")
	return len(lines)
}

func countWords(fileContent []byte) int {
	scanner := bufio.NewScanner(bytes.NewReader(fileContent))
	scanner.Split(bufio.ScanWords)

	wordCount := 0
	firstNonEmpty := false

	whitespaceRegex := regexp.MustCompile(`[\s\p{C}]`)

	for scanner.Scan() {
		word := scanner.Text()

		cleanedWord := whitespaceRegex.ReplaceAllString(word, " ")
		trimmedWord := strings.TrimSpace(cleanedWord)

		if !firstNonEmpty && trimmedWord == "" {

			continue
		}

		firstNonEmpty = true
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

	return utf8.RuneCountInString(getContentString(fileContent, false))
}

func getContentString(fileContent []byte, isTrim bool) string {
	if isTrim {
		return strings.TrimSpace(string(fileContent))
	}

	return string(fileContent)

}

func printFileStatistics(fileContent []byte, fileNameSuffix string) map[string]interface{} {
	byteCount := countBytes(fileContent)
	lineCount := countLines(fileContent)
	wordCount := countWords(fileContent)
	characterCount := countCharacters(fileContent)

	fmt.Printf("Byte count%s: %d\n", fileNameSuffix, byteCount)
	fmt.Printf("Line count%s: %d\n", fileNameSuffix, lineCount)
	fmt.Printf("Word count%s: %d\n", fileNameSuffix, wordCount)
	fmt.Printf("Character count%s: %d\n", fileNameSuffix, characterCount)

	stats := make(map[string]interface{})
	stats["byteCount"] = byteCount
	stats["lineCount"] = lineCount
	stats["wordCount"] = wordCount
	stats["characterCount"] = characterCount

	return stats
}
