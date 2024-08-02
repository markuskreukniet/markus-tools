package main

// import "testing"

func filterComments(htmlDocument string) string {
	var filteredHTMLDocument []rune
	inHTMLComment := false
	inJSCommentSingleLine := false
	inJSCommentMultiLine := false

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
		} else if inJSCommentMultiLine {
			if iPlusOne < count && runes[i] == '*' && runes[iPlusOne] == '/' {
				inJSCommentMultiLine = false
				i = iPlusOne
			}
		} else {
			if iPlusThree < count && runes[i] == '<' && runes[iPlusOne] == '!' && runes[iPlusTwo] == '-' && runes[iPlusThree] == '-' {
				inHTMLComment = true
				i = iPlusThree
			} else if iPlusOne < count && runes[i] == '/' && runes[iPlusOne] == '/' {
				inJSCommentSingleLine = true
				i = iPlusOne
			} else if iPlusOne < count && runes[i] == '/' && runes[iPlusOne] == '*' {
				inJSCommentMultiLine = true
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
