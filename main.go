package main

import (
	"fmt"
	"io/ioutil"

	"github.com/honmaple/org-golang/parser"
	"github.com/honmaple/org-golang/render"
)

// https://orgmode.org/worg/dev/org-syntax.html
func main() {
	content, _ := ioutil.ReadFile("test.org")

	r := render.HTML{
		Toc:      true,
		Document: parser.Parse(string(content)),
	}
	fmt.Println(r.String())

	r1 := render.Debug{
		Document: parser.Parse(string(content)),
	}
	fmt.Println(r1.String())
}
