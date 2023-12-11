# Wc Tool

## Overview

The `wc-tool` is a command-line utility written in Go that provides word counting and file statistics functionalities. It allows you to analyze the content of text files, either through standard input (pipe) or by specifying a file name and an option.

This tool was inspired by the [Word Count Challenge](https://codingchallenges.fyi/challenges/challenge-wc/), providing a practical implementation of word counting and file statistics with additional features.

## Features

- Count bytes, lines, words, and characters in text files.
- Handle both standard input and file input.

## Usage

### Standard Input (Pipe)

To use the `wc-tool` with standard input, you can pipe the content to the tool and specify an option:

```bash
<some_command> | wc-tool -option <option>
```
### Example

```bash
echo 'Hello, world!' | wc-tool -option w
>>> Word count: 2
```
### File Input

To use the wc-tool with a file, provide the file name and an option:

```bash
wc-tool -option <option> <filename>
```
### Example

```bash
wc-tool -option b example.txt
>>> Byte count in example.txt: 51
```
## Available Options

- `b`: Print byte count.
- `l`: Print line count.
- `w`: Print word count.
- `m`: Print character count.
- No option: Print byte, line, word, and character count.
