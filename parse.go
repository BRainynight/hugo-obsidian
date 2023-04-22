package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"strings"
	"path/filepath"
	"net/url"
)

// parse single file for links
func parse(dir, pathPrefix string) []Link {
	// read file
	source, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}

	// parse md
	var links []Link
	fmt.Printf("[Parsing note] %s => ", trim(dir, pathPrefix, ".md"))

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	var n int
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		target, ok := s.Attr("href")
		if !ok {
			target = "#"
		}

		target = processTarget(target)
		source := processSource(trim(dir, pathPrefix, ".md"))

		abs_target_path := target
		u, err := url.ParseRequestURI(abs_target_path)
		if err != nil || u.Scheme == "" || u.Host == "" {
			abs_target_path = filepath.Join(source,"..",target)
			abs_target_path = filepath.ToSlash(abs_target_path)
		} else {
			abs_target_path = target
		}

		// fmt.Printf("  '%s' => %s\n", source, target)
		if !strings.HasPrefix(text, "^"){
			links = append(links, Link{
				Source: source,
				Target: abs_target_path, // target
				Text:   text,
			})
			n++
		}
	})
	fmt.Printf("found: %d links\n", n)

	return links
}
