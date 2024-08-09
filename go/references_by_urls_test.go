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

func finishCreatingStartTag(htmlDocumentPart []rune) (int, bool) {
	startTagEndPartLength := 0
	var quoteRune rune
	inAttributeName := false
	inAttributeValue := false

	length := len(htmlDocumentPart)

	for i := 0; i < length; i++ {
		switch {
		case inAttributeName:
			if htmlDocumentPart[i] == '=' {
				startTagEndPartLength++
				inAttributeName = false
				inAttributeValue = true
			} else if isLetterDigitHyphenOrUnderscore(htmlDocumentPart[i]) {
				startTagEndPartLength++
			} else if unicode.IsSpace(htmlDocumentPart[i]) {
				startTagEndPartLength++
				inAttributeName = false
			}
		case inAttributeValue:
			if htmlDocumentPart[i] == '"' || htmlDocumentPart[i] == '\'' {
				quoteRune = htmlDocumentPart[i]
				startTagEndPartLength++
				for j := i + 1; j < length; j++ {
					startTagEndPartLength++
					if htmlDocumentPart[j] == quoteRune {
						i = j
						break
					}
				}
				inAttributeValue = false
			} else {
				// TODO: results in endless loop? return 0, false
			}
		default:
			if isLetter(htmlDocumentPart[i]) {
				startTagEndPartLength++
				inAttributeName = true
			} else if htmlDocumentPart[i] == '>' {
				startTagEndPartLength++
				return startTagEndPartLength, false
			} else if hasPrefix, length := hasStringPrefix(htmlDocumentPart, i, "/>"); hasPrefix {
				startTagEndPartLength += length
				return startTagEndPartLength, true
			} else if unicode.IsSpace(htmlDocumentPart[i]) {
				startTagEndPartLength++
				continue
			}
		}
	}

	return 0, false
}

func finishCreatingHTMLElement(htmlDocumentPart []rune, endTag string) int {
	inEndTag := false

	for i := 0; i < len(htmlDocumentPart); i++ {
		if inEndTag {
			if unicode.IsSpace(htmlDocumentPart[i]) {
				continue
			} else if htmlDocumentPart[i] == '>' {
				return i + 1
			} else {
				return 0
			}
		} else if updateIndexIfPrefixMatches(htmlDocumentPart, endTag, &i) {
			inEndTag = true
		}
	}

	return 0
}

func findTitleAndH1Elements(htmlDocument string) ([]string, []string) {
	var titleElements []string
	var h1Elements []string
	var htmlElementPart []rune
	creatingTitleStartTag := false
	creatingH1StartTag := false
	finishCreatingTitleElement := false
	finishCreatingH1Element := false

	runes := []rune(htmlDocument)

	for i := 0; i < len(runes); i++ {
		switch {
		case creatingTitleStartTag:
			length, tagIsClosed := finishCreatingStartTag(runes[i:])
			htmlElementPart = append(htmlElementPart, runes[i:i+length]...)
			i += length - 1
			if tagIsClosed {
				titleElements = append(titleElements, string(htmlElementPart))
				htmlElementPart = nil
			} else {
				finishCreatingTitleElement = true
			}
			creatingTitleStartTag = false
		case creatingH1StartTag:
			length, tagIsClosed := finishCreatingStartTag(runes[i:])
			htmlElementPart = append(htmlElementPart, runes[i:i+length]...)
			i += length - 1
			if tagIsClosed {
				h1Elements = append(h1Elements, string(htmlElementPart))
				htmlElementPart = nil
			} else {
				finishCreatingH1Element = true
			}
			creatingH1StartTag = false
		case finishCreatingTitleElement:
			length := finishCreatingHTMLElement(runes[i:], "</title")
			htmlElementPart = append(htmlElementPart, runes[i:i+length]...)
			i += length - 1
			titleElements = append(titleElements, string(htmlElementPart))
			htmlElementPart = nil
			finishCreatingTitleElement = false
		case finishCreatingH1Element:
			length := finishCreatingHTMLElement(runes[i:], "</h1")
			htmlElementPart = append(htmlElementPart, runes[i:i+length]...)
			i += length - 1
			h1Elements = append(h1Elements, string(htmlElementPart))
			htmlElementPart = nil
			finishCreatingH1Element = false
		default:
			if hasPrefix, length := hasStringPrefix(runes, i, "<title"); hasPrefix {
				creatingTitleStartTag = true
				htmlElementPart = append(htmlElementPart, runes[i:i+length]...)
				i += length - 1
				continue
			}
			if hasPrefix, length := hasStringPrefix(runes, i, "<h1"); hasPrefix {
				creatingH1StartTag = true
				htmlElementPart = append(htmlElementPart, runes[i:i+length]...)
				i += length - 1
			}
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
