package main

// import "testing"

// func toUncommentedHTMLElements(htmlDocument string) {
// 	var htmlElements []string
// 	var htmlElement []rune
// 	var wantedNextRune rune
// 	isCommentBuilding := false
// 	isCommentActive := false
// 	countToFindElementOrComment := 0

// 	htmlDocument = strings.TrimSpace(htmlDocument)

// 	// don't forget escaping \

// 	for _, r := range htmlDocument {
// 		if isCommentBuilding || isCommentActive {
// 			if r == wantedNextRune {
// 				if r == '!' {
// 					wantedNextRune = '-'
// 					countToFindElementOrComment++
// 				} else if r == '-' {
// 					if countToFindElementOrComment == 2 {
// 						wantedNextRune = '-'
// 					} else {
// 						isCommentActive = true
// 					}
// 				}
// 			} else {
// 				isCommentBuilding = false
// 				countToFindElementOrComment = 0
// 			}
// 		} else {
// 			if r == '<' {
// 				isCommentBuilding = true
// 				wantedNextRune = '!'
// 				countToFindElementOrComment++
// 			}
// 		}
// 	}

// }

// func TestReferencesByURLs(t *testing.T) {

// }
