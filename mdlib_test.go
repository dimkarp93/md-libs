package mdlib

import (
	"reflect"
	"testing"

	"github.com/dimkarp93/md-libs/markdown"
)

const src = "---\nfront: matter\n---\n# One\n\nfirst page\n\n---\n\n## Result\n\nsecond page\n\n---\n\n# Three\n\nthird page\n"

func TestPrepareDefaultsToAllPagesExceptZero(t *testing.T) {
	got, err := Prepare(src, Options{})
	if err != nil {
		t.Fatalf("Prepare() unexpected error: %v", err)
	}
	want := []markdown.Block{
		{Kind: markdown.Heading, Level: 1, Text: "One"},
		{Kind: markdown.Paragraph, Text: "first page"},
		{Kind: markdown.PageBreak},
		{Kind: markdown.Heading, Level: 2, Text: "Result"},
		{Kind: markdown.Paragraph, Text: "second page"},
		{Kind: markdown.PageBreak},
		{Kind: markdown.Heading, Level: 1, Text: "Three"},
		{Kind: markdown.Paragraph, Text: "third page"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Prepare() = %#v, want %#v", got, want)
	}
}

func TestPrepareSelectsPages(t *testing.T) {
	got, err := Prepare(src, Options{Pages: "0,3"})
	if err != nil {
		t.Fatalf("Prepare() unexpected error: %v", err)
	}
	want := []markdown.Block{
		{Kind: markdown.Paragraph, Text: "front: matter"},
		{Kind: markdown.PageBreak},
		{Kind: markdown.Heading, Level: 1, Text: "Three"},
		{Kind: markdown.Paragraph, Text: "third page"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Prepare() = %#v, want %#v", got, want)
	}
}

func TestPrepareIgnoresOutOfRangePages(t *testing.T) {
	got, err := Prepare(src, Options{Pages: "2,99"})
	if err != nil {
		t.Fatalf("Prepare() unexpected error: %v", err)
	}
	want := []markdown.Block{
		{Kind: markdown.Heading, Level: 2, Text: "Result"},
		{Kind: markdown.Paragraph, Text: "second page"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Prepare() = %#v, want %#v", got, want)
	}
}

func TestPrepareAppliesHeadFilters(t *testing.T) {
	got, err := Prepare(src, Options{Heads: "h2:result", HideMatchedHeads: true})
	if err != nil {
		t.Fatalf("Prepare() unexpected error: %v", err)
	}
	want := []markdown.Block{
		{Kind: markdown.Paragraph, Text: "second page"},
		{Kind: markdown.PageBreak},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Prepare() = %#v, want %#v", got, want)
	}
}

func TestPrepareRejectsBadPages(t *testing.T) {
	if _, err := Prepare(src, Options{Pages: "5-1"}); err == nil {
		t.Error("Prepare() = nil error, want error for reversed range")
	}
}
