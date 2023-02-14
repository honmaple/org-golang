package render

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTML(t *testing.T) {
	files := testFiles()
	for _, file := range files {
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		expect, err := ioutil.ReadFile(file[:len(file)-4] + ".html")
		if err != nil {
			panic(err)
		}
		out := HTML{
			Document: toDocument(buf),
			Toc:      true,
		}
		assert.Equal(t, string(expect), out.String())
	}
}
