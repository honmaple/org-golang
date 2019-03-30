package org

import (
	"bytes"
	"strings"
)

// Options ..
type Options struct {
	Toc    bool
	Offset int
	Escape bool
}

var options = &Options{
	Toc:    true,
	Escape: false,
	Offset: 0,
}

// SetOptions ..
func SetOptions(conf *Options) {
	options = conf
}

// ToHTML ..
func ToHTML(text string, conf *Options) string {
	var buffer bytes.Buffer

	if conf != nil {
		options = conf
	}

	for _, str := range strings.Split(text, "\n") {
		org.Append(str)
	}
	if options.Toc {
		options.Escape = false
		buffer.WriteString(toc.HTML())
	}
	for _, str := range org.Children {
		buffer.WriteString(str.HTML())
	}
	// for _, i := range org.Debug() {
	//	fmt.Println(i)
	// }
	return buffer.String()
}
