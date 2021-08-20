/*********************************************************************************
 Copyright © 2020 lin.jiang
 File Name: main.go
 Author: lin.jiang
 Email: mail@honmaple.com
 Created: 2020-12-28 17:52:36 (CST)
 Last Update: Thursday 2021-10-14 12:03:51 (CST)
		  By:
 Description:
 *********************************************************************************/
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

	r := render.HTMLRender{
		Document: parser.Parse(string(content)),
	}
	fmt.Println(r.String())

	r1 := render.DebugRender{
		Document: parser.Parse(string(content)),
	}
	fmt.Println(r1.String())
}
