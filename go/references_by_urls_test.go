package main

import (
	"reflect"
	"testing"
	"unicode"
)

func isLetter(r rune) bool {
	return unicode.IsUpper(r) || unicode.IsLower(r)
}

func isLetterOrUnderscore(r rune) bool {
	return isLetter(r) || r == '_'
}

func isLetterDigitHyphenOrUnderscore(r rune) bool {
	return isLetterOrUnderscore(r) || unicode.IsDigit(r) || r == '-'
}

// returns: tagPartLength, tagIsClosed, tagIsFound
func finishCreatingStartTag(document []rune, documentLength, index int) (int, bool, bool) {
	startTagEndPartLength := 0
	closingTagPart := []rune("/>")
	var quoteRune rune
	inAttributeName := false
	inAttributeValue := false

	closingTagPartLength := len(closingTagPart)

	for ; index < documentLength; index++ {
		if inAttributeName {
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
		} else if inAttributeValue {
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
		} else {
			if isLetter(document[index]) {
				startTagEndPartLength++
				inAttributeName = true
			} else if document[index] == '>' {
				startTagEndPartLength++
				return startTagEndPartLength, false, true
			} else if hasPrefix(document, closingTagPart, documentLength, closingTagPartLength, index) {
				startTagEndPartLength += closingTagPartLength
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

// returns: elementPartLength, elementIsFound
func getTheOtherHTMLElementPartLength(document, startTagPart, endTagPart []rune, documentLength, startTagPartLength, endTagPartLength, index int) (int, bool) {
	tagPartLength := 0
	numberOfOpenStartTags := 1

	for i := index; i < documentLength; i++ {
		if updateTagPartLengthAndIndexIfHasPrefix(document, startTagPart, documentLength, startTagPartLength, &tagPartLength, &i) {
			if length, tagIsClosed, tagIsFound := finishCreatingStartTag(document, documentLength, i); tagIsFound {
				tagPartLength += length
				i += length - 1
				if !tagIsClosed {
					numberOfOpenStartTags++
				}
			} else {
				return 0, false
			}
		} else if updateTagPartLengthAndIndexIfHasPrefix(document, endTagPart, documentLength, endTagPartLength, &tagPartLength, &i) {
			for ; i < documentLength; i++ {
				tagPartLength++
				if document[i] == '>' {
					numberOfOpenStartTags--
					if numberOfOpenStartTags == 0 {
						return tagPartLength, true
					}
				}
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

// TODO: Finding HTML elements should happen for every element name in a complete HTML document since an element could be a child element of another element,
// which is possible to fix.
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
				elementPartLength, elementIsFound := getTheOtherHTMLElementPartLength(documentRunes, startTagPartRunes, endTagPartRunes, documentLength, startTagPartLength, len(endTagPartRunes), i)
				if elementIsFound {
					elementPart = append(elementPart, documentRunes[i:i+elementPartLength]...)
					i += elementPartLength
				} else {
					return nil
				}
			}
			elements = append(elements, string(elementPart))
			i--
		}
	}

	return elements
}
func removeTagsFromElements(elements []string) []string {
	var elementsWithoutTags []string

	for _, element := range elements {
		elementWithoutTags := removeTagsFromElement(element)
		elementsWithoutTags = append(elementsWithoutTags, elementWithoutTags)
	}

	return elementsWithoutTags
}

// TODO: function is useless?
// TODO: WIP
func getTagLength(elementPart []rune, elementPartLength, index int) (int, bool) {
	tagLength := 2

	if elementPart[index] != '<' || !isLetterOrUnderscore(elementPart[index+1]) {
		return 0, false
	}

	index += 2

	for ; index < elementPartLength; index++ {

		// 	for ; i < elementRunesLength; i++ {
		// 		if elementRunes[i] == '<' && elementRunes[i+1] == '/' && isLetterOrUnderscore(elementRunes[i+2]) {
		// 			i += 3
		// 			for ; i < elementRunesLength; i++ {

		// 			}
		// 		} else if unicode.IsSpace(elementRunes[i]) {
		// 			i++
		// 			for ; i < elementRunesLength; i++ {

		// 			}
		// 		}
		// 	}
	}

	return tagLength, false
}

func removeTagsFromElement(element string) string {
	var elementWithoutTags []rune
	elementRunes := []rune(element)
	elementRunesLength := len(elementRunes)

	for i := 0; i < elementRunesLength; i++ {
		if tagLength, foundTag := getTagLength(elementRunes, elementRunesLength, i); foundTag {
			i += tagLength
		} else {
			elementWithoutTags = append(elementWithoutTags, elementRunes[i])
		}
	}

	return string(elementWithoutTags)
}

func findTitleAndH1ElementsAndRemoveTags(htmlDocument string) ([]string, []string) {
	titles, h1s := findTitleAndH1Elements(htmlDocument)
	return removeTagsFromElements(titles), removeTagsFromElements(h1s)
}

func findTitleAndH1Elements(htmlDocument string) ([]string, []string) {
	return findHTMLElements(htmlDocument, "title"), findHTMLElements(htmlDocument, "h1")
}

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

func updateTagPartLengthAndIndexIfHasPrefix(runes, prefix []rune, runesLength, prefixLength int, tagPartLength, index *int) bool {
	if updateIndexIfHasPrefix(runes, prefix, runesLength, prefixLength, index) {
		*tagPartLength += prefixLength
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

func appendAndIncrementIfEscape(htmlDocument []rune, filteredHTMLDocument *[]rune, htmlDocumentLength int, index *int) bool {
	if htmlDocument[*index] == '\\' && *index+1 < htmlDocumentLength {
		*filteredHTMLDocument = append(*filteredHTMLDocument, htmlDocument[*index])
		*index++
		*filteredHTMLDocument = append(*filteredHTMLDocument, htmlDocument[*index])
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
				if hasPrefix(htmlDocumentRunes, delimiter.endDelimiter, htmlDocumentRunesLength, delimiter.endDelimiterLength, *index) {
					*index += delimiter.endDelimiterLength - 1
					return true
				}
			}
		}
	}

	return false
}

func appendAndIncrement(filteredHTMLDocument *[]rune, r rune, index *int) {
	*filteredHTMLDocument = append(*filteredHTMLDocument, r)
	*index++
}

func filterComments(htmlDocument string) string {
	var filteredHTMLDocument []rune
	var jsStringRune rune

	backtickRune := rune('`')
	jsBacktickInterpolationEnd := rune('}')

	htmlDocumentRunes := []rune(htmlDocument)
	htmlCommentStart := []rune("<!--")
	htmlCommentEnd := []rune("-->")
	jsCommentSingleLineStart := []rune("//")
	jsCommentSingleLineEnd := []rune("\n")
	commentMultiLineStart := []rune("/*")
	commentMultiLineEnd := []rune("*/")
	jsBacktickInterpolationStart := []rune("${")

	htmlDocumentRunesLength := len(htmlDocumentRunes)
	jsBacktickInterpolationStartLength := len(jsBacktickInterpolationStart)

	jsCommentSingleLineDelimiter := createCommentDelimiter(jsCommentSingleLineStart, jsCommentSingleLineEnd)
	commentMultiLineDelimiter := createCommentDelimiter(commentMultiLineStart, commentMultiLineEnd)

	commentDelimiters := []commentDelimiter{
		createCommentDelimiter(htmlCommentStart, htmlCommentEnd),
		jsCommentSingleLineDelimiter,
		commentMultiLineDelimiter,
	}
	commentMultiLineDelimiters := []commentDelimiter{commentMultiLineDelimiter}

	for i := 0; i < htmlDocumentRunesLength; i++ {
		// string escape or comment
		if appendAndIncrementIfEscape(htmlDocumentRunes, &filteredHTMLDocument, htmlDocumentRunesLength, &i) ||
			updateIndexIfComment(htmlDocumentRunes, htmlDocumentRunesLength, &i, commentDelimiters) {
			continue
		}

		// JavaScript string
		if htmlDocumentRunes[i] == '"' || htmlDocumentRunes[i] == '\'' {
			jsStringRune = htmlDocumentRunes[i]
			appendAndIncrement(&filteredHTMLDocument, htmlDocumentRunes[i], &i)
			for ; i < htmlDocumentRunesLength; i++ {
				if appendAndIncrementIfEscape(htmlDocumentRunes, &filteredHTMLDocument, htmlDocumentRunesLength, &i) {
					continue
				}

				filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
				if htmlDocumentRunes[i] == jsStringRune {
					break
				}
			}
			// JavaScript backtick string
		} else if htmlDocumentRunes[i] == backtickRune {
			appendAndIncrement(&filteredHTMLDocument, htmlDocumentRunes[i], &i)
			for ; i < htmlDocumentRunesLength; i++ {
				if appendAndIncrementIfEscape(htmlDocumentRunes, &filteredHTMLDocument, htmlDocumentRunesLength, &i) {
					continue
				}

				if hasPrefix(htmlDocumentRunes, jsBacktickInterpolationStart, htmlDocumentRunesLength, jsBacktickInterpolationStartLength, i) {
					jsBacktickInterpolationIsClosed := false
					if i+jsBacktickInterpolationStartLength < htmlDocumentRunesLength {
						for j := 0; j < jsBacktickInterpolationStartLength; j++ {
							filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i+j])
						}
						i += jsBacktickInterpolationStartLength
					}
					for ; i < htmlDocumentRunesLength; i++ {
						if updateIndexIfComment(htmlDocumentRunes, htmlDocumentRunesLength, &i, commentMultiLineDelimiters) {
							continue
						}

						if updateIndexIfHasPrefix(htmlDocumentRunes, jsCommentSingleLineStart, htmlDocumentRunesLength, jsCommentSingleLineDelimiter.startDelimiterLength, &i) {
							for ; i < htmlDocumentRunesLength; i++ {
								if htmlDocumentRunes[i] == jsBacktickInterpolationEnd {
									filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
									jsBacktickInterpolationIsClosed = true
									break
								}
							}
						} else {
							filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
						}
						if htmlDocumentRunes[i] == jsBacktickInterpolationEnd || jsBacktickInterpolationIsClosed {
							break
						}
					}
				} else {
					filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
					if htmlDocumentRunes[i] == backtickRune {
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
			actual := filterComments(tt.input)
			if actual != tt.expected {
				t.Errorf("filterComments() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
