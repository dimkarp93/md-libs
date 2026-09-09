package markdown

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []Block
	}{
		{
			name: "headings and paragraphs",
			src:  "# One\n\ntext a\ntext b\n\n### Three\n",
			want: []Block{
				{Kind: Heading, Level: 1, Text: "One"},
				{Kind: Paragraph, Text: "text a\ntext b"},
				{Kind: Heading, Level: 3, Text: "Three"},
			},
		},
		{
			name: "fenced code keeps blank lines and hashes",
			src:  "para\n\n```go\nfunc f() {\n\n# not a heading\n}\n```\n\nafter\n",
			want: []Block{
				{Kind: Paragraph, Text: "para"},
				{Kind: Code, Lines: []string{"func f() {", "", "# not a heading", "}"}},
				{Kind: Paragraph, Text: "after"},
			},
		},
		{
			name: "unterminated code fence still flushes",
			src:  "```\nx\n",
			want: []Block{
				{Kind: Code, Lines: []string{"x"}},
			},
		},
		{
			name: "heading requires space after hashes",
			src:  "#nospace\n",
			want: []Block{
				{Kind: Paragraph, Text: "#nospace"},
			},
		},
		{
			name: "empty input",
			src:  "",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.src)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestParseInline(t *testing.T) {
	tests := []struct {
		name            string
		text            string
		allowUnderscore bool
		want            []Run
	}{
		{
			name: "bold italic combo",
			text: "a ***bi*** b **b** c *i* d",
			want: []Run{
				{Text: "a "},
				{Text: "bi", Bold: true, Italic: true},
				{Text: " b "},
				{Text: "b", Bold: true},
				{Text: " c "},
				{Text: "i", Italic: true},
				{Text: " d"},
			},
		},
		{
			name:            "underscores honored when allowed",
			text:            "__b__ and _i_",
			allowUnderscore: true,
			want: []Run{
				{Text: "b", Bold: true},
				{Text: " and "},
				{Text: "i", Italic: true},
			},
		},
		{
			name:            "underscores literal when disallowed",
			text:            "__b__ and _i_",
			allowUnderscore: false,
			want: []Run{
				{Text: "__b__ and _i_"},
			},
		},
		{
			name: "plain text",
			text: "nothing special",
			want: []Run{{Text: "nothing special"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseInline(tt.text, tt.allowUnderscore)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseInline() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestPlainText(t *testing.T) {
	if got, want := PlainText("**Result** and *more*"), "Result and more"; got != want {
		t.Errorf("PlainText() = %q, want %q", got, want)
	}
	if got, want := PlainText("_kept_"), "_kept_"; got != want {
		t.Errorf("PlainText() = %q, want %q", got, want)
	}
}
