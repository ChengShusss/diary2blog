package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
description = "阿树%s的阅读记录，谨供参考"
+++

> "本文是阿树%s的阅读记录，仅供参考，不对真实性和有效性作任何保障。"
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

	for i := 0; i < len(files); i++ {
		file := files[len(files)-1-i]
		date := re.FindString(file)
		if date == "" {
			continue
		}

		month := date[:6]
		outputFile := filepath.Join(destDir, month+".md")
		err := AppendReadList(file, outputFile, date)
		if err != nil {
			return err
		}
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
	t, err := normalizeDate(date)
	if err != nil {
		return err
	}
	yearStr := date[:4]
	monthStr := date[4:6]
	formatStr := fmt.Sprintf(" %s 年 %s 月", yearStr, monthStr)
	s := fmt.Sprintf(DocHeader, date[:6], t.Format(TimeFormat), formatStr, formatStr)
	// fmt.Printf("List: %s\n", s)
	_, err = file.WriteString(s)
	return err
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

		// Only extract readlist section
		if strings.HasPrefix(line, "## ") {
			if hasStart {
				break
			}
			if strings.HasPrefix(line, "## Readlist") {
				hasStart = true
				continue
			}
		}

		// if meet "---" inline then omit this line
		if strings.Contains(line, "---") {
			continue
		}

		// omit empty line(have no content but a dash '-')
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
		err = MakeMonthFile(dst, date)
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

func normalizeDate(date string) (time.Time, error) {

	month, err := strconv.ParseInt(date[:6], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now()
	nowMonth := int64(now.Year()*100 + int(now.Month()))

	fmt.Printf("Now: %v, month: %v\n", nowMonth, month)

	if month >= nowMonth {
		return time.Now(), nil
	}

	t, err := time.Parse("200601", date[:6])
	if err != nil {
		return time.Time{}, err
	}
	nextMonth := t.AddDate(0, 1, 0)
	nextMonthStart := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

	return nextMonthStart, nil
}
