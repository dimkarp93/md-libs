package markdown

import (
	"reflect"
	"testing"
)

func TestSplitPages(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []string
	}{
		{
			name: "front matter becomes page zero",
			src:  "---\ntitle: t\n---\nbody one\n---\nbody two\n",
			want: []string{"title: t", "body one", "body two"},
		},
		{
			name: "no front matter leaves page zero empty",
			src:  "body one\n---\nbody two\n",
			want: []string{"", "body one", "body two"},
		},
		{
			name: "delimiter inside code fence does not split",
			src:  "before\n```\n---\n```\nafter\n",
			want: []string{"", "before\n```\n---\n```\nafter"},
		},
		{
			name: "unclosed front matter: delimiter falls back to a page break",
			src:  "---\nnot closed\n",
			want: []string{"", "", "not closed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitPages(tt.src)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitPages() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestParsePageRanges(t *testing.T) {
	tests := []struct {
		name    string
		spec    string
		want    []int
		wantErr bool
	}{
		{name: "single and range", spec: "1,3-5", want: []int{1, 3, 4, 5}},
		{name: "dedup and sort", spec: "5,1,3-4,4", want: []int{1, 3, 4, 5}},
		{name: "zero page", spec: "0", want: []int{0}},
		{name: "spaces tolerated", spec: " 2 , 4 ", want: []int{2, 4}},
		{name: "reversed range", spec: "5-1", wantErr: true},
		{name: "not a number", spec: "a", wantErr: true},
		{name: "broken range", spec: "1-x", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePageRanges(tt.spec)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParsePageRanges(%q) = %v, want error", tt.spec, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParsePageRanges(%q) unexpected error: %v", tt.spec, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePageRanges(%q) = %v, want %v", tt.spec, got, tt.want)
			}
		})
	}
}
