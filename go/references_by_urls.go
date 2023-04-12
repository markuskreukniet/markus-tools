package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func referencesByUrlsToJSON(rawUrls string) string {
	return resultToJSONFunctionResult(referencesByUrls(rawUrls))
}

func referencesByUrls(rawUrls string) string {
	hostNames, titles, errors, hasOnlyErrors := extractHostNamesAndTitlesOrdered(rawUrls)

	if hasOnlyErrors {
		return "No references found"
	}

	createSource := func(title, hostName string) string {
		return fmt.Sprintf("\"%s\" by %s", title, hostName)
	}

	var builder strings.Builder
	index := 0
	for i, title := range titles {
		if errors[i] == nil {
			builder.WriteString(fmt.Sprintf("(sources: %s", createSource(title, hostNames[i])))
			index = i + 1
			break
		}
	}

	for _, title := range titles[index:] {
		builder.WriteString(fmt.Sprintf(", %s", createSource(title, hostNames[index])))
		index++
	}

	builder.WriteString(")")

	return builder.String()
}

func splitTrimmedNonEmptyLines(s string) []string {
	var result []string
	scanner := bufio.NewScanner(strings.NewReader(s))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			result = append(result, line)
		}
	}

	return result
}

func extractHostNamesAndTitlesOrdered(rawUrls string) ([]string, []string, []error, bool) {
	findTitle := func(elements []string) (string, bool) {
		for _, element := range elements {
			title := strings.TrimSpace(removeTagsFromElement(element))
			if !utils.IsBlank(title) {
				return title, true
			}
		}

		return "", false
	}

	lines := splitTrimmedNonEmptyLines(rawUrls)
	length := len(lines)
	hostNames, titles, errors := make([]string, length), make([]string, length), make([]error, length) // TODO: use this multiple assign style also on other places
	errorCount := 0
	var group sync.WaitGroup

	for i, line := range lines {
		group.Add(1)

		go func(index int, rawUrl string) {
			defer group.Done()

			parsedURL, err := url.Parse(rawUrl)
			if err != nil {
				errorCount++
				return
			}

			page, err := downloadWebPage(rawUrl)
			if err != nil {
				errorCount++
				return
			}

			titleElements, h1Elements := findTitleAndH1Elements(filterComments(page))

			title, titleFound := findTitle(h1Elements)
			if !titleFound {
				title, _ = findTitle(titleElements)
			}

			hostNames[index], titles[index], errors[index] = parsedURL.Hostname(), title, err
		}(i, line)
	}

	group.Wait()

	return hostNames, titles, errors, length == errorCount
}

func downloadWebPage(url string) (string, error) {
	client := http.Client{Timeout: 8 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func isPeriod(r rune) bool {
	return r == '.'
}

func isForwardSlash(r rune) bool {
	return r == '/'
}

func isBackSlash(r rune) bool {
	return r == '\\'
}

func isGreaterThan(r rune) bool {
	return r == '>'
}

func isEqual(r rune) bool {
	return r == '='
}

func isSingleOrDoubleQuote(r rune) bool {
	return r == '\'' || r == '"'
}

func isLetter(r rune) bool {
	return unicode.IsUpper(r) || unicode.IsLower(r)
}

func isLetterOrUnderscore(r rune) bool {
	return isLetter(r) || r == '_'
}

func isLetterDigitHyphenOrUnderscore(r rune) bool {
	return isLetterOrUnderscore(r) || unicode.IsDigit(r) || r == '-'
}

// returns: tagPartLength, tagIsClosed, tagIsFound
func finishCreatingStartTag(document []rune, documentLength, index int) (int, bool, bool) {
	closingTagPart := []rune("/>")
	startTagEndPartLength, closingTagPartLength := 0, len(closingTagPart)

	var quoteRune rune
	inAttributeName, inAttributeValue := false, false

	for ; index < documentLength; index++ {
		if inAttributeName {
			if isEqual(document[index]) {
				startTagEndPartLength++
				inAttributeName = false
				inAttributeValue = true
			} else if isLetterDigitHyphenOrUnderscore(document[index]) {
				startTagEndPartLength++
			} else if unicode.IsSpace(document[index]) {
				startTagEndPartLength++
				inAttributeName = false
			} else {
				return 0, false, false
			}
		} else if inAttributeValue {
			if isSingleOrDoubleQuote(document[index]) {
				quoteRune = document[index]
				startTagEndPartLength++
				for i := index + 1; i < documentLength; i++ {
					startTagEndPartLength++
					if document[i] == quoteRune {
						index = i
						break
					}
				}
				inAttributeValue = false
			} else {
				return 0, false, false
			}
		} else {
			if isLetter(document[index]) {
				startTagEndPartLength++
				inAttributeName = true
			} else if isGreaterThan(document[index]) {
				startTagEndPartLength++
				return startTagEndPartLength, false, true
			} else if hasPrefix(document, closingTagPart, documentLength, closingTagPartLength, index) {
				startTagEndPartLength += closingTagPartLength
				return startTagEndPartLength, true, true
			} else if unicode.IsSpace(document[index]) {
				startTagEndPartLength++
			} else {
				return 0, false, false
			}
		}
	}

	return 0, false, false
}

// returns: elementPartLength, elementIsFound
func getTheOtherHTMLElementPartLength(document, startTagPart, endTagPart []rune, documentLength, startTagPartLength, endTagPartLength, index int) (int, bool) {
	tagPartLength := 0
	numberOfOpenStartTags := 1

	for i := index; i < documentLength; i++ {
		if updateTagPartLengthAndIndexIfHasPrefix(document, startTagPart, documentLength, startTagPartLength, &tagPartLength, &i) {
			if length, tagIsClosed, tagIsFound := finishCreatingStartTag(document, documentLength, i); tagIsFound {
				tagPartLength += length
				i += length - 1
				if !tagIsClosed {
					numberOfOpenStartTags++
				}
			} else {
				return 0, false
			}
		} else if updateTagPartLengthAndIndexIfHasPrefix(document, endTagPart, documentLength, endTagPartLength, &tagPartLength, &i) {
			for ; i < documentLength; i++ {
				tagPartLength++
				if isGreaterThan(document[i]) {
					numberOfOpenStartTags--
					if numberOfOpenStartTags == 0 {
						return tagPartLength, true
					}
				}
			}
		} else {
			tagPartLength++
		}
	}

	return 0, false
}

func hasOpenOrSelfClosingHTMLTagPrefix(document, prefix []rune, documentLength, prefixLength, index int) (int, bool, bool) {
	if hasPrefix(document, prefix, documentLength, prefixLength, index) {
		if tagPartLength, tagIsClosed, tagIsFound := finishCreatingStartTag(document, documentLength, index+prefixLength); tagIsFound {
			return prefixLength + tagPartLength, tagIsClosed, true
		}
	}

	return 0, false, false
}

// TODO: Finding HTML elements should happen for every element name in a complete HTML document since an element could be a child element of another element,
// which is possible to fix.
func findHTMLElements(document, elementName string) []string {
	var elements []string

	documentRunes := []rune(document)
	startTagPartRunes := append([]rune("<"), []rune(elementName)...)
	endTagPartRunes := append([]rune("</"), []rune(elementName)...)

	documentLength, startTagPartLength := len(documentRunes), len(startTagPartRunes)

	for i := 0; i < documentLength; i++ {
		length, tagIsClosed, hasPrefix := hasOpenOrSelfClosingHTMLTagPrefix(documentRunes, startTagPartRunes, documentLength, startTagPartLength, i)
		if hasPrefix {
			elementPart := documentRunes[i : i+length]
			i += length
			if !tagIsClosed {
				elementPartLength, elementIsFound := getTheOtherHTMLElementPartLength(documentRunes, startTagPartRunes, endTagPartRunes, documentLength, startTagPartLength, len(endTagPartRunes), i)
				if elementIsFound {
					elementPart = append(elementPart, documentRunes[i:i+elementPartLength]...)
					i += elementPartLength
				} else {
					return nil
				}
			}
			elements = append(elements, string(elementPart))
			i--
		}
	}

	return elements
}
func removeTagsFromElements(elements []string) []string {
	var elementsWithoutTags []string

	for _, element := range elements {
		elementWithoutTags := removeTagsFromElement(element)
		elementsWithoutTags = append(elementsWithoutTags, elementWithoutTags)
	}

	return elementsWithoutTags
}

func incrementIfPlusOneIsSmaller(i *int, j int) bool {
	if *i+1 < j {
		*i++
		return true
	}

	return false
}

func getTagLength(elementPart []rune, elementPartLength, indexArgument int) (int, bool) {
	index := indexArgument

	if elementPart[indexArgument] != '<' {
		return 0, false
	}

	if isLetterOrUnderscore(elementPart[indexArgument+1]) {
		index += 2
	} else if isForwardSlash(elementPart[indexArgument+1]) && isLetterOrUnderscore(elementPart[indexArgument+2]) {
		index += 3
	} else {
		return 0, false
	}

	var quoteRune rune
	tagFound := false

	for ; index < elementPartLength; index++ {
		if isLetterDigitHyphenOrUnderscore(elementPart[index]) {
			continue
		}

		if isGreaterThan(elementPart[index]) {
			tagFound = true
			break
		}

		indexPlusOne := index + 1
		if indexPlusOne < elementPartLength {
			if isForwardSlash(elementPart[index]) && isGreaterThan(elementPart[indexPlusOne]) {
				index++
				tagFound = true
				break
			}

			if isPeriod(elementPart[index]) && isLetter(elementPart[indexPlusOne]) {
				index++
				continue
			}
		}

		if unicode.IsSpace(elementPart[index]) {
			if !incrementIfPlusOneIsSmaller(&index, elementPartLength) {
				return 0, false
			}

			for ; index < elementPartLength; index++ {
				if !unicode.IsSpace(elementPart[index]) {
					break
				}
			}

			if !incrementIfPlusOneIsSmaller(&index, elementPartLength) {
				return 0, false
			}

			if isLetter(elementPart[index]) {
				index++
				for ; index < elementPartLength; index++ {
					if unicode.IsSpace(elementPart[index]) || isLetterDigitHyphenOrUnderscore(elementPart[index]) {
						continue
					}

					indexPlusOne = index + 1
					if indexPlusOne >= elementPartLength {
						return 0, false
					}

					if isPeriod(elementPart[index]) && isLetter(elementPart[indexPlusOne]) {
						index++
						continue
					}

					if isEqual(elementPart[index]) && isSingleOrDoubleQuote(elementPart[indexPlusOne]) {
						index++
						quoteRune = elementPart[index]
						index++
						for ; index < elementPartLength; index++ {
							if elementPart[index] == quoteRune {
								break
							} else if isBackSlash(elementPart[index]) {
								index++
							}
						}
						continue
					}

					if isGreaterThan(elementPart[index]) {
						tagFound = true
						break
					}

					if isForwardSlash(elementPart[index]) && isGreaterThan(elementPart[index+1]) {
						index++
						tagFound = true
						break
					}

					return 0, false
				}
			}

			if tagFound {
				break
			}

			continue
		}

		if tagFound {
			break
		}

		return 0, false
	}

	return index - indexArgument, true
}

func removeTagsFromElement(element string) string {
	var elementWithoutTags []rune
	elementRunes := []rune(element)
	elementRunesLength := len(elementRunes)

	for i := 0; i < elementRunesLength; i++ {
		if tagLength, tagFound := getTagLength(elementRunes, elementRunesLength, i); tagFound {
			i += tagLength
		} else {
			elementWithoutTags = append(elementWithoutTags, elementRunes[i])
		}
	}

	return string(elementWithoutTags)
}

func findTitleAndH1ElementsAndRemoveTags(htmlDocument string) ([]string, []string) {
	titles, h1s := findTitleAndH1Elements(htmlDocument)
	return removeTagsFromElements(titles), removeTagsFromElements(h1s)
}

func findTitleAndH1Elements(htmlDocument string) ([]string, []string) {
	return findHTMLElements(htmlDocument, "title"), findHTMLElements(htmlDocument, "h1")
}

func hasPrefix(runes, prefix []rune, runesLength, prefixLength, index int) bool {
	if runesLength-index < prefixLength {
		return false
	}

	for _, r := range prefix {
		if runes[index] != r {
			return false
		}
		index++
	}

	return true
}

func updateTagPartLengthAndIndexIfHasPrefix(runes, prefix []rune, runesLength, prefixLength int, tagPartLength, index *int) bool {
	if updateIndexIfHasPrefix(runes, prefix, runesLength, prefixLength, index) {
		*tagPartLength += prefixLength
		return true
	}

	return false
}

func updateIndexIfHasPrefix(runes, prefix []rune, runesLength, prefixLength int, index *int) bool {
	if hasPrefix(runes, prefix, runesLength, prefixLength, *index) {
		*index += prefixLength
		return true
	}

	return false
}

func appendAndIncrementIfEscape(htmlDocument []rune, filteredHTMLDocument *[]rune, htmlDocumentLength int, index *int) bool {
	if isBackSlash(htmlDocument[*index]) && *index+1 < htmlDocumentLength {
		*filteredHTMLDocument = append(*filteredHTMLDocument, htmlDocument[*index])
		*index++
		*filteredHTMLDocument = append(*filteredHTMLDocument, htmlDocument[*index])
		return true
	}

	return false
}

type commentDelimiter struct {
	startDelimiter       []rune
	endDelimiter         []rune
	startDelimiterLength int
	endDelimiterLength   int
}

func createCommentDelimiter(startDelimiter, endDelimiter []rune) commentDelimiter {
	return commentDelimiter{
		startDelimiter:       startDelimiter,
		endDelimiter:         endDelimiter,
		startDelimiterLength: len(startDelimiter),
		endDelimiterLength:   len(endDelimiter),
	}
}

func updateIndexIfComment(htmlDocumentRunes []rune, htmlDocumentRunesLength int, index *int, commentDelimiters []commentDelimiter) bool {
	for _, delimiter := range commentDelimiters {
		if updateIndexIfHasPrefix(htmlDocumentRunes, delimiter.startDelimiter, htmlDocumentRunesLength, delimiter.startDelimiterLength, index) {
			for ; *index < htmlDocumentRunesLength; *index++ {
				if hasPrefix(htmlDocumentRunes, delimiter.endDelimiter, htmlDocumentRunesLength, delimiter.endDelimiterLength, *index) {
					*index += delimiter.endDelimiterLength - 1
					return true
				}
			}
		}
	}

	return false
}

func appendAndIncrement(filteredHTMLDocument *[]rune, r rune, index *int) {
	*filteredHTMLDocument = append(*filteredHTMLDocument, r)
	*index++
}

func filterComments(htmlDocument string) string {
	var filteredHTMLDocument []rune
	var jsStringRune rune

	backtickRune := rune('`')
	jsBacktickInterpolationEnd := rune('}')

	htmlDocumentRunes := []rune(htmlDocument)
	htmlCommentStart, htmlCommentEnd := []rune("<!--"), []rune("-->")
	jsCommentSingleLineStart, jsCommentSingleLineEnd := []rune("//"), []rune("\n")
	commentMultiLineStart, commentMultiLineEnd := []rune("/*"), []rune("*/")
	jsBacktickInterpolationStart := []rune("${")

	htmlDocumentRunesLength := len(htmlDocumentRunes)
	jsBacktickInterpolationStartLength := len(jsBacktickInterpolationStart)

	jsCommentSingleLineDelimiter := createCommentDelimiter(jsCommentSingleLineStart, jsCommentSingleLineEnd)
	commentMultiLineDelimiter := createCommentDelimiter(commentMultiLineStart, commentMultiLineEnd)

	commentDelimiters := []commentDelimiter{
		createCommentDelimiter(htmlCommentStart, htmlCommentEnd),
		jsCommentSingleLineDelimiter,
		commentMultiLineDelimiter,
	}
	commentMultiLineDelimiters := []commentDelimiter{commentMultiLineDelimiter}

	for i := 0; i < htmlDocumentRunesLength; i++ {
		// string escape or comment
		if appendAndIncrementIfEscape(htmlDocumentRunes, &filteredHTMLDocument, htmlDocumentRunesLength, &i) ||
			updateIndexIfComment(htmlDocumentRunes, htmlDocumentRunesLength, &i, commentDelimiters) {
			continue
		}

		// JavaScript string
		if isSingleOrDoubleQuote(htmlDocumentRunes[i]) {
			jsStringRune = htmlDocumentRunes[i]
			appendAndIncrement(&filteredHTMLDocument, htmlDocumentRunes[i], &i)
			for ; i < htmlDocumentRunesLength; i++ {
				if appendAndIncrementIfEscape(htmlDocumentRunes, &filteredHTMLDocument, htmlDocumentRunesLength, &i) {
					continue
				}

				filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
				if htmlDocumentRunes[i] == jsStringRune {
					break
				}
			}
			// JavaScript backtick string
		} else if htmlDocumentRunes[i] == backtickRune {
			appendAndIncrement(&filteredHTMLDocument, htmlDocumentRunes[i], &i)
			for ; i < htmlDocumentRunesLength; i++ {
				if appendAndIncrementIfEscape(htmlDocumentRunes, &filteredHTMLDocument, htmlDocumentRunesLength, &i) {
					continue
				}

				if hasPrefix(htmlDocumentRunes, jsBacktickInterpolationStart, htmlDocumentRunesLength, jsBacktickInterpolationStartLength, i) {
					jsBacktickInterpolationIsClosed := false
					if i+jsBacktickInterpolationStartLength < htmlDocumentRunesLength {
						for j := 0; j < jsBacktickInterpolationStartLength; j++ {
							filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i+j])
						}
						i += jsBacktickInterpolationStartLength
					}
					for ; i < htmlDocumentRunesLength; i++ {
						if updateIndexIfComment(htmlDocumentRunes, htmlDocumentRunesLength, &i, commentMultiLineDelimiters) {
							continue
						}

						if updateIndexIfHasPrefix(htmlDocumentRunes, jsCommentSingleLineStart, htmlDocumentRunesLength, jsCommentSingleLineDelimiter.startDelimiterLength, &i) {
							for ; i < htmlDocumentRunesLength; i++ {
								if htmlDocumentRunes[i] == jsBacktickInterpolationEnd {
									filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
									jsBacktickInterpolationIsClosed = true
									break
								}
							}
						} else {
							filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
						}
						if htmlDocumentRunes[i] == jsBacktickInterpolationEnd || jsBacktickInterpolationIsClosed {
							break
						}
					}
				} else {
					filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
					if htmlDocumentRunes[i] == backtickRune {
						break
					}
				}
			}
		} else {
			filteredHTMLDocument = append(filteredHTMLDocument, htmlDocumentRunes[i])
		}
	}

	return string(filteredHTMLDocument)
}
