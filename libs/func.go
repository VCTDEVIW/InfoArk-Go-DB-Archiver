package vastcom_infoshare

import (
	. "fmt"
	"strconv"
	"io/ioutil"
	"os"
    "path/filepath"
	"strings"
)

func init() {
	Print("")
}

func ToStr(param interface{}) string {
	return Sprintf("%s", param)
}

func ToInt(param interface{}) int {
	raw_param := Sprintf("%s", param)
	num, err := strconv.Atoi(raw_param)
    if err != nil {
        Println("Error:", err)
        return -1
    }

	return num
}

func DirListFiles(dir string) ([]string, error) {
    var files []string

    entries, err := ioutil.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    for _, entry := range entries {
        if !entry.IsDir() {
            filePath := filepath.Join(dir, entry.Name())
            files = append(files, filePath)
        }
    }

    return files, nil
}

func DirListFilesCount(dir string) int {
	var countFs int

    workDir := dir
	Println("Inspect DirListFilesCount() at directory path: " + workDir)

    files, err := DirListFiles(workDir)
    if err != nil {
        Println("Error:", err)
        return -1
    }

    for _, file := range files {
        Println(file)
		countFs++
    }

	ret := Sprintf("Total files count: =%d at %s", countFs, workDir)
	Println(ret)
	return countFs
}

func Explode(str string, delimiter string) []string {
	var parts []string
	i := 0

	for i < len(str) {
		x := strings.Index(str[i:], delimiter) 
		if x < 0 {
			parts = append(parts, str[i:])
			break
		}
		parts = append(parts, str[i:i+x])
		i = i + x + len(delimiter)
	}

	return parts
}

func MoveFile(oldPath string, newPath string, filename string) {
	// Source file path
	sourcePath := oldPath

    // Destination directory path
    destDir := newPath

    // Destination file path
    destPath := filepath.Join(destDir, filename)

    // Check if the destination directory exists
    _, err := os.Stat(destDir)
    if os.IsNotExist(err) {
        // Create the destination directory
        err = os.MkdirAll(destDir, 0755)
        if err != nil {
            Println("Error creating directory:", err)
            return
        }
        Println("Directory created:", destDir)
    }

    // Move the file
    err = os.Rename(sourcePath, destPath)
    if err != nil {
        Println("Error moving file:", err)
        return
    }

    Println("File moved successfully!")
}

func GetFilename(path string) string {
    // Get the base name (filename) from the path
    filename := filepath.Base(path)
    return filename
}

