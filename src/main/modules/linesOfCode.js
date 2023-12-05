import { readFileSync } from 'node:fs'
import { Either, toEitherLeftResult, toEitherRightResult } from '../../preload/monads/either'

const endOfLine = '\n'

// TODO: functions this function starting from GUI are still async, which is not needed
export default function linesOfCode(filePaths) {
  let numberOfLines = 0
  for (const path of filePaths) {
    const result = numberOfFileLinesWithoutCommentsAndEmptyLines(path.value)
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
    return Either.left(error)
  }
  const code = removeCommentsAndEmptyLines(fileContents)
  const lines = code.split(endOfLine)
  return Either.right(lines.length)
}

function removeCommentsAndEmptyLines(code) {
  // should remove all JavaScript, HTML, and CSS comments
  code = code.replace(/\/\*[\s\S]*?\*\/|([^\\:]|^)\/\/.*$|<!--(.|\s)*?-->/gm, '')
  let lines = code.split(endOfLine)
  lines = lines.filter((line) => line.trim() !== '')
  return lines.join(endOfLine)
}
