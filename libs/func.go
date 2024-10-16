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

func StrExist(src string, target string) bool {
    if strings.Contains(src, target) {
		return true
	} else {
        return false
    }
}

func Strlen(str string) int {
    return len(str)
}

func Upper(str string) string {
    return strings.ToUpper(str)
}

func Lower(str string) string {
    return strings.ToLower(str)
}

func ReverseStr(s string) string {
	runes := []rune(s) // Convert string to rune slice to handle Unicode characters
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i] // Swap characters
	}
	return string(runes) // Convert rune slice back to string
}

func str_replace(src string, target string, replaceAs string) string {
    return strings.ReplaceAll(src, target, replaceAs)
}

func StrRep(src string, target string, replaceAs string) string {
    return str_replace(src, target, replaceAs)
}

func Strpos(src string, target string) int {
    // String slice counting from [0]...[n]
	find_idx := strings.Index(src, target)
	return find_idx
}

func Substr(args ...any) string {
	switch len(args) {
	case 2:
		str := args[0].(string)
		count_str := Strlen(str)
		start_idx := args[1].(int)

		if (start_idx < 0) || (start_idx > count_str) {
			Println("Error using Substr()! start_idx must in between 0 and max strlen.")
			return ""
		}

		return str[start_idx:]
	case 3:
		str := args[0].(string)
		count_str := Strlen(str)
		start_idx := args[1].(int)
		end_idx := args[2].(int)

		if (start_idx < 0) || (start_idx > count_str) {
			Println("Error using Substr()! start_idx must in between 0 and max strlen.")
			return ""
		}

		if (end_idx < start_idx) || (end_idx > count_str) {
			Println("Error using Substr()! Please inspect the underlying co-relation between start_idx and end_idx before use.")
			return ""
		}

		return str[start_idx:end_idx]
	default:
		Println("Error using Substr()!")
		return ""
	}
}

