package render

import (
	"bytes"
	"fmt"
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

// BenchmarkSprintf-8		11847894			99.43 ns/op
// BenchmarkPlus-8			1000000000			 0.2529 ns/op
// BenchmarkBuilder-8		22237069			52.56 ns/op
func BenchmarkSprintf(b *testing.B) {
	const quote = "<blockquote>\n%[1]s\n</blockquote>"
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf(quote, "test")
	}
}

func BenchmarkPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = "<blockquote>\n" + "test" + "\n</blockquote>"
	}
}

func BenchmarkBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var w strings.Builder
		w.WriteString("<blockquote>\n")
		w.WriteString("test")
		w.WriteString("\n</blockquote>")
		_ = w.String()
	}
}
