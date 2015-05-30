package main

import (
	"fmt"
	"github.com/sisteamnik/guseful/chpu"
	"hash/adler32"
	"io"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var inventory = map[string]interface{}{}

func HasChanged(name string, value interface{}) bool {
	changed := true

	if inventory[name] == value {
		changed = false
	} else {
		inventory[name] = value
	}

	return changed
}

func Cut(value, begin, end string) (string, error) {
	bre, err := regexp.Compile(begin)
	if err != nil {
		return "", err
	}
	ere, err := regexp.Compile(end)
	if err != nil {
		return "", err
	}

	bloc := bre.FindIndex([]byte(value))
	eloc := ere.FindIndex([]byte(value))

	if bloc == nil {
		bloc = []int{0, 0}
	}
	if eloc == nil {
		eloc = []int{len(value)}
	}

	return value[bloc[1]:eloc[0]], nil
}

func Hash(value string) string {
	h := adler32.New()
	io.WriteString(h, value)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Versionize(current *Page, value string) string {
	page := current.Site.Pages.ByPath(value)
	if page == nil {
		errhandle(fmt.Errorf(
			"trying to versionize page which does not exist: %s, current: %s",
			value, current.Path))
	}
	c := page.Process().Content()
	h := Hash(c)
	return current.UrlTo(page) + "?v=" + h
}

func Truncate(length int, value string) string {
	if length > len(value) {
		length = len(value)
	}
	return value[0:length]
}

func StripHTML(value string) string {
	return regexp.MustCompile("<[^>]+>").ReplaceAllString(value, "")
}

func MustDate(t time.Time) string {
	n := time.Now()
	if t.Day() == n.Day() && t.Month() == n.Month() && t.Year() == n.Year() {
		return t.Format("15:04")
	}
	if t.Day() == n.Day()-1 && t.Month() == n.Month() && t.Year() == n.Year() {
		return t.Format("вчера в 15:04")
	}
	return t.Format("02.01.2006")
}

var TemplateFuncMap = template.FuncMap{
	"changed":    HasChanged,
	"cut":        Cut,
	"hash":       Hash,
	"version":    Versionize,
	"truncate":   Truncate,
	"strip_html": StripHTML,
	"split":      strings.Split,
	"date":       MustDate,
	"slug":       chpu.Chpu,
}
