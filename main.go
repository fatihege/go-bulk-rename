package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const sep = string(os.PathSeparator)

func renameExtension(filesFlag, newFlag *string, appendFlag *bool) {
	if strings.TrimSpace(*filesFlag) == "" {
		fmt.Println("ext: the names of the files whose extensions will be renamed are required")
		return
	}

	files := strings.Split(*filesFlag, ",")
	ext := strings.TrimSpace(*newFlag)

	if len(files) == 1 && strings.HasSuffix(strings.TrimSpace(files[0]), "*") {
		tempDir := strings.Split(files[0], sep)
		dir := tempDir[:len(tempDir)-1]
		files = []string{}
		filesInDir, err := os.ReadDir(strings.Join(dir, sep))
		if err != nil {
			panic(err)
		}
		for _, file := range filesInDir {
			files = append(files, fmt.Sprintf("%s%s%s", strings.Join(dir, sep), sep, file.Name()))
		}
	}

	for _, file := range files {
		file = strings.TrimSpace(file)

		statFile, err := os.Stat(file)

		if err != nil {
			fmt.Printf("ext: %v", err)
		}

		if statFile.IsDir() {
			return
		}

		if ext == "" {
			splits := strings.Split(file, ".")
			err = os.Rename(file, strings.Join(splits[:len(splits)-1], "."))
		} else {
			if *appendFlag {
				err = os.Rename(file, fmt.Sprintf("%s.%s", file, ext))
			} else {
				path := strings.Split(file, sep)
				originalFile := file
				realFileName := path[len(path)-1]
				if strings.HasPrefix(realFileName, ".") && strings.Count(realFileName, ".") == 1 {
					file = fmt.Sprintf("%s._", file)
				}
				separatedName := strings.Split(file, ".")
				separatedName[len(separatedName)-1] = ext
				newName := strings.Join(separatedName, ".")

				file = originalFile
				err = os.Rename(originalFile, newName)
			}
		}

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("ext: %s not exists\n", file)
			} else {
				fmt.Printf("ext: %v", err)
			}
		}
	}
}

func renameFilename(filesFlag, newFlag *string, keepExtFlag *bool) {
	if strings.TrimSpace(*filesFlag) == "" {
		fmt.Println("name: files to be renamed are required")
		return
	}

	if strings.TrimSpace(*newFlag) == "" {
		fmt.Println("name: new filename is required")
		return
	}

	files := strings.Split(*filesFlag, ",")
	counters := map[string]int{
		".": 0,
	}

	if len(files) == 1 && strings.HasSuffix(strings.TrimSpace(files[0]), "*") {
		tempDir := strings.Split(files[0], sep)
		dir := tempDir[:len(tempDir)-1]
		files = []string{}
		filesInDir, err := os.ReadDir(strings.Join(dir, sep))
		if err != nil {
			panic(err)
		}
		for _, file := range filesInDir {
			files = append(files, fmt.Sprintf("%s%s%s", strings.Join(dir, sep), sep, file.Name()))
		}
	}

	for _, file := range files {
		file = strings.TrimSpace(file)
		dirSlice := strings.Split(file, sep)
		filename := dirSlice[len(dirSlice)-1]
		filenameSlice := strings.Split(filename, ".")
		extension := ""
		if len(filenameSlice) > 1 && filenameSlice[0] != "" {
			extension = strings.TrimSpace(filenameSlice[len(filenameSlice)-1])
		}
		newFile := fmt.Sprintf("%s%s%s", strings.Join(dirSlice[:len(dirSlice)-1], sep), sep, *newFlag)
		if _, ok := counters[extension]; !ok {
			counters[extension] = 0
		}
		if !*keepExtFlag && counters["."] != 0 {
			newFile = fmt.Sprintf("%s_%d", newFile, counters["."])
		} else if *keepExtFlag && counters[extension] != 0 {
			newFile = fmt.Sprintf("%s_%d", newFile, counters[extension])
		}

		statFile, err := os.Stat(file)

		if err != nil {
			fmt.Printf("name: %v", err)
		}

		if statFile.IsDir() {
			return
		}

		if !*keepExtFlag {
			err = os.Rename(file, newFile)
		} else {
			if extension != "" {
				err = os.Rename(file, fmt.Sprintf("%s.%s", newFile, extension))
			} else {
				err = os.Rename(file, newFile)
			}
		}

		if !*keepExtFlag {
			counters["."]++
		} else {
			counters[extension]++
		}

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("name: %s not exists\n", file)
			} else {
				fmt.Printf("name: %v", err)
			}
		}
	}
}

func main() {
	extMode := flag.Bool(
		"ext",
		false,
		"if you want to rename file extensions then use this option",
	)
	filesFlag := flag.String("files", "", "files to be renamed")
	newFlag := flag.String("new", "", "new extension name")
	appendFlag := flag.Bool("append", false, "specify whether to append or replace the file extension to be renamed")
	keepExtFlag := flag.Bool("keep-ext", false, "if you are renaming filenames choose whether to keep the file extension")

	flag.Usage = func() {
		fmt.Println("\nusage of bulk-rename: ")
		flag.PrintDefaults()
		fmt.Printf(`
examples:
  bulk-rename -files "./*" -new "new-name"
  bulk-rename -files "./*" -new "new-name" -keep-ext
  bulk-rename -files "./file1,./file2" -new "new-name"

  bulk-rename -ext -files "./*" -new "go"
  bulk-rename -ext -files "./*" -new "go" -append
  bulk-rename -ext -files "./file1.cpp,./file2.cpp" -new "go"

`)
	}

	flag.Parse()

	if *filesFlag == "" && *newFlag == "" {
		flag.Usage()
		return
	}

	if *extMode {
		renameExtension(filesFlag, newFlag, appendFlag)
	} else {
		renameFilename(filesFlag, newFlag, keepExtFlag)
	}
}
