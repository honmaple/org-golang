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
	Children []Node
}

type TableRow struct {
	Children  []Node
	Separator bool
	Infos     []string
}

type TableColumn struct {
	Align    string
	Width    int
	IsHeader bool
	Children []Node
}

func (Table) Name() string {
	return TableName
}

func (TableRow) Name() string {
	return TableRowName
}

func (TableColumn) Name() string {
	return TableColumnName
}

func (s *parser) ParseTableRow(d *Document, lines []string) (*TableRow, int) {
	line := lines[0]
	match := tableRowRegexp.FindStringSubmatch(line)
	if match == nil {
		return nil, 0
	}
	if tableSepRegexp.MatchString(line) {
		return &TableRow{Separator: true}, 1
	}

	infos := make([]string, 0)
	texts := make([]string, 0)
	for _, text := range strings.FieldsFunc(match[2], func(r rune) bool { return r == '|' }) {
		text = strings.TrimSpace(text)
		texts = append(texts, text)
		if m := tableAlignRegexp.FindStringSubmatch(text); m != nil {
			infos = append(infos, m[1])
		}
	}
	// if not equal, infos is not infos, just tablecolumn
	if len(infos) == len(texts) {
		return &TableRow{Infos: infos}, 1
	}
	children := make([]Node, len(texts))
	for i, text := range texts {
		children[i] = &TableColumn{Children: s.ParseAllInline(d, text, false)}
	}
	return &TableRow{Children: children}, 1
}

func (s *parser) ParseTable(d *Document, lines []string) (*Table, int) {
	var (
		rows   = make([]Node, 0)
		infos  []string
		header int
	)

	idx, end := 0, len(lines)
	for idx < end {
		row, rowIdx := s.ParseTableRow(d, lines[idx:])
		if row == nil {
			break
		}
		idx = idx + rowIdx
		if header == 0 && row.Separator {
			header = len(rows)
		}
		if len(row.Infos) > 0 {
			infos = row.Infos
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil, 0
	}
	for i, info := range infos {
		align := ""
		width := 0
		switch info {
		case "l":
			align = "left"
		case "r":
			align = "right"
		case "c":
			align = "center"
		}
		if n, err := strconv.Atoi(info); err == nil {
			width = n
		}
		for _, node := range rows {
			row := node.(*TableRow)
			if row.Separator || len(row.Children) <= i {
				continue
			}
			column := row.Children[i].(*TableColumn)
			column.Align = align
			column.Width = width
		}
	}
	for i, node := range rows[:header] {
		row := node.(*TableRow)
		if row.Separator || len(row.Children) == 0 {
			continue
		}
		for _, column := range row.Children {
			column.(*TableColumn).IsHeader = i < header
		}
	}
	b := &Table{
		Children: rows,
	}
	return b, idx
}
