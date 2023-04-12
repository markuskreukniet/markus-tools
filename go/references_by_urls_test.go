package main

import (
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func TestReferencesByURLs(t *testing.T) {
	tests := []struct {
		htmlDocument   string
		expectedTitles []string
		expectedH1s    []string
	}{
		{
			htmlDocument:   "<html><head><title>Test Title</title></head><body><h1>Heading 1</h1></body></html>",
			expectedTitles: []string{"<title>Test Title</title>"},
			expectedH1s:    []string{"<h1>Heading 1</h1>"},
		},
		{
			htmlDocument:   "<html><head><title>Another Title</title></head><body><h1>First Heading</h1><h1>Second Heading</h1></body></html>",
			expectedTitles: []string{"<title>Another Title</title>"},
			expectedH1s:    []string{"<h1>First Heading</h1>", "<h1>Second Heading</h1>"},
		},
		{
			htmlDocument:   "<html><head></head><body><h1>Only Heading</h1></body></html>",
			expectedTitles: nil,
			expectedH1s:    []string{"<h1>Only Heading</h1>"},
		},
		{
			htmlDocument:   "<html><head><title>Empty Title</title></head><body><h1></h1></body></html>",
			expectedTitles: []string{"<title>Empty Title</title>"},
			expectedH1s:    []string{"<h1></h1>"},
		},
		{
			htmlDocument:   "<html><head><title>    </title></head><body><h1>     </h1></body></html>",
			expectedTitles: []string{"<title>    </title>"},
			expectedH1s:    []string{"<h1>     </h1>"},
		},
		{
			htmlDocument:   "<html><head><title te_s-t   test2-lol   test3=\"asdf-asdf2\">    </title ></head><body><h1 asdf asdf=\"asdf89-asdf99\"   >     </h1     ></body></html>",
			expectedTitles: []string{"<title te_s-t   test2-lol   test3=\"asdf-asdf2\">    </title >"},
			expectedH1s:    []string{"<h1 asdf asdf=\"asdf89-asdf99\"   >     </h1     >"},
		},
		{
			htmlDocument:   "<html><head><title/></head><body><h1   /> <h1   test-a=\"a_b-c\"  /></body></html>",
			expectedTitles: []string{"<title/>"},
			expectedH1s:    []string{"<h1   />", "<h1   test-a=\"a_b-c\"  />"},
		},
		{
			htmlDocument:   "<html><head><title te-st=\"a---test\"  lol  ><title/></title   ></head><body></body></html>",
			expectedTitles: []string{"<title te-st=\"a---test\"  lol  ><title/></title   >"},
			expectedH1s:    nil,
		},
		{
			htmlDocument:   "<html><head></head><body><h1 te-st=\"a---test\"  lol  ><h1 asdf-l=\"test\"  /></h1   ><h1/></body></html>",
			expectedTitles: nil,
			expectedH1s:    []string{"<h1 te-st=\"a---test\"  lol  ><h1 asdf-l=\"test\"  /></h1   >", "<h1/>"},
		},
		{
			htmlDocument:   "<html><head></head><body><h1><h1></h1></h1></body></html>",
			expectedTitles: nil,
			expectedH1s:    []string{"<h1><h1></h1></h1>"},
		},
	}

	for _, test := range tests {
		titles, h1s := findTitleAndH1Elements(test.htmlDocument)

		utils.TMustAssertDeepEqual(t, test.expectedTitles, titles)
		utils.TMustAssertDeepEqual(t, test.expectedH1s, h1s)
	}
}

func TestFilterComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No comments",
			input:    "<html><body>Hello World</body></html>",
			expected: "<html><body>Hello World</body></html>",
		},
		{
			name:     "HTML comment",
			input:    "<html><!-- This is a comment --><body>Hello World</body></html>",
			expected: "<html><body>Hello World</body></html>",
		},
		{
			name:     "HTML empty comment",
			input:    "<html><!----><body>Hello World</body></html>",
			expected: "<html><body>Hello World</body></html>",
		},
		{
			name:     "Single-line JS comment",
			input:    "<script>// This is a comment\nvar x = 1;</script>",
			expected: "<script>var x = 1;</script>",
		},
		{
			name:     "Single-line JS empty comment",
			input:    "<script>//\nvar x = 1;</script>",
			expected: "<script>var x = 1;</script>",
		},
		{
			name:     "Multi-line JS comment",
			input:    "<script>/* This is a \n multi-line comment */var x = 1;</script>",
			expected: "<script>var x = 1;</script>",
		},
		{
			name:     "Multi-line empty comment",
			input:    "<script>/**/var x = 1;</script>",
			expected: "<script>var x = 1;</script>",
		},
		{
			name:     "Mixed comments",
			input:    "<html><!-- HTML comment --><script>// JS single-line comment\n/* JS multi-line comment */var x = 1;</script></html>",
			expected: "<html><script>var x = 1;</script></html>",
		},
		{
			name:     "Escaped characters",
			input:    "<script>var str = \"This is not a comment: \\\" /* not a comment */\";</script>",
			expected: "<script>var str = \"This is not a comment: \\\" /* not a comment */\";</script>",
		},
		{
			name:     "Comment with escaped newline",
			input:    "<script>// JS comment \\n still comment\nvar x = 1;</script>",
			expected: "<script>var x = 1;</script>",
		},
		{
			name:     "Single-line JS comment in backtick string interpolation",
			input:    "<script>let test = `t ${A}, asdf ${C // A comment} another test.`;</script>",
			expected: "<script>let test = `t ${A}, asdf ${C } another test.`;</script>",
		},
		{
			name:     "Multi-line JS comment in backtick string interpolation",
			input:    "<script>let test = `t ${A}, asdf ${C /*A comment*/ a/* A comment 2*/} another test.`;</script>",
			expected: "<script>let test = `t ${A}, asdf ${C  a} another test.`;</script>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.TMustAssertEqualStrings(t, tt.expected, filterComments(tt.input))
		})
	}
}

func TestRemoveTagsFromElement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<div>Hello</div>", "Hello"},
		{"<p>This is a <strong>test</strong>.</p>", "This is a test."},
		{"<a href=\"#\">Link</a>", "Link"},
		{"<span class='class-name'>Text</span>", "Text"},
		{"<img src=\"image.jpg\" alt=\"image\"/>", ""},
		{"<div>Nested <span>tags</span> example</div>", "Nested tags example"},
		{"No tags here", "No tags here"},
		{"<div>Incomplete tag", "Incomplete tag"},
	}

	for _, tt := range tests {
		utils.TMustAssertEqualStrings(t, tt.expected, removeTagsFromElement(tt.input))
	}
}

func TestRemoveTagsFromElements(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			[]string{"<div>Hello</div>", "<p>World</p>"},
			[]string{"Hello", "World"},
		},
		{
			[]string{"<a href=\"#\">Link 1</a>", "<a href=\"#\">Link 2</a>"},
			[]string{"Link 1", "Link 2"},
		},
		{
			[]string{"<span>Text</span>", "Plain text"},
			[]string{"Text", "Plain text"},
		},
	}

	for _, tt := range tests {
		utils.TMustAssertDeepEqual(t, tt.expected, removeTagsFromElements(tt.input))
	}
}
