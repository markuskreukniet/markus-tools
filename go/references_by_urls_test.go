package main

// TODO: does the escaping work with \n?
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
