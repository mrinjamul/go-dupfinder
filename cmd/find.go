/* Package cmd ...
Copyright © 2021 Injamul Mohammad Mollah <mrinjamul@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mrinjamul/go-dupfinder/app"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find duplicate files",
	Long:  `Find duplicate files`,
	Run:   findRun,
}

var (
	flagDelete bool
)

func init() {
	rootCmd.AddCommand(findCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	findCmd.Flags().BoolVarP(&flagDelete, "delete", "d", false, "Find and delete files")
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func findRun(cmd *cobra.Command, args []string) {
	// check if argument exist or not
	if len(args) == 0 {
		fmt.Println("Please provide a path")
		return
	}
	// chech if argument is a valid path
	if _, err := app.IsValidPath(args[0]); err != nil {
		fmt.Println(err)
		return
	}
	// check if argument is more than 1
	if len(args) > 1 {
		fmt.Println("Please provide only one path")
		return
	}
	// get path
	path := args[0]
	// all files list
	var allFiles []string
	// unique files list
	var uniqueFiles []string
	// duplicate files list
	var duplicateFiles []string
	// unique hashes
	var uniqueHash []string
	// hash Maps
	var hashMap = make(map[string]string)

	// get all files from the given path
	allFiles, err := app.GetFiles(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print finding duplicates
	spin := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
	spin.Suffix = " Finding Duplicate files..."
	spin.FinalMSG = "\nProcess Complete !\n"
	spin.Color("green", "bold") // Set the spinner color to a bold green
	spin.Start()                // Start the spinner

	for _, file := range allFiles {
		sum, err := app.Sha256sum(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !app.ContainsString(uniqueHash, sum) {
			uniqueHash = append(uniqueHash, sum)
			hashMap[sum] = file
		}
	}

	spin.Stop() // stop the spinner

	uniqueFiles = app.GetUniqueFiles(hashMap, uniqueHash)
	duplicateFiles = app.GetDuplicateFiles(allFiles, uniqueFiles)

	// Print file statistics
	fmt.Println("Total file(s):", (len(allFiles)),
		" Unique file(s):", len(uniqueFiles),
		" Duplicate file(s):", len(duplicateFiles))

	// print duplicate files
	if len(duplicateFiles) != 0 {
		if ok := app.Confirm("Press y to view duplicate file(s)"); ok {
			app.PrintFiles(duplicateFiles)
		}
	}
	// print unique file(s)
	if len(duplicateFiles) != 0 {
		if ok := app.Confirm("Press y to view unique file(s)"); ok {
			app.PrintFiles(uniqueFiles)
		}
	}

	// For flagDelete
	if flagDelete {
		app.DeleteAllFiles(duplicateFiles)
	} else if len(duplicateFiles) != 0 {
		// prompt to ask user if want to remove duplicates
		ok := app.Confirm("Do you want to delete duplicate files? (y/n)")
		if ok {
			app.DeleteAllFiles(duplicateFiles)
		}
	}
}
