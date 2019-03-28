package main

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
	OrgToHTML(text)
}
