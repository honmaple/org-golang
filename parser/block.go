package parser

import (
	"regexp"
	"strings"
)

const (
	BlockName       = "Block"
	BlockResultName = "BlockResult"
)

var (
	beginBlockRegexp         = regexp.MustCompile(`(?i)^(\s*)#\+BEGIN_(\w+)(.*)`)
	endBlockRegexp           = regexp.MustCompile(`(?i)^(\s*)#\+END_(\w+)`)
	resultRegexp             = regexp.MustCompile(`(?i)^(\s*)#\+RESULTS:`)
	exampleBlockEscapeRegexp = regexp.MustCompile(`(^|\n)([ \t]*),([ \t]*)(\*|,\*|#\+|,#\+)`)
)

type Block struct {
	Type       string
	Parameters []string
	Result     *BlockResult
	Children   []Node
}

func (Block) Name() string {
	return BlockName
}

type BlockResult struct {
	Children []Node
}

func (BlockResult) Name() string {
	return BlockResultName
}

func (s *parser) ParseBlock(d *Document, lines []string) (*Block, int) {
	match := beginBlockRegexp.FindStringSubmatch(lines[0])
	if match == nil {
		return nil, 0
	}
	blockType := strings.ToUpper(match[2])

	idx, end := 1, len(lines)
	for idx < end {
		if m := endBlockRegexp.FindStringSubmatch(lines[idx]); m != nil && strings.ToUpper(m[2]) == blockType {
			b := &Block{
				Type:       blockType,
				Parameters: strings.Split(strings.TrimSpace(match[3]), " "),
			}
			switch blockType {
			case "VERSE":
				b.Children = s.ParseAllInline(d, strings.Join(lines[1:idx], "\n"), false)
			case "SRC", "EXAMPLE":
				b.Children = s.ParseAll(d, lines[1:idx], true)
			default:
				b.Children = s.ParseAll(d, lines[1:idx], false)
			}
			return b, idx + 1
		}
		idx++
	}
	return nil, 0
}

func (s *parser) ParseBlockResult(d *Document, lines []string) (*BlockResult, int) {
	match := resultRegexp.FindStringSubmatch(lines[0])
	if match == nil {
		return nil, 0
	}

	idx, end := 1, len(lines)
	for idx < end {
		if match := resultRegexp.FindStringSubmatch(lines[idx]); match == nil {
			return &BlockResult{s.ParseAll(d, lines[1:idx], false)}, idx + 1
		}
		idx++
	}
	return nil, 0
}
