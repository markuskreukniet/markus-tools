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
func finishCreatingStartTag(htmlDocumentPart []rune, index int) (int, bool, bool) {
	startTagEndPartLength := 0
	var quoteRune rune
	inAttributeName := false
	inAttributeValue := false

	length := len(htmlDocumentPart)

	for ; index < length; index++ {
		switch {
		case inAttributeName:
			if htmlDocumentPart[index] == '=' {
				startTagEndPartLength++
				inAttributeName = false
				inAttributeValue = true
			} else if isLetterDigitHyphenOrUnderscore(htmlDocumentPart[index]) {
				startTagEndPartLength++
			} else if unicode.IsSpace(htmlDocumentPart[index]) {
				startTagEndPartLength++
				inAttributeName = false
			} else {
				return 0, false, false
			}
		case inAttributeValue:
			if htmlDocumentPart[index] == '"' || htmlDocumentPart[index] == '\'' {
				quoteRune = htmlDocumentPart[index]
				startTagEndPartLength++
				for i := index + 1; i < length; i++ {
					startTagEndPartLength++
					if htmlDocumentPart[i] == quoteRune {
						index = i
						break
					}
				}
				inAttributeValue = false
			} else {
				return 0, false, false
			}
		default:
			if isLetter(htmlDocumentPart[index]) {
				startTagEndPartLength++
				inAttributeName = true
			} else if htmlDocumentPart[index] == '>' {
				startTagEndPartLength++
				return startTagEndPartLength, false, true
			} else if hasPrefix, length := hasStringPrefix(htmlDocumentPart, index, "/>"); hasPrefix {
				startTagEndPartLength += length
				return startTagEndPartLength, true, true
			} else if unicode.IsSpace(htmlDocumentPart[index]) {
				startTagEndPartLength++
			} else {
				return 0, false, false
			}
		}
	}

	return 0, false, false
}

// TODO: should work with index instead of sub slice
// TODO: should receive startTagPart and endTagPart
func finishCreatingHTMLElement(htmlDocumentPart []rune, startTagPart, endTagPart string) int {
	numberOfOpenStartTags := 1
	htmlDocumentPartLength := len(htmlDocumentPart)

	for i := 0; i < htmlDocumentPartLength; i++ {
		if hasPrefix, length := hasStringPrefix(htmlDocumentPart, i, startTagPart); hasPrefix {
			i += length
			length, tagIsClosed, _ := finishCreatingStartTag(htmlDocumentPart, i) // TODO _
			i += length - 1
			if !tagIsClosed {
				numberOfOpenStartTags++
			}
		} else if updateIndexIfPrefixMatches(htmlDocumentPart, endTagPart, &i) {
			for ; i < htmlDocumentPartLength; i++ {
				if htmlDocumentPart[i] == '>' {
					numberOfOpenStartTags--
					if numberOfOpenStartTags == 0 {
						return i + 1
					}
				}
			}
		} else if hasPrefix, length := hasStringPrefix(htmlDocumentPart, i, "/>"); hasPrefix {
			i += length - 1
			numberOfOpenStartTags--
			if numberOfOpenStartTags == 0 {
				return i + 1
			}
		}
	}

	return 0
}

// returns: tagLength, tagIsClosed, hasPrefix
func hasOpenOrSelfClosingHTMLTagPrefix(runes []rune, index int, prefix string) (int, bool, bool) {
	if hasPrefix, prefixLength := hasStringPrefix(runes, index, prefix); hasPrefix {
		if tagPartLength, tagIsClosed, tagIsFound := finishCreatingStartTag(runes, index+prefixLength); tagIsFound {
			return prefixLength + tagPartLength, tagIsClosed, hasPrefix
		}
		return 0, false, false
	}

	return 0, false, false
}

func findTitleAndH1Elements(htmlDocument string) ([]string, []string) {
	var titleElements []string
	var h1Elements []string
	var htmlElementPart []rune

	titleStartTagPart := "<title"
	titleEndTagPart := "</title"
	h1StartTagPart := "<h1"
	h1EndTagPart := "</h1"

	runes := []rune(htmlDocument)

	for i := 0; i < len(runes); i++ {
		if tagLength, tagIsClosed, hasPrefix := hasOpenOrSelfClosingHTMLTagPrefix(runes, i, titleStartTagPart); hasPrefix {
			htmlElementPart = append(htmlElementPart, runes[i:i+tagLength]...)
			i += tagLength
			if tagIsClosed {
				titleElements = append(titleElements, string(htmlElementPart))
				htmlElementPart = nil
			} else {
				length := finishCreatingHTMLElement(runes[i:], titleStartTagPart, titleEndTagPart)
				htmlElementPart = append(htmlElementPart, runes[i:i+length]...)
				i += length
				titleElements = append(titleElements, string(htmlElementPart))
				htmlElementPart = nil
			}
			i--
		} else if tagLength, tagIsClosed, hasPrefix := hasOpenOrSelfClosingHTMLTagPrefix(runes, i, h1StartTagPart); hasPrefix {
			htmlElementPart = append(htmlElementPart, runes[i:i+tagLength]...)
			i += tagLength
			if tagIsClosed {
				h1Elements = append(h1Elements, string(htmlElementPart))
				htmlElementPart = nil
			} else {
				length := finishCreatingHTMLElement(runes[i:], h1StartTagPart, h1EndTagPart)
				htmlElementPart = append(htmlElementPart, runes[i:i+length]...)
				i += length
				h1Elements = append(h1Elements, string(htmlElementPart))
				htmlElementPart = nil
			}
			i--
		}
	}

	return titleElements, h1Elements
}

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

func updateIndexIfPrefixMatches(runes []rune, prefix string, index *int) bool {
	hasPrefix, length := hasStringPrefix(runes, *index, prefix)

	if hasPrefix {
		*index += length - 1
	}

	return hasPrefix
}

func filterComments(htmlDocument string) string {
	var filteredHTMLDocument []rune
	inHTMLComment := false
	inJSCommentSingleLine := false
	inCommentMultiLine := false
	escaped := false

	runes := []rune(htmlDocument)

	for i := 0; i < len(runes); i++ {
		switch {
		case inHTMLComment:
			if updateIndexIfPrefixMatches(runes, "-->", &i) {
				inHTMLComment = false
			}
		case inJSCommentSingleLine:
			if runes[i] == '\n' {
				inJSCommentSingleLine = false
			}
		case inCommentMultiLine:
			if updateIndexIfPrefixMatches(runes, "*/", &i) {
				inCommentMultiLine = false
			}
		case escaped:
			filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			escaped = false
		default:
			if runes[i] == '\\' {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
				escaped = true
			} else if updateIndexIfPrefixMatches(runes, "<!--", &i) {
				inHTMLComment = true
			} else if updateIndexIfPrefixMatches(runes, "//", &i) {
				inJSCommentSingleLine = true
			} else if updateIndexIfPrefixMatches(runes, "/*", &i) {
				inCommentMultiLine = true
			} else {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			}
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
			expectedTitles: []string{},
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
			expectedH1s:    []string{},
		},
		{
			htmlDocument:   "<html><head></head><body><h1 te-st=\"a---test\"  lol  ><h1 asdf-l=\"test\"  /></h1   ></body></html>",
			expectedTitles: []string{},
			expectedH1s:    []string{"<h1 te-st=\"a---test\"  lol  ><h1 asdf-l=\"test\"  /></h1   >"},
		},
		{
			htmlDocument:   "<html><head></head><body><h1><h1></h1></h1></body></html>",
			expectedTitles: []string{},
			expectedH1s:    []string{"<h1><h1></h1></h1>"},
		},
	}

	for _, test := range tests {
		titles, h1s := findTitleAndH1Elements(test.htmlDocument)

		if titles != nil && test.expectedTitles != nil {
			if !reflect.DeepEqual(titles, test.expectedTitles) {
				t.Errorf("findTitleAndH1Elements(%q) titles = %v; want %v", test.htmlDocument, titles, test.expectedTitles)
			}
		}

		if h1s != nil && test.expectedH1s != nil {
			if !reflect.DeepEqual(h1s, test.expectedH1s) {
				t.Errorf("findTitleAndH1Elements(%q) h1s = %v; want %v", test.htmlDocument, h1s, test.expectedH1s)
			}
		}
	}
}
