package render

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/honmaple/org-golang/parser"
	"github.com/stretchr/testify/assert"
)

func testFiles() []string {
	fs := make([]string, 0)
	dir := "./testdata"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !strings.HasSuffix(name, ".org") {
			continue
		}
		fs = append(fs, filepath.Join(dir, name))
	}
	return fs
}

func toDocument(buf []byte) *parser.Document {
	d := &parser.Document{
		Sections: &parser.Section{},
		Keywords: map[string]string{
			"TODO": "TODO | DONE | CANCELED",
		},
		Hyperlinks:      []string{"http", "https", "file"},
		TimestampFormat: "2006-01-02 Mon 15:04",
	}
	d.Children = parser.Parse(d, bytes.NewBuffer(buf))
	return d
}

func TestOrg(t *testing.T) {
	files := testFiles()
	for _, file := range files {
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		out := &Org{Document: toDocument(buf)}
		assert.Equal(t, string(buf), out.String())
	}
}
