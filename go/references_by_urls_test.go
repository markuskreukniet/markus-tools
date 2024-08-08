package main

import "unicode"

func isLetter(r rune) bool {
	if unicode.IsUpper(r) || unicode.IsLower(r) {
		return true
	}

	return false
}

func isLetterDigitHyphenOrUnderscore(r rune) bool {
	if isLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
		return true
	}

	return false
}

func finishCreatingStartTag(htmlDocumentPart []rune) (int, bool) {
	var startTagEndPart []rune
	var quoteRune rune
	inAttributeName := false
	inAttributeValue := false

	length := len(htmlDocumentPart)

	for i := 0; i < length; i++ {
		switch {
		case inAttributeName:
			if htmlDocumentPart[i] == '=' {
				startTagEndPart = append(startTagEndPart, htmlDocumentPart[i])
				inAttributeName = false
				inAttributeValue = true
			} else if isLetterDigitHyphenOrUnderscore(htmlDocumentPart[i]) {
				startTagEndPart = append(startTagEndPart, htmlDocumentPart[i])
			} else if unicode.IsSpace(htmlDocumentPart[i]) {
				startTagEndPart = append(startTagEndPart, htmlDocumentPart[i])
				inAttributeName = false
			}
		case inAttributeValue:
			if htmlDocumentPart[i] == '"' || htmlDocumentPart[i] == '\'' {
				quoteRune = htmlDocumentPart[i]
				startTagEndPart = append(startTagEndPart, htmlDocumentPart[i])
				for j := i + 1; j < length; j++ {
					startTagEndPart = append(startTagEndPart, htmlDocumentPart[j])
					if htmlDocumentPart[j] == quoteRune {
						i = j
						break
					}
				}
				inAttributeValue = false
			}
		default:
			if isLetter(htmlDocumentPart[i]) {
				startTagEndPart = append(startTagEndPart, htmlDocumentPart[i])
				inAttributeName = true
			} else if htmlDocumentPart[i] == '>' {
				startTagEndPart = append(startTagEndPart, htmlDocumentPart[i])
				return len(startTagEndPart), false // TODO: should be an int instead of slice?
			} else if hasPrefix, length := hasStringPrefix(htmlDocumentPart[i:], "/>"); hasPrefix {
				startTagEndPart = append(startTagEndPart, htmlDocumentPart[i:length+2]...)
				return len(startTagEndPart) + length, true
			} else if unicode.IsSpace(htmlDocumentPart[i]) {
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
		} else if hasPrefix, length := hasStringPrefix(htmlDocumentPart[i:], endTag); hasPrefix {
			i += length - 1
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
			if hasPrefix, length := hasStringPrefix(runes[i:], "<title"); hasPrefix {
				creatingTitleStartTag = true
				htmlElementPart = append(htmlElementPart, runes[i:i+length]...)
				i += length - 1
				continue
			}
			if hasPrefix, length := hasStringPrefix(runes[i:], "<h1"); hasPrefix {
				creatingH1StartTag = true
				htmlElementPart = append(htmlElementPart, runes[i:i+length]...)
				i += length - 1
			}
		}
	}

	return titleElements, h1Elements
}

// TODO: uses a runes sub slice, which is not optimal
func hasStringPrefix(runes []rune, prefix string) (bool, int) {
	length := len(prefix)

	if len(runes) < length {
		return false, 0
	}

	for i, r := range prefix {
		if runes[i] != r {
			return false, 0
		}
	}

	return true, length
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
			if hasPrefix, length := hasStringPrefix(runes[i:], "-->"); hasPrefix {
				inHTMLComment = false
				i += length - 1
			}
		case inJSCommentSingleLine:
			if runes[i] == '\n' {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
				inJSCommentSingleLine = false
			}
		case inCommentMultiLine:
			if hasPrefix, length := hasStringPrefix(runes[i:], "*/"); hasPrefix {
				inCommentMultiLine = false
				i += length - 1
			}
		case escaped:
			filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			escaped = false
		default:
			if runes[i] == '\\' {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
				escaped = true
			} else if hasPrefix, length := hasStringPrefix(runes[i:], "<!--"); hasPrefix {
				inHTMLComment = true
				i += length - 1
			} else if hasPrefix, length := hasStringPrefix(runes[i:], "//"); hasPrefix {
				inJSCommentSingleLine = true
				length--
				i += length
			} else if hasPrefix, length := hasStringPrefix(runes[i:], "/*"); hasPrefix {
				inCommentMultiLine = true
				i += length - 1
			} else {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			}
		}
	}

	return string(filteredHTMLDocument)
}

// func TestReferencesByURLs(t *testing.T) {

// }
