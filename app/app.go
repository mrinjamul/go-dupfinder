package app

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
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

// Sha256sumChunks will hash last few chunks of the file and return the sha256sum
func Sha256sumChunks(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.CopyN(hash, file, 1024*1024); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
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
func GetFiles(path string, recursive bool) ([]string, error) {
	var filepaths []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	for _, f := range files {
		if isDir, _ := IsDir(path + "/" + f.Name()); isDir {
			if !recursive {
				continue
			}
			// ignore git files
			if f.Name() == ".git" {
				continue
			}
			filesfromSubDir, err := GetFiles(path+"/"+f.Name(), recursive)
			if err != nil {
				continue
			}
			filepaths = append(filepaths, filesfromSubDir...)
			continue
		}
		filepaths = append(filepaths, path+"/"+f.Name())
	}
	return filepaths, nil
}

// GetExcludeFiles will return the files to exclude
func GetExcludeFiles(exclude string) []string {
	var excludeFiles []string
	if exclude == "" {
		return excludeFiles
	}
	excludeFiles = strings.Split(exclude, ",")
	return excludeFiles
}

// IsExcluded check if a file is in excluded list
func IsExcluded(file string, excludedFiles []string, excludeEmpty bool) bool {
	if excludeEmpty {
		fileSize, err := GetFileSize(file)
		if err != nil {
			panic(err)
		}
		if fileSize == 0 {
			return true
		}
	}

	for _, f := range excludedFiles {
		if strings.Contains(file, f) {
			return true
		}
	}
	return false
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
	for id, f := range filePaths {
		cfmt := color.New(color.FgCyan)
		if id%2 == 0 {
			cfmt.Printf(" " + f)
			continue
		}
		cfmt = color.New(color.FgBlue)
		cfmt.Print(" " + f)
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
			if isDir, _ := IsDir(f); isDir {
				continue
			}
			duplicateFiles = append(duplicateFiles, f)
		}
	}
	return duplicateFiles
}

// Confirm will prompt to user for yes or no
func Confirm(message string) bool {
	var response string
	fmt.Print(message + " :")
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

// GetFileSize return file size of given file
func GetFileSize(filename string) (size int64, err error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// SoftLink will create a soft link with src and dest
func SoftLink(src, dest string, force bool) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		if !force {
			return errors.New("file already exists")
		}
		if err := os.Remove(dest); err != nil {
			return err
		}
	}
	if err := os.Symlink(src, dest); err != nil {
		return err
	}
	return nil
}

// HardLink will create a hard link with src and dest
func HardLink(src, dest string, force bool) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		if !force {
			return errors.New("file already exists")
		}
		if err := os.Remove(dest); err != nil {
			return err
		}
	}
	if err := os.Link(src, dest); err != nil {
		return err
	}
	return nil
}

// ReplaceWithLink will replace a file with a link
func ReplaceWithLink(src, dest string, hard bool) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		if err := os.Remove(dest); err != nil {
			return err
		}
	}
	if hard {
		return HardLink(src, dest, false)
	}
	return SoftLink(src, dest, false)
}

// CheckChunks will check hash of chunks between two files and compare them
func CheckChunks(firstFile, secondFile string) (bool, error) {
	firstFileHash, err := Sha256sumChunks(firstFile)
	if err != nil {
		return false, err
	}
	secondFileHash, err := Sha256sumChunks(secondFile)
	if err != nil {
		return false, err
	}
	if firstFileHash == secondFileHash {
		return true, nil
	}
	return false, nil
}
