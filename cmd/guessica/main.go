package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/xyproto/guessica"
	"github.com/xyproto/textoutput"
)

const versionString = "guessica 0.0.4"

func main() {
	o := textoutput.New()
	if appErr := (&cli.App{
		Name:  "guessica",
		Usage: "Update a PKGBUILD by guessing the new pkgbuild and source array",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "version", Aliases: []string{"V"}},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				o.Println(versionString)
				os.Exit(0)
			}
			pkgbuildFilenames := []string{"PKGBUILD"}

			// Check if any arguments are given
			if c.NArg() > 0 {
				pkgbuildFilenames = c.Args().Slice()
			} else {
				o.Printf("<lightblue>%s</lightblue> <white>%s</white> <lightblue>%s</lightblue>\n", "Please provide one or more", "PKGBUILD", "filenames")
				os.Exit(1)
			}

			// Treat all arguments as PKGBUILD files that are to be updated
			var err error
			for _, pkgbuildFilename := range pkgbuildFilenames {
				o.Printf("<lightblue>Updating</lightblue> <white>%s</white>... ", pkgbuildFilename)

				data, err := ioutil.ReadFile(pkgbuildFilename)
				if err != nil {
					o.Printf("<darkred>%s</darkred>\n", err)
					continue
				}
				pkgbuildContents := string(data)
				pkgverLine, sourceLine, err := guessica.GuessSourceString(pkgbuildContents)
				if err != nil {
					o.Printf("<darkred>%s</darkred>\n", err)
					continue

				}
				var sb strings.Builder
				for _, line := range strings.Split(pkgbuildContents, "\n") {
					if strings.HasPrefix(line, "pkgver=") {
						sb.WriteString(pkgverLine + "\n")
					} else if strings.HasPrefix(line, "source=") {
						sb.WriteString(sourceLine + "\n")
					} else {
						sb.WriteString(line + "\n")
					}
				}
				err = ioutil.WriteFile(pkgbuildFilename, []byte(strings.TrimSpace(sb.String())), 0664)
				if err != nil {
					o.Printf("<darkred>%s</darkred>\n", err)
					continue
				}
				o.Printf("<lightgreen>%s</lightgreen>\n", "ok")
			}
			return err
		},
	}).Run(os.Args); appErr != nil {
		o.ErrExit(appErr.Error())
	}
}
