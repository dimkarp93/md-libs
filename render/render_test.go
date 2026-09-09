package render

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dimkarp93/md-libs/markdown"
)

type recorder struct {
	calls []string
}

func (r *recorder) Heading(level int, text string) {
	r.calls = append(r.calls, fmt.Sprintf("heading(%d,%s)", level, text))
}

func (r *recorder) Paragraph(text string) {
	r.calls = append(r.calls, fmt.Sprintf("paragraph(%s)", text))
}

func (r *recorder) Code(lines []string) {
	r.calls = append(r.calls, fmt.Sprintf("code(%v)", lines))
}

func (r *recorder) PageBreak() {
	r.calls = append(r.calls, "pagebreak()")
}

func TestBlocksDispatchesEveryKindWithItsPayload(t *testing.T) {
	blocks := []markdown.Block{
		{Kind: markdown.Heading, Level: 3, Text: "Result"},
		{Kind: markdown.Paragraph, Text: "line one\nline two"},
		{Kind: markdown.Code, Lines: []string{"a", "b"}},
		{Kind: markdown.PageBreak},
	}

	r := &recorder{}
	Blocks(r, blocks)

	want := []string{
		"heading(3,Result)",
		"paragraph(line one\nline two)",
		"code([a b])",
		"pagebreak()",
	}
	if !reflect.DeepEqual(r.calls, want) {
		t.Errorf("Blocks() = %#v, want %#v", r.calls, want)
	}
}

func TestBlocksOnEmptyInput(t *testing.T) {
	r := &recorder{}
	Blocks(r, nil)
	if len(r.calls) != 0 {
		t.Errorf("Blocks(nil) made %d calls, want none", len(r.calls))
	}
}

func TestBlocksTreatsUnknownKindAsParagraph(t *testing.T) {
	r := &recorder{}
	Blocks(r, []markdown.Block{{Kind: markdown.BlockKind(42), Text: "fallback"}})

	want := []string{"paragraph(fallback)"}
	if !reflect.DeepEqual(r.calls, want) {
		t.Errorf("Blocks() = %#v, want %#v", r.calls, want)
	}
}

func TestHeadingSizePt(t *testing.T) {
	want := map[int]float64{1: 18, 2: 16, 3: 14, 4: 13, 5: 12, 6: 11}
	for level, size := range want {
		if got := HeadingSizePt(level); got != size {
			t.Errorf("HeadingSizePt(%d) = %v, want %v", level, got, size)
		}
	}
}

func TestHeadingSizeDecreasesWithLevel(t *testing.T) {
	for level := 1; level < 6; level++ {
		if HeadingSizePt(level) < HeadingSizePt(level+1) {
			t.Errorf("HeadingSizePt(%d) = %v is smaller than HeadingSizePt(%d) = %v",
				level, HeadingSizePt(level), level+1, HeadingSizePt(level+1))
		}
	}
}
