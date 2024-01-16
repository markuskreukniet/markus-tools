import { readFileSync } from 'node:fs'
import { Either, toEitherLeftResult, toEitherRightResult } from '../../preload/monads/either'
import { removeHtmlCssJavaScriptComments } from './modifyString.js'

// TODO:
const endOfLine = '\n'

export default function linesOfCode(filePaths) {
  let numberOfLines = 0
  for (const path of filePaths) {
    const result = numberOfFileLinesWithoutCommentsAndEmptyLines(path.path)
    if (result.isRight()) {
      numberOfLines += result.value
    } else {
      return toEitherLeftResult(result)
    }
  }
  return toEitherRightResult(numberOfLines)
}

function numberOfFileLinesWithoutCommentsAndEmptyLines(filePath) {
  let content = ''
  try {
    content = readFileSync(filePath, { encoding: 'utf8' })
  } catch (error) {
    return Either.left(error.message)
  }
  content = removeHtmlCssJavaScriptComments(content)
  return Either.right(content.split(endOfLine).filter((line) => line.trim() !== '').length)
}
