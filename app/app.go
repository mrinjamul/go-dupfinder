package app

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

// Sha256sum will hash the file and return the sha256sum
func Sha256sum(filename string) (string, error) {
	// check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", errors.New("file does not exist")
	}
	// check if a file or directory
	if isDir, _ := IsDir(filename); isDir {
		return "", nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// IsDir check if a given path is file or directory
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// IsValidPath check if the path is valid
func IsValidPath(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, errors.New("path does not exist")
	}
	return true, nil
}

// GetFileName will return the file name from a directory path
func GetFileName(path string) string {
	return strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
}

// ListDir will walk through directory
func ListDir(path string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []fs.FileInfo{}, err
	}
	return files, nil
}

// GetFiles will return the file name from a directory path
func GetFiles(path string) ([]string, error) {
	var filepaths []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	for _, f := range files {
		if isDir, _ := IsDir(path + "/" + f.Name()); isDir {
			continue
		}
		filepaths = append(filepaths, path+"/"+f.Name())
	}
	return filepaths, nil
}

// ContainsString check if string exists in array of string
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// PrintFiles will print files in pretter style
func PrintFiles(filePaths []string) {
	for _, f := range filePaths {
		fmt.Print(" " + GetFileName(f))
	}
	fmt.Println()
}

// DeleteFile will delete a file
func DeleteFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	fmt.Println(GetFileName(path) + " deleted")
	return nil
}

// DeleteAllFiles will delete all files in array
func DeleteAllFiles(files []string) {
	for _, f := range files {
		DeleteFile(f)
	}
}

// GetUniqueFiles return unique files
func GetUniqueFiles(hashMap map[string]string, hash []string) []string {
	var uniqueFiles []string
	for _, f := range hash {
		uniqueFiles = append(uniqueFiles, hashMap[f])
	}
	return uniqueFiles
}

// GetDuplicateFiles return duplicate files
func GetDuplicateFiles(AllFiles []string, uniqueFiles []string) []string {
	// subtract unique files from all files
	var duplicateFiles []string
	for _, f := range AllFiles {
		if !ContainsString(uniqueFiles, f) {
			duplicateFiles = append(duplicateFiles, f)
		}
	}
	return duplicateFiles
}

// Confirm will prompt to user for yes or no
func Confirm(message string) bool {
	var response string
	fmt.Print(message + " (yes/no) :")
	fmt.Scanln(&response)

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return false
	}
}
