package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/xyproto/guessica"
	"github.com/xyproto/textoutput"
)

func main() {
	o := textoutput.New()
	if appErr := (&cli.App{
		Name:  "guessica",
		Usage: "Update a PKGBUILD by guessing the new pkgbuild and source array",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "version", Aliases: []string{"V"}, Usage: "verbose"},
			&cli.BoolFlag{Name: "i", Usage: "write changes"},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				o.Println(guessica.VersionString)
				os.Exit(0)
			}

			// Check if any arguments are given
			if c.NArg() <= 0 {
				o.Printf("<lightblue>%s</lightblue> <white>%s</white> <lightblue>%s</lightblue>\n", "Please provide one or more", "PKGBUILD", "filenames")
				os.Exit(1)
			}

			pkgbuildFilenames := c.Args().Slice()

			// Treat all arguments as PKGBUILD files that are to be updated
			var err error
			for _, pkgbuildFilename := range pkgbuildFilenames {
				if c.Bool("i") { // Write changes
					o.Printf("<darkgray>[<white>%s<darkgray>] <lightblue>Updating to version</lightblue>... ", filepath.Base(pkgbuildFilename))
				}

				data, err := os.ReadFile(pkgbuildFilename)
				if err != nil {
					o.Printf("<darkred>%s</darkred>\n", err)
					continue
				}
				var (
					pkgbuildContents = string(data)
					ver              string
					sourceLine       string
				)
				for _, site := range guessica.SpecificSites {
					if strings.Contains(pkgbuildContents, site) {
						ver, sourceLine, err = guessica.GuessSourceString(pkgbuildContents, site)
						if err == nil {
							break
						}
					}
				}
				if err != nil {
					o.Printf("<darkred>%s</darkred>\n", err)
					continue
				}

				foundNewVersion := !strings.Contains(pkgbuildContents, "pkgver="+ver)

				var sb strings.Builder
				for _, line := range strings.Split(pkgbuildContents, "\n") {
					if strings.HasPrefix(line, "pkgver=") {
						sb.WriteString("pkgver=" + ver + "\n")
					} else if foundNewVersion && strings.HasPrefix(line, "pkgrel=") {
						sb.WriteString("pkgrel=1\n")
					} else if strings.HasPrefix(line, "source=") {
						sb.WriteString(sourceLine + "\n")
					} else {
						sb.WriteString(line + "\n")
					}
				}
				if c.Bool("i") {
					// Write changes
					err = os.WriteFile(pkgbuildFilename, []byte(strings.TrimSpace(sb.String())), 0664)
					if err != nil {
						o.Printf("<darkred>%s</darkred>\n", err)
						continue
					}
					o.Printf("<cyan>%s</cyan>\n", ver)
				} else {
					// Just output the version number, if found
					o.Printf("<white>Found version: <lightblue>%s</lightblue>\n", ver)
				}
			}
			return err
		},
	}).Run(os.Args); appErr != nil {
		o.ErrExit(appErr.Error())
	}
}
