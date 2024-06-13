package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// import pflag to parse args
// import flag "github.com/spf13/pflag"

const (
	MarkDownExt = ".md"
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

	fmt.Print(files)

	TransferReadList(files, outputDir)

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
		files = append(files, item.Name())
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

	visitedMonth := map[string]struct{}{}
	for _, file := range files {
		date := re.FindString(file)
		if date == "" {
			continue
		}

		month := date[:6]
		outputFile := filepath.Join(destDir, month)
		if _, ok := visitedMonth[month]; !ok {
			fmt.Printf("Need to create new file: %s\n", outputFile)
			visitedMonth[month] = struct{}{}
		}
		fmt.Printf("Transfer file: %s to %s\n", file, outputFile)

	}

	return nil
}

func MakeMonthFile(destDir, month string) error {
	// create a file
	file, err := os.Create(filepath.Join(destDir, month+".md"))
	if err != nil {
		return err
	}
	defer file.Close()

	str := "# Read List of " + month + "\n"

	// This is hugo file header
	str += fmt.Sprintf(" %s\n", month)
	// write some content to the file
	_, err = file.WriteString(fmt.Sprintf("# Read List of %s\n", month))
	if err != nil {
		return err
	}
	return nil
}