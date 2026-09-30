package main

import (
	"bufio"
	"fmt"
	"hw3/models"
	"io"
	"os"
	"regexp"
	"strings"
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile("@")
	var seenBrowsers []string // #4 Сделать слайс с заданной размерностью?
	uniqueBrowsers := 0
	var foundUsers strings.Builder

	reader := bufio.NewReader(file)
	i := -1

	for {
		var user models.User
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			i++
			// user := make(map[string]interface{})
			if err := user.UnmarshalJSON(line); err != nil {
				panic(err)
			}

			isAndroid := false
			isMSIE := false

			browsers := user.Browsers

			for _, browser := range browsers {
				if strings.Contains(browser, "Android") { // #2 УЗКОЕ МЕСТО
					isAndroid = true
					notSeenBefore := true
					for _, item := range seenBrowsers {
						if item == browser {
							notSeenBefore = false
						}
					}
					if notSeenBefore {
						// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
						seenBrowsers = append(seenBrowsers, browser)
						uniqueBrowsers++
					}
				}
			}

			for _, browser := range browsers {
				if strings.Contains(browser, "MSIE") { // УЗКОЕ МЕСТО
					isMSIE = true
					notSeenBefore := true
					for _, item := range seenBrowsers {
						if item == browser {
							notSeenBefore = false
						}
					}
					if notSeenBefore {
						// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
						seenBrowsers = append(seenBrowsers, browser)
						uniqueBrowsers++
					}
				}
			}
			if !(isAndroid && isMSIE) {
				continue
			}

			// log.Println("Android and MSIE user:", user["name"], user["email"])
			email := r.ReplaceAllString(user.Email, " [at] ")
			fmt.Fprintf(&foundUsers, "[%d] %s <%s>\n", i, user.Name, email)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

	}
	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
