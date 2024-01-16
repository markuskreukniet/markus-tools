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
  let fileContents = ''
  try {
    fileContents = readFileSync(filePath, { encoding: 'utf8' })
  } catch (error) {
    return Either.left(error.message)
  }
  const code = removeCommentsAndEmptyLines(fileContents)
  const lines = code.split(endOfLine)
  return Either.right(lines.length)
}

function removeCommentsAndEmptyLines(code) {
  code = removeHtmlCssJavaScriptComments(code)
  return code
    .split(endOfLine)
    .filter((line) => line.trim() !== '')
    .join(endOfLine)
}
