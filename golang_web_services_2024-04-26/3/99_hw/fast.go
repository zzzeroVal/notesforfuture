package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	var seenBrowsers []string
	uniqueBrowsers := 0
	foundUsers := ""

	reader := bufio.NewReader(file)
	i := -1

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			i++
			user := make(map[string]interface{})
			err := json.Unmarshal(line, &user)
			if err != nil {
				panic(err)
			}

			isAndroid := false
			isMSIE := false

			browsers, ok := user["browsers"].([]interface{})
			if !ok {
				// log.Println("cant cast browsers")
				continue
			}

			for _, browserRaw := range browsers {
				browser, ok := browserRaw.(string)
				if !ok {
					// log.Println("cant cast browser to string")
					continue
				}
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

			for _, browserRaw := range browsers {
				browser, ok := browserRaw.(string)
				if !ok {
					// log.Println("cant cast browser to string")
					continue
				}
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
			email := r.ReplaceAllString(user["email"].(string), " [at] ")
			foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email) // тут надо сделать через Builder, потому что мы каждый раз делаем аллокацию
		}
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

	}
	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
