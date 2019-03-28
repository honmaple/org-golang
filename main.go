package main

import (
	"strings"
	"fmt"
	"bytes"
)

// OrgToHTML ..
func OrgToHTML(text string) string {
	var buffer bytes.Buffer

	for _, str := range strings.Split(text, "\n") {
		org.Append(str)
	}
	for _, str := range org.Children {
		s := str.HTML()
		fmt.Println(s)
		buffer.WriteString(s)
	}
	// for _, i := range org.Debug() {
	//	fmt.Println(i)
	// }
	return buffer.String()
}

// main ..
func main() {
	text := `*bold* bold* *bold\* \*bold\* \*bold*

* sssss
	**italic** italic** **italic\** \**italic\** \**italic**
	=code= code= =code\= \=code\= \=code=
		#+begin_example
		ssssss
		exam
		#+end_example
	~code1~ code1~ ~code1\~ \~code1\~ \~code~

	#+begin_example
	ssssss
	exam
	#+end_example

	-----
	- title1
	  sssss
		#+begin_example
		ssssss
		exam
		#+end_example

	- title2
	  vvvva
	  - title3
		vvvvc

| kk  | kv   |
|-----+------|
| sss | cccc |
|     | cccc |
| ddd | vvv  |

** heading2
	#+begin_src bash
	ssssss
	cccccaa
	11cccccaa
	aaa
	#+end_src`
	// TOHTML(text)
	//	text = `
	// ** heading2
	//	#+begin_src bash
	//	ssssss
	//	cccccaa
	//	11cccccaa
	//	aaa
	//	#+end_src`
	OrgToHTML(text)
}
