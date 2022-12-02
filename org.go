package org

import (
	"io"

	"github.com/honmaple/org-golang/parser"
	"github.com/honmaple/org-golang/render"
)

const (
	todoKeywords    = "TODO | DONE | CANCELED"
	timestampFormat = "2006-01-02 Mon 15:04"
)

type Option func(*parser.Document)

func New(r io.Reader, opts ...Option) *parser.Document {
	d := &parser.Document{
		Sections: &parser.Section{},
		Keywords: map[string]string{
			"TODO": todoKeywords,
		},
		Hyperlinks:      []string{"http", "https", "file"},
		TimestampFormat: timestampFormat,
	}
	for _, opt := range opts {
		opt(d)
	}
	d.Children = parser.Parse(d, r)
	return d
}

func HTML(r io.Reader, opts ...Option) string {
	out := render.HTML{
		Document: New(r, opts...),
	}
	return out.String()
}

func Debug(r io.Reader, opts ...Option) string {
	out := render.Debug{
		Document: New(r, opts...),
	}
	return out.String()
}
