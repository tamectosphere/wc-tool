package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	argsLength := len(os.Args)

	stat, _ := os.Stdin.Stat()
	isPipe := (stat.Mode() & os.ModeCharDevice) == 0

	if isPipe == true {
		if argsLength != 1 && argsLength != 2 {
			log.Fatal("Usage: <std> | gowc <option>")
		}

		scanner := bufio.NewScanner(os.Stdin)
		var input strings.Builder

		for scanner.Scan() {
			input.WriteString(scanner.Text() + "\n")
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		option := getOptionFromArgsPipe(argsLength, os.Args)
		if option == nil {
			if err := printFileStatisticsPipe(input); err != nil {
				log.Fatal(err)
			}
			return
		}

		switch *option {
		case "-c":
			printByteCountPipe(input)
		case "-l":
			printLineCountPipe(input)
		case "-w":
			printWordCountPipe(input)
		case "-m":
			printCharacterCountPipe(input)
		default:
			log.Fatal("Invalid option")
		}

	} else {

		if argsLength != 3 && argsLength != 2 {
			log.Fatal("Usage: gowc <option> <filename>")
		}

		filename, option := getFileNameAndOption(argsLength, os.Args)

		if option == nil {
			if err := printFileStatistics(filename); err != nil {
				log.Fatal(err)
			}
			return
		}

		switch *option {
		case "-c":
			printByteCount(filename)
		case "-l":
			printLineCount(filename)
		case "-w":
			printWordCount(filename)
		case "-m":
			printCharacterCount(filename)
		default:
			log.Fatal("Invalid option")
		}
	}
}

// Priv function for pipe
func getOptionFromArgsPipe(argsLength int, args []string) *string {
	if argsLength == 1 {
		return nil
	}

	return &args[1]
}

func printFileStatisticsPipe(content strings.Builder) error {
	byteCount, err := countBytesPipe(content)
	if err != nil {
		log.Fatal(err)
	}

	lineCount, err := countLinesPipe(content)
	if err != nil {
		log.Fatal(err)
	}

	wordCount, err := countWordsPipe(content)
	if err != nil {
		log.Fatal(err)
	}

	characterCount, err := countCharactersPipe(content)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Byte count: %d\n", byteCount)
	fmt.Printf("Line count: %d\n", lineCount)
	fmt.Printf("Word count: %d\n", wordCount)
	fmt.Printf("Character count: %d\n", characterCount)

	return nil
}

func printByteCountPipe(content strings.Builder) {
	byteCount, err := countBytesPipe(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Byte count: %d\n", byteCount)
}

func printLineCountPipe(content strings.Builder) {
	lineCount, err := countLinesPipe(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Line count: %d\n", lineCount)
}

func printWordCountPipe(content strings.Builder) {
	wordCount, err := countWordsPipe(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Word count: %d\n", wordCount)
}

func printCharacterCountPipe(content strings.Builder) {
	characterCount, err := countCharactersPipe(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Character count: %d\n", characterCount)
}

func countBytesPipe(content strings.Builder) (int, error) {
	return len(trimmedContent(content)), nil
}

func countLinesPipe(content strings.Builder) (int, error) {
	lines := strings.Split(trimmedContent(content), "\n")
	return len(lines), nil
}

func countWordsPipe(content strings.Builder) (int, error) {
	var wordCount int

	scanner := bufio.NewScanner(strings.NewReader(trimmedContent(content)))

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		wordCount += len(words)
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return wordCount, nil
}

func countCharactersPipe(content strings.Builder) (int, error) {
	var count int

	// Trim BOM if present
	trimmedContent := trimBOM(trimmedContent(content))

	for len(trimmedContent) > 0 {
		_, size := utf8.DecodeRuneInString(trimmedContent)
		count++
		trimmedContent = trimmedContent[size:]
	}

	return count, nil
}

func trimBOM(s string) string {

	if len(s) >= 3 && s[0] == 0xef && s[1] == 0xbb && s[2] == 0xbf {
		return s[3:]
	}
	return s
}

func trimmedContent(content strings.Builder) string {
	return strings.TrimRight(content.String(), "\n")
}

// Priv function for non-pipe
func getFileNameAndOption(argsLength int, args []string) (string, *string) {
	if argsLength == 2 {
		return args[1], nil
	}
	return args[2], &args[1]
}

func printFileStatistics(filename string) error {
	byteCount, err := countBytes(filename)
	if err != nil {
		log.Fatal(err)
	}

	lineCount, err := countLines(filename)
	if err != nil {
		log.Fatal(err)
	}

	wordCount, err := countWords(filename)
	if err != nil {
		log.Fatal(err)
	}

	characterCount, err := countCharacters(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Byte count in %s: %d\n", filename, byteCount)
	fmt.Printf("Line count in %s: %d\n", filename, lineCount)
	fmt.Printf("Word count in %s: %d\n", filename, wordCount)
	fmt.Printf("Character count in %s: %d\n", filename, characterCount)

	return nil
}

func printByteCount(filename string) {
	byteCount, err := countBytes(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Byte count in %s: %d\n", filename, byteCount)
}

func printLineCount(filename string) {
	lineCount, err := countLines(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Line count in %s: %d\n", filename, lineCount)
}

func printWordCount(filename string) {
	wordCount, err := countWords(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Word count in %s: %d\n", filename, wordCount)
}

func printCharacterCount(filename string) {
	characterCount, err := countCharacters(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Character count in %s: %d\n", filename, characterCount)
}

func countBytes(filename string) (int, error) {
	fileContent, err := os.ReadFile(filename)

	if err != nil {
		return 0, err
	}
	return len(fileContent), nil
}

func countLines(filename string) (int, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(fileContent), "\n")

	return len(lines), nil
}

func countWords(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	wordCount := 0
	for scanner.Scan() {
		wordCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return wordCount, nil
}

func countCharacters(filename string) (int, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	if len(fileContent) >= 3 && fileContent[0] == 0xef && fileContent[1] == 0xbb && fileContent[2] == 0xbf {
		fileContent = fileContent[3:]
	}

	count := 0

	for len(fileContent) > 0 {

		_, size := utf8.DecodeRune(fileContent)
		count++

		fileContent = fileContent[size:]
	}

	return count, nil
}
