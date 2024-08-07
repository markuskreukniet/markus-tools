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

func finishCreatingStartTag(htmlDocumentPart []rune) int {
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
				return len(startTagEndPart) // TODO: should be an int instead of slice?
			} else if hasPrefix, length := runesHaveStringPrefixAndGetPrefixLengthMinusOne(htmlDocumentPart[i:], "/>"); hasPrefix {
				startTagEndPart = append(startTagEndPart, htmlDocumentPart[i:length+2]...)
				return len(startTagEndPart) + length
			} else if unicode.IsSpace(htmlDocumentPart[i]) {
				continue
			}
		}
	}

	return 0
}

func findTitleAndH1Elements(htmlDocument string) ([]string, []string) {
	var titleElements []string
	var h1Elements []string
	var htmlElementPart []rune
	creatingTitleStartTag := false
	// creatingH1StartTag := false

	runes := []rune(htmlDocument)

	for i := range runes {
		switch {
		case creatingTitleStartTag:
			// length := finishCreatingStartTag(runes[i:])
		default:
			if hasPrefix, length := runesHaveStringPrefixAndGetPrefixLengthMinusOne(runes[i:], "<title"); hasPrefix {
				creatingTitleStartTag = true
				htmlElementPart = append(htmlElementPart, runes[i:length+2]...)
				i += length
			} else if hasPrefix, length := runesHaveStringPrefixAndGetPrefixLengthMinusOne(runes[i:], "<h1"); hasPrefix {
				// creatingH1StartTag = true
				htmlElementPart = append(htmlElementPart, runes[i:length+2]...)
				i += length
			}
		}
	}

	return titleElements, h1Elements
}

// TODO: uses a runes sub slice, which is not optimal
func runesHaveStringPrefixAndGetPrefixLengthMinusOne(runes []rune, prefix string) (bool, int) {
	prefixRunes := []rune(prefix)
	length := len(prefixRunes)

	if len(runes) < length {
		return false, 0
	}

	for i := range prefixRunes {
		if runes[i] != prefixRunes[i] {
			return false, 0
		}
	}

	return true, length - 1
}

func filterComments(htmlDocument string) string {
	var filteredHTMLDocument []rune
	inHTMLComment := false
	inJSCommentSingleLine := false
	inCommentMultiLine := false
	escaped := false

	runes := []rune(htmlDocument)

	for i := range runes {
		switch {
		case inHTMLComment:
			if hasPrefix, length := runesHaveStringPrefixAndGetPrefixLengthMinusOne(runes[i:], "-->"); hasPrefix {
				inHTMLComment = false
				i += length
			}
		case inJSCommentSingleLine:
			if runes[i] == '\n' {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
				inJSCommentSingleLine = false
			}
		case inCommentMultiLine:
			if hasPrefix, length := runesHaveStringPrefixAndGetPrefixLengthMinusOne(runes[i:], "*/"); hasPrefix {
				inCommentMultiLine = false
				i += length
			}
		case escaped:
			filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			escaped = false
		default:
			if runes[i] == '\\' {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
				escaped = true
			} else if hasPrefix, length := runesHaveStringPrefixAndGetPrefixLengthMinusOne(runes[i:], "<!--"); hasPrefix {
				inHTMLComment = true
				i += length
			} else if hasPrefix, length := runesHaveStringPrefixAndGetPrefixLengthMinusOne(runes[i:], "//"); hasPrefix {
				inJSCommentSingleLine = true
				i += length
			} else if hasPrefix, length := runesHaveStringPrefixAndGetPrefixLengthMinusOne(runes[i:], "/*"); hasPrefix {
				inCommentMultiLine = true
				i += length
			} else {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			}
		}
	}

	return string(filteredHTMLDocument)
}

// func TestReferencesByURLs(t *testing.T) {

// }
