package parser

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	TableName       = "Table"
	TableRowName    = "TableRow"
	TableColumnName = "TableColumn"
)

var (
	tableSepRegexp   = regexp.MustCompile(`^(\s*)(\|[+-|]*)\s*$`)
	tableRowRegexp   = regexp.MustCompile(`^(\s*)(\|.*)`)
	tableAlignRegexp = regexp.MustCompile(`^<(l|c|r)>$`)
)

type Table struct {
	numeric  bool
	Aligns   []string
	Children []Node
}

type TableRow struct {
	Children  []Node
	Separator bool
}

type TableColumn struct {
	IsHeader bool
	Align    string
	Children []Node
}

func (Table) Name() string {
	return TableName
}

func (TableRow) Name() string {
	return TableRowName
}

func (s TableRow) IsAlign() bool {
	return false
}

func (TableColumn) Name() string {
	return TableColumnName
}

func isNumeric(text string) bool {
	if _, err := strconv.Atoi(text); err != nil {
		return false
	}
	return true
}
func repeatString(text string, count int) []string {
	s := make([]string, count)
	for i := 0; i < count; i++ {
		s[i] = text
	}
	return s
}

func (s *parser) TableRow(d *Document, lines []string) (*TableRow, int) {
	line := lines[0]
	match := tableRowRegexp.FindStringSubmatch(line)
	if match == nil {
		return nil, 0
	}
	if tableSepRegexp.MatchString(line) {
		return &TableRow{Separator: true}, 1
	}

	aligns := make([]string, 0)
	children := make([]Node, 0)
	for _, text := range strings.FieldsFunc(match[2], func(r rune) bool { return r == '|' }) {
		text = strings.TrimSpace(text)
		if m := tableAlignRegexp.FindStringSubmatch(text); m != nil {
			aligns = append(aligns, m[1])
		}
		children = append(children, &TableColumn{Children: s.ParseAllInline(d, text, false)})
	}
	if len(aligns) == len(children) {

	}
	return &TableRow{Children: children}, 1
}

func (s *parser) Table(d *Document, lines []string) (*Table, int) {
	rows := make([]Node, 0)

	idx, end := 0, len(lines)
	for idx < end {
		row, rowIdx := s.TableRow(d, lines[idx:])
		if row == nil {
			break
		}
		idx = idx + rowIdx
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil, 0
	}
	b := &Table{
		numeric:  true,
		Children: rows,
	}
	return b, idx
}
