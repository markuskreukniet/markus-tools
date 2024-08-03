package main

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

func filterComments(htmlDocument string) string {
	var filteredHTMLDocument []rune
	inHTMLComment := false
	inJSCommentSingleLine := false
	inCommentMultiLine := false
	escaped := false

	runes := []rune(htmlDocument)
	count := len(runes)

	for i := 0; i < count; i++ {
		iPlusOne := i + 1
		iPlusTwo := i + 2
		iPlusThree := i + 3

		if inHTMLComment {
			if iPlusTwo < count && runes[i] == '-' && runes[iPlusOne] == '-' && runes[iPlusTwo] == '>' {
				inHTMLComment = false
				i = iPlusTwo
			}
		} else if inJSCommentSingleLine {
			if runes[i] == '\n' {
				inJSCommentSingleLine = false
			}
		} else if inCommentMultiLine {
			if iPlusOne < count && runes[i] == '*' && runes[iPlusOne] == '/' {
				inCommentMultiLine = false
				i = iPlusOne
			}
		} else if escaped {
			filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			escaped = false
		} else {
			if runes[i] == '\\' {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
				escaped = true
			} else if iPlusThree < count && runes[i] == '<' && runes[iPlusOne] == '!' && runes[iPlusTwo] == '-' && runes[iPlusThree] == '-' {
				inHTMLComment = true
				i = iPlusThree
			} else if iPlusOne < count && runes[i] == '/' && runes[iPlusOne] == '/' {
				inJSCommentSingleLine = true
				i = iPlusOne
			} else if iPlusOne < count && runes[i] == '/' && runes[iPlusOne] == '*' {
				inCommentMultiLine = true
				i = iPlusOne
			} else {
				filteredHTMLDocument = append(filteredHTMLDocument, runes[i])
			}
		}
	}

	return string(filteredHTMLDocument)
}

// func TestReferencesByURLs(t *testing.T) {

// }
