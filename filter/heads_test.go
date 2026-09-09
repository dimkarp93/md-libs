package filter

import (
	"reflect"
	"testing"

	"github.com/dimkarp93/md-libs/markdown"
)

func TestParseHeads(t *testing.T) {
	tests := []struct {
		name string
		spec string
		want []HeadFilter
	}{
		{
			name: "levelled entries",
			spec: "h2:result,h3:resume",
			want: []HeadFilter{{Level: 2, Name: "result"}, {Level: 3, Name: "resume"}},
		},
		{
			name: "bare name matches any level",
			spec: "result",
			want: []HeadFilter{{Level: 0, Name: "result"}},
		},
		{
			name: "case of the name is preserved",
			spec: "H2:Result",
			want: []HeadFilter{{Level: 2, Name: "Result"}},
		},
		{
			name: "empty entries skipped",
			spec: "h1:a,,h2:b",
			want: []HeadFilter{{Level: 1, Name: "a"}, {Level: 2, Name: "b"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHeads(tt.spec)
			if err != nil {
				t.Fatalf("ParseHeads(%q) unexpected error: %v", tt.spec, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseHeads(%q) = %#v, want %#v", tt.spec, got, tt.want)
			}
		})
	}
}

func TestByHeads(t *testing.T) {
	blocks := []markdown.Block{
		{Kind: markdown.Heading, Level: 1, Text: "Intro"},
		{Kind: markdown.Paragraph, Text: "intro text"},
		{Kind: markdown.Heading, Level: 2, Text: "Result"},
		{Kind: markdown.Paragraph, Text: "result text"},
		{Kind: markdown.Heading, Level: 3, Text: "Nested"},
		{Kind: markdown.Paragraph, Text: "nested text"},
		{Kind: markdown.Heading, Level: 2, Text: "Other"},
		{Kind: markdown.Paragraph, Text: "other text"},
	}

	tests := []struct {
		name        string
		filters     []HeadFilter
		hideMatched bool
		want        []markdown.Block
	}{
		{
			name:    "section keeps nested content until same level",
			filters: []HeadFilter{{Level: 2, Name: "result"}},
			want: []markdown.Block{
				{Kind: markdown.Heading, Level: 2, Text: "Result"},
				{Kind: markdown.Paragraph, Text: "result text"},
				{Kind: markdown.Heading, Level: 3, Text: "Nested"},
				{Kind: markdown.Paragraph, Text: "nested text"},
			},
		},
		{
			name:        "hideMatched drops the heading itself",
			filters:     []HeadFilter{{Level: 2, Name: "result"}},
			hideMatched: true,
			want: []markdown.Block{
				{Kind: markdown.Paragraph, Text: "result text"},
				{Kind: markdown.Heading, Level: 3, Text: "Nested"},
				{Kind: markdown.Paragraph, Text: "nested text"},
			},
		},
		{
			name:    "level mismatch selects nothing",
			filters: []HeadFilter{{Level: 1, Name: "result"}},
			want:    nil,
		},
		{
			name:    "bare filter matches any level",
			filters: []HeadFilter{{Level: 0, Name: "nested"}},
			want: []markdown.Block{
				{Kind: markdown.Heading, Level: 3, Text: "Nested"},
				{Kind: markdown.Paragraph, Text: "nested text"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ByHeads(blocks, tt.filters, tt.hideMatched)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByHeads() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestMatchesIgnoresInlineMarkup(t *testing.T) {
	b := markdown.Block{Kind: markdown.Heading, Level: 2, Text: "**Result**"}
	if !Matches(b, []HeadFilter{{Level: 2, Name: "result"}}) {
		t.Error("Matches() = false, want true for bold heading matched case-insensitively")
	}
}
