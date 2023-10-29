/*
Copyright © 2023 Sam Wolfs be.samwolfs@gmail.com

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
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Link struct {
	Name string `mapstructure:"name"`
	Tags string `mapstructure:"tags"`
	Url  string `mapstructure:"url"`
}

var list string
var links []Link

func init() {
	showCommand.Flags().StringVarP(&list, "list", "l", "", "Name of list (required)")
	showCommand.MarkFlagRequired("list")

	rootCmd.AddCommand(showCommand)
}

var showCommand = &cobra.Command{
	Use: "show",
	// TODO: descriptions
	Short: "List all links defined under the provided key and open selection.",
	Long: `List all links defined under the provided key and open selection.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			mapstructure.Decode(metadata.Get(list), &links)
			rows := make([]string, len(links))
			for i, link := range links {
				rows[i] = link.format()
			}
			// Enable Markup
			fmt.Println("\x00markup-rows\x1ftrue")
			fmt.Print(strings.Join(rows, "\n"))
		} else {
			var link Link

			response := os.Getenv("ROFI_INFO")
			if response == "" {
				mapstructure.Decode(metadata.Get(list), &links)
				name := strings.Split(args[0], " <span ")[0]

				idx := slices.IndexFunc(links, func(l Link) bool { return l.Name == name })
				link = links[idx]
			} else {
				json.Unmarshal([]byte(response), &link)
			}

			link.open()
		}
	},
}

func (l Link) format() string {
	color := viper.GetString("fgColor")
	if color == "" {
		color = "#928374"
	}
	info, _ := json.Marshal(l)
	return fmt.Sprintf("%s <span color=\"%s\" size=\"10pt\" style=\"italic\">(%s)</span>\x00info\x1f%s", l.Name, color, l.Tags, string(info[:]))
}

func (l Link) open() error {
    var cmd string
    var args []string

    switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default: // "linux", "freebsd", "openbsd", "netbsd"
        cmd = "xdg-open"
    }
    args = append(args, l.Url)
    return exec.Command(cmd, args...).Start()
}
