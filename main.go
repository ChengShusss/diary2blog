package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// import pflag to parse args
// import flag "github.com/spf13/pflag"

const (
	MarkDownExt = ".md"

	TimeFormat = "2006-01-02T15:04:05+08:00"
	DocHeader  = `+++
title = 'ReadList - %s'
date = %s
draft = false
+++

`
)

var (
	visitedMonth = map[string]struct{}{}
)

func main() {

	if len(os.Args) < 3 {
		panic("Please provide a directory path and output directory")
	}

	dir := os.Args[1]
	outputDir := os.Args[2]
	files, err := GetFiles(dir)
	if err != nil {
		panic(err)
	}

	// fmt.Print(files)
	fmt.Printf("Total %d files\n", len(files))

	err = TransferReadList(files, outputDir)
	if err != nil {
		fmt.Printf("Failed to transfer, err: %v\n", err)
	}
}

func GetFiles(dir string) ([]string, error) {
	// get all direct files in dir
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, item := range entries {
		if item.IsDir() {
			continue
		}
		name := item.Name()
		if !strings.HasSuffix(name, MarkDownExt) {
			continue
		}
		files = append(files, filepath.Join(dir, item.Name()))
	}

	return files, nil
}

func TransferReadList(files []string, destDir string) error {
	// TODO maybe could support specific filename format to parse

	// This must complie, just for normalize the code style
	re, err := regexp.Compile(`\d{8}`)
	if err != nil {
		return err
	}

	for _, file := range files {
		date := re.FindString(file)
		if date == "" {
			continue
		}

		month := date[:6]
		outputFile := filepath.Join(destDir, month+".md")
		AppendReadList(file, outputFile, date)
	}

	return nil
}

func MakeMonthFile(name, date string) error {
	// create a file
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	// This is hugo file header
	current := time.Now().Format(TimeFormat)
	_, err = file.WriteString(fmt.Sprintf(DocHeader, date, current))
	if err != nil {
		return err
	}

	return nil
}

func AppendReadList(src, dst, date string) error {
	// append the content to the file
	// fmt.Printf("Transfer file: %s to %s\n", src, dst)

	// open the src file
	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	hasStart := false
	scanner := bufio.NewScanner(srcFile)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") {
			if hasStart {
				break
			}
			if strings.HasPrefix(line, "## Readlist") {
				hasStart = true
				continue
			}
		}
		if hasStart {
			trimed := strings.TrimSpace(line)
			if len(trimed) == 0 || trimed == "-" {
				continue
			}
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return nil
	}

	month := date[:6]
	if _, ok := visitedMonth[month]; !ok {
		fmt.Printf("Need to create new file: %s\n", filepath.Base(dst))
		visitedMonth[month] = struct{}{}
		err = MakeMonthFile(dst, month)
		if err != nil {
			return err
		}
	}

	// write to dst file
	dstFile, err := os.OpenFile(dst, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	dstFile.WriteString(fmt.Sprintf("## %s\n\n%s\n\n", date, strings.Join(lines, "\n")))
	fmt.Printf("Src: %v\n  Total %d lines\n", filepath.Base(src), len(lines))

	return nil
}
