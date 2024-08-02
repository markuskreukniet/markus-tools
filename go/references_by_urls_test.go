package main

// import "testing"

// func toUncommentedHTMLElements(htmlDocument string) {
// 	var htmlElements []string
// 	var htmlElement []rune
// 	var wantedNextRune rune
// 	isCommentBuilding := false
// 	isCommentActive := false
// 	countToFindElementOrComment := 0

// 	var filteredHTMLDocument []rune

// 	htmlDocument = strings.TrimSpace(htmlDocument)

// 	// don't forget escaping \

// 	for _, r := range htmlDocument {
// 		if isCommentActive {
// 			if isCommentBuilding {
// 				if r == wantedNextRune {
// 					if r == '-' {
// 						wantedNextRune = '>'
// 					} else if r == '>' {
// 						isCommentActive = false
// 					}
// 				}
// 			} else {
// 				if r == '-' {
// 					wantedNextRune = '-'
// 				} else {
// 					filteredHTMLDocument = append(filteredHTMLDocument, r)
// 				}
// 			}
// 		} else {
// 			if isCommentBuilding {
// 				if r == wantedNextRune {
// 					if r == '!' {
// 						wantedNextRune = '-'
// 						countToFindElementOrComment++
// 					} else if r == '-' {
// 						if countToFindElementOrComment == 2 {
// 							wantedNextRune = '-'
// 							countToFindElementOrComment++
// 						} else {
// 							isCommentActive = true
// 						}
// 					}
// 				} else {
// 					isCommentBuilding = false
// 					countToFindElementOrComment = 0
// 				}
// 			} else {
// 				if r == '<' {
// 					isCommentBuilding = true
// 					wantedNextRune = '!'
// 					countToFindElementOrComment++
// 				} else {
// 					filteredHTMLDocument = append(filteredHTMLDocument, r)
// 				}
// 			}
// 		}
// 	}
// }

// func TestReferencesByURLs(t *testing.T) {

// }
