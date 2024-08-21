package main

import (
	"reflect"
	"testing"
	"unicode"
)

func isLetter(r rune) bool {
	return unicode.IsUpper(r) || unicode.IsLower(r)
}

func isLetterDigitHyphenOrUnderscore(r rune) bool {
	return isLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
}

// returns: tagPartLength, tagIsClosed, tagIsFound
func finishCreatingStartTag(document []rune, documentLength, index int) (int, bool, bool) {
	startTagEndPartLength := 0
	var quoteRune rune
	inAttributeName := false
	inAttributeValue := false

	for ; index < documentLength; index++ {
		switch {
		case inAttributeName:
			if document[index] == '=' {
				startTagEndPartLength++
				inAttributeName = false
				inAttributeValue = true
			} else if isLetterDigitHyphenOrUnderscore(document[index]) {
				startTagEndPartLength++
			} else if unicode.IsSpace(document[index]) {
				startTagEndPartLength++
				inAttributeName = false
			} else {
				return 0, false, false
			}
		case inAttributeValue:
			if document[index] == '"' || document[index] == '\'' {
				quoteRune = document[index]
				startTagEndPartLength++
				for i := index + 1; i < documentLength; i++ {
					startTagEndPartLength++
					if document[i] == quoteRune {
						index = i
						break
					}
				}
				inAttributeValue = false
			} else {
				return 0, false, false
			}
		default:
			if isLetter(document[index]) {
				startTagEndPartLength++
				inAttributeName = true
			} else if document[index] == '>' {
				startTagEndPartLength++
				return startTagEndPartLength, false, true
			} else if hasPrefix, length := hasStringPrefix(document, index, "/>"); hasPrefix {
				startTagEndPartLength += length
				return startTagEndPartLength, true, true
			} else if unicode.IsSpace(document[index]) {
				startTagEndPartLength++
			} else {
				return 0, false, false
			}
		}
	}

	return 0, false, false
}

// TODO: use tagIsFound
// returns: htmlElementPartLength, htmlElementIsFound
func getTheOtherHTMLElementPartLength(document, startTagPart, endTagPart []rune, index, documentLength, startTagPartLength, endTagPartLength int) (int, bool) {
	tagPartLength := 0
	numberOfOpenStartTags := 1

	closingTagPart := []rune("/>") // TODO: /> is duplicate
	closingTagPartLength := len(closingTagPart)

	for i := index; i < documentLength; i++ {
		if hasPrefix(document, startTagPart, documentLength, startTagPartLength, i) {
			tagPartLength += startTagPartLength
			i += startTagPartLength
			length, tagIsClosed, _ := finishCreatingStartTag(document, documentLength, i)
			tagPartLength += length
			i += length - 1
			if !tagIsClosed {
				numberOfOpenStartTags++
			}
		} else if hasPrefix(document, endTagPart, documentLength, endTagPartLength, i) {
			tagPartLength += endTagPartLength
			i += endTagPartLength
			for ; i < documentLength; i++ {
				tagPartLength++
				if document[i] == '>' {
					numberOfOpenStartTags--
					if numberOfOpenStartTags == 0 {
						return tagPartLength, true
					}
				}
			}
		} else if hasPrefix(document, closingTagPart, documentLength, closingTagPartLength, i) {
			tagPartLength += closingTagPartLength
			i += closingTagPartLength - 1
			numberOfOpenStartTags--
			if numberOfOpenStartTags == 0 {
				return tagPartLength, true
			}
		} else {
			tagPartLength++
		}
	}

	return 0, false
}

func hasOpenOrSelfClosingHTMLTagPrefix(document, prefix []rune, documentLength, prefixLength, index int) (int, bool, bool) {
	if hasPrefix(document, prefix, documentLength, prefixLength, index) {
		if tagPartLength, tagIsClosed, tagIsFound := finishCreatingStartTag(document, documentLength, index+prefixLength); tagIsFound {
			return prefixLength + tagPartLength, tagIsClosed, true
		}
	}

	return 0, false, false
}

// Finding HTML elements should happen for every element name in a complete HTML document since an element could be a child element of another element.
func findHTMLElements(document, elementName string) []string {
	var elements []string

	documentRunes := []rune(document)
	startTagPartRunes := append([]rune("<"), []rune(elementName)...)
	endTagPartRunes := append([]rune("</"), []rune(elementName)...)

	documentLength := len(documentRunes)
	startTagPartLength := len(startTagPartRunes)

	for i := 0; i < documentLength; i++ {
		length, tagIsClosed, hasPrefix := hasOpenOrSelfClosingHTMLTagPrefix(documentRunes, startTagPartRunes, documentLength, startTagPartLength, i)
		if hasPrefix {
			elementPart := documentRunes[i : i+length]
			i += length
			if !tagIsClosed {
				// TODO: use htmlElementIsFound?
				elementPartLength, _ := getTheOtherHTMLElementPartLength(documentRunes, startTagPartRunes, endTagPartRunes, i, documentLength, startTagPartLength, len(endTagPartRunes))
				elementPart = append(elementPart, documentRunes[i:i+elementPartLength]...)
				i += elementPartLength
			}
			elements = append(elements, string(elementPart))
			i--
		}
	}

	return elements
}

func findTitleAndH1Elements(htmlDocument string) ([]string, []string) {
	return findHTMLElements(htmlDocument, "title"), findHTMLElements(htmlDocument, "h1")
}

// TODO: not efficient since it is used in a loop
func hasStringPrefix(runes []rune, index int, prefix string) (bool, int) {
	prefixLength := len(prefix)

	if len(runes)-index < prefixLength {
		return false, 0
	}

	for _, r := range prefix {
		if runes[index] != r {
			return false, 0
		}
		index++
	}

	return true, prefixLength
}

// // //

func hasPrefix(runes, prefix []rune, runesLength, prefixLength, index int) bool {
	if runesLength-index < prefixLength {
		return false
	}

	for _, r := range prefix {
		if runes[index] != r {
			return false
		}
		index++
	}

	return true
}

func updateIndexMinusOneIfHasPrefix(runes, prefix []rune, runesLength, prefixLength int, index *int) bool {
	if hasPrefix(runes, prefix, runesLength, prefixLength, *index) {
		*index += prefixLength - 1
		return true
	}

	return false
}

func updateIndexIfHasPrefix(runes, prefix []rune, runesLength, prefixLength int, index *int) bool {
	if hasPrefix(runes, prefix, runesLength, prefixLength, *index) {
		*index += prefixLength
		return true
	}

	return false
}

func appendIfEscape(htmlDocument []rune, filteredHTMLDocument *[]rune, index, htmlDocumentLength int) bool {
	indexPlusOne := index + 1

	if htmlDocument[index] == '\\' && indexPlusOne < htmlDocumentLength {
		*filteredHTMLDocument = append(*filteredHTMLDocument, htmlDocument[index])
		*filteredHTMLDocument = append(*filteredHTMLDocument, htmlDocument[indexPlusOne])
		return true
	}

	return false
}

type commentDelimiter struct {
	startDelimiter       []rune
	endDelimiter         []rune
	startDelimiterLength int
	endDelimiterLength   int
}

func createCommentDelimiter(startDelimiter, endDelimiter []rune) commentDelimiter {
	return commentDelimiter{
		startDelimiter:       startDelimiter,
		endDelimiter:         endDelimiter,
		startDelimiterLength: len(startDelimiter),
		endDelimiterLength:   len(endDelimiter),
	}
}

func updateIndexIfComment(htmlDocumentRunes []rune, htmlDocumentRunesLength int, index *int, commentDelimiters []commentDelimiter) bool {
	for _, delimiter := range commentDelimiters {
		if updateIndexIfHasPrefix(htmlDocumentRunes, delimiter.startDelimiter, htmlDocumentRunesLength, delimiter.startDelimiterLength, index) {
			for ; *index < htmlDocumentRunesLength; *index++ {
				if updateIndexMinusOneIfHasPrefix(htmlDocumentRunes, delimiter.endDelimiter, htmlDocumentRunesLength, delimiter.endDelimiterLength, index) {
					return true
				}
			}
		}
	}

	return false
}

func filterComments(htmlDocument string) string {
	var filteredHTMLDocument []rune
	var jsStringRune rune

	htmlDocumentRunes := []rune(htmlDocument)
	htmlCommentStart := []rune("<!--")
	htmlCommentEnd := []rune("-->")
	jsCommentSingleLineStart := []rune("//")
	jsCommentSingleLineEnd := []rune("\n")
	commentMultiLineStart := []rune("/*")
	commentMultiLineEnd := []rune("*/")

	htmlDocumentRunesLength := len(htmlDocumentRunes)

	commentDelimiters := []commentDelimiter{
		createCommentDelimiter(htmlCommentStart, htmlCommentEnd),
		createCommentDelimiter(jsCommentSingleLineStart, jsCommentSingleLineEnd),
		createCommentDelimiter(commentMultiLineStart, commentMultiLineEnd),
	}

	for i := 0; i < htmlDocumentRunesLength; i++ {
		// string escape
		if appendIfEscape(htmlDocumentRunes, &filteredHTMLDocument, i, htmlDocumentRunesLength) {
			i++
		} else if updateIndexIfComment(htmlDocumentRunes, htmlDocumentRunesLength, &i, commentDelimiters) {
			continue
			// JavaScript string
		} else if htmlDocumentRunes[i] == '"' || htmlDocumentRunes[i] == '\'' { // TODO: also add backtick strings
			filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
			jsStringRune = htmlDocumentRunes[i]
			i++
			for ; i < htmlDocumentRunesLength; i++ {
				if appendIfEscape(htmlDocumentRunes, &filteredHTMLDocument, i, htmlDocumentRunesLength) {
					i++
				} else {
					filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
					if htmlDocumentRunes[i] == jsStringRune {
						break
					}
				}
			}
		} else {
			filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
		}
	}

	return string(filteredHTMLDocument)
}

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

		if !reflect.DeepEqual(titles, test.expectedTitles) {
			t.Errorf("findTitleAndH1Elements(%q) titles = %v; want %v", test.htmlDocument, titles, test.expectedTitles)
		}

		if !reflect.DeepEqual(h1s, test.expectedH1s) {
			t.Errorf("findTitleAndH1Elements(%q) h1s = %v; want %v", test.htmlDocument, h1s, test.expectedH1s)
		}
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := filterComments(tt.input)
			if actual != tt.expected {
				t.Errorf("filterComments() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
