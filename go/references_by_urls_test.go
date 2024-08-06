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

// TODO: should also return number of runes used to end the tag
func findHTMLTagAttributes(htmlDocumentPart string) ([]string, int) {
	var attributes []string
	var attributePart []rune
	var quoteRune rune
	numberOfRunesUsed := 0
	inAttributeName := false
	inAttributeValue := false

	runes := []rune(htmlDocumentPart)
	count := len(runes)

	for i := 0; i < count; i++ {
		iPlusOne := i + 1
		switch {
		case inAttributeName:
			if runes[i] == '=' {
				attributePart = append(attributePart, runes[i])
				inAttributeName = false
				inAttributeValue = true
			} else if isLetterDigitHyphenOrUnderscore(runes[i]) {
				attributePart = append(attributePart, runes[i])
			} else {
				inAttributeName = false
				attributePart = nil
			}
		case inAttributeValue:
			if runes[i] == '"' || runes[i] == '\'' {
				quoteRune = runes[i]
				attributePart = append(attributePart, runes[i])
				for j := iPlusOne; j < count; j++ {
					attributePart = append(attributePart, runes[j])
					if runes[j] == quoteRune {
						i = j
						break
					}
				}
				attributes = append(attributes, string(attributePart))
				inAttributeValue = false
				attributePart = nil
			}
		default:
			if isLetter(runes[i]) {
				attributePart = append(attributePart, runes[i])
				inAttributeName = true
			} else if iPlusOne < count && runes[i] == '/' && runes[iPlusOne] == '>' {
				numberOfRunesUsed = iPlusOne
				break
			} else if runes[i] == '>' {
				numberOfRunesUsed = i
				break
			} else if unicode.IsSpace(runes[i]) {
				continue
			}
		}
	}

	return attributes, numberOfRunesUsed
}

// TODO: WIP
func findTitleAndH1Elements(htmlDocument string) ([]string, []string) {
	var titleElements []string
	var h1Elements []string
	var htmlElementPart []rune
	creatingTitleElement := false
	// fidingTitleStartTag := false
	// fidingH1StartTag := false

	runes := []rune(htmlDocument)
	count := len(runes)

	for i := 0; i < count; i++ {
		iPlusOne := i + 1
		iPlusTwo := i + 2
		iPlusThree := i + 3
		iPlusFour := i + 4
		iPlusFive := i + 5
		iPlusSix := i + 6
		iPlusSeven := i + 7

		if creatingTitleElement {
			if iPlusSeven < count &&
				runes[i] == '<' && runes[iPlusOne] == '/' &&
				runes[iPlusTwo] == 't' && runes[iPlusThree] == 'i' && runes[iPlusFour] == 't' && runes[iPlusFive] == 'l' && runes[iPlusSix] == 'e' &&
				runes[iPlusSeven] == '>' {
				creatingTitleElement = false
				htmlElementPart = append(htmlElementPart, runes[i:iPlusSeven+1]...)
				titleElements = append(titleElements, string(htmlElementPart))
				htmlElementPart = nil
				i = iPlusSeven
			} else {
				htmlElementPart = append(htmlElementPart, runes[i])
			}
		} else {
			if iPlusSix < count &&
				runes[i] == '<' &&
				runes[iPlusOne] == 't' && runes[iPlusTwo] == 'i' && runes[iPlusThree] == 't' && runes[iPlusFour] == 'l' && runes[iPlusFive] == 'e' &&
				runes[iPlusSix] == '>' {
				creatingTitleElement = true
				htmlElementPart = append(htmlElementPart, runes[i:iPlusSix+1]...)
				i = iPlusSix
			} else if iPlusTwo < count &&
				runes[i] == '<' && runes[iPlusOne] == 'h' && runes[iPlusTwo] == '1' {
				i = iPlusTwo
			}
		}
	}

	return titleElements, h1Elements
}

// TODO: uses a runes sub slice, which is not optimal
func runesHaveStringPrefix(runes []rune, prefix string) (bool, int) {
	prefixRunes := []rune(prefix)
	length := len(prefixRunes)

	if len(runes) < length {
		return false, 0
	}

	for i := 0; i < length; i++ {
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
	count := len(runes)

	for i := 0; i < count; i++ {
		if inHTMLComment {
			if hasPrefix, prefixLength := runesHaveStringPrefix(runes[i:], "-->"); hasPrefix {
				inHTMLComment = false
				i += prefixLength
			}
		} else if inJSCommentSingleLine {
			if runes[i] == '\n' {
				inJSCommentSingleLine = false
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			}
		} else if inCommentMultiLine {
			if hasPrefix, prefixLength := runesHaveStringPrefix(runes[i:], "*/"); hasPrefix {
				inCommentMultiLine = false
				i += prefixLength
			}
		} else if escaped {
			filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			escaped = false
		} else {
			if runes[i] == '\\' {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
				escaped = true
			} else if hasPrefix, prefixLength := runesHaveStringPrefix(runes[i:], "<!--"); hasPrefix {
				inHTMLComment = true
				i += prefixLength
			} else if hasPrefix, prefixLength := runesHaveStringPrefix(runes[i:], "//"); hasPrefix {
				inJSCommentSingleLine = true
				i += prefixLength
			} else if hasPrefix, prefixLength := runesHaveStringPrefix(runes[i:], "/*"); hasPrefix {
				inCommentMultiLine = true
				i += prefixLength
			} else {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			}
		}
	}

	return string(filteredHTMLDocument)
}

// func TestReferencesByURLs(t *testing.T) {

// }
