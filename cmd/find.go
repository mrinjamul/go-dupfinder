/*
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

	"github.com/mrinjamul/go-dupfinder/app"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find duplicate files by sha256sum",
	Long:  `Find duplicate files by sha256sum`,
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
	path := args[0]
	var filepath []string
	// var hashMap = make(map[string][]string)
	var hashs []string
	dupCount := 0

	filepath, err := app.GetFiles(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	app.PrintFiles(filepath)

	if flagDelete {
		fmt.Println("Finding and Deleting duplicates...")
		for _, file := range filepath {
			sum, err := app.Sha256sum(file)
			if err != nil {
				fmt.Println(err)
				return
			}
			if app.ContainsString(hashs, sum) {
				dupCount++
				err := app.DeleteFile(file)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				hashs = append(hashs, sum)
			}

			// print hash sum
			// fmt.Println(app.GetFileName(file) + " : " + sum)
		}
		// print how many duplicate founds
		fmt.Println("Duplicate founds :", dupCount)
		return
	}

	// print finding duplicates
	fmt.Println("Finding duplicates...")

	for _, file := range filepath {
		sum, err := app.Sha256sum(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		if app.ContainsString(hashs, sum) {
			dupCount++
		} else {
			hashs = append(hashs, sum)
		}

		// print hash sum
		// fmt.Println(app.GetFileName(file) + " : " + sum)
	}
	// print how many duplicate founds
	fmt.Println("Duplicate founds :", dupCount)
}
