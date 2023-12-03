import { promises } from 'fs'
import { Either } from '../../preload/monads/either'

const endOfLine = '\n'

export default async function linesOfCode(filePaths) {
  // TODO: use error handling in GUI
  // When one numberOfFileLinesWithoutCommentsAndEmptyLines fails, the function should stop immediately, which is impossible with a promise.all solution.
  let numberOfLines = 0
  for (const path of filePaths) {
    const result = await numberOfFileLinesWithoutCommentsAndEmptyLines(path.value)
    if (result.isRight()) {
      numberOfLines += result.value
    } else {
      return result
    }
  }
  return numberOfLines
}

async function numberOfFileLinesWithoutCommentsAndEmptyLines(filePath) {
  const result = await getUtf8FileContents(filePath)
  if (result.isRight()) {
    const code = removeCommentsAndEmptyLines(result.value)
    const lines = code.split(endOfLine)
    return Either.right(lines.length)
  } else {
    return result
  }
}

async function getUtf8FileContents(filePath) {
  try {
    return Either.right(await promises.readFile(filePath, { encoding: 'utf8' }))
  } catch (error) {
    return Either.left(error)
  }
}

function removeCommentsAndEmptyLines(code) {
  // should remove all JavaScript, HTML, and CSS comments
  code = code.replace(/\/\*[\s\S]*?\*\/|([^\\:]|^)\/\/.*$|<!--(.|\s)*?-->/gm, '')
  let lines = code.split(endOfLine)
  lines = lines.filter((line) => line.trim() !== '')
  return lines.join(endOfLine)
}
