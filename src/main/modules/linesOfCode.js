import { Either } from '../../preload/monads/either'
import { getUtf8FileContents } from './filePaths.js'

const endOfLine = '\n'

export default async function linesOfCode(filePaths) {
  // TODO: use error handling in GUI
  // TODO: use promise.all ?
  // const promises = filePaths.map((path) =>
  //   numberOfFileLinesWithoutCommentsAndEmptyLines(path.value)
  // )
  // const results = await Promise.all(promises)
  // return results.reduce((accumulator, currentValue) => accumulator + currentValue, 0)

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

function removeCommentsAndEmptyLines(code) {
  // should remove all JavaScript, HTML, and CSS comments
  code = code.replace(/\/\*[\s\S]*?\*\/|([^\\:]|^)\/\/.*$|<!--(.|\s)*?-->/gm, '')
  let lines = code.split(endOfLine)
  lines = lines.filter((line) => line.trim() !== '')
  return lines.join(endOfLine)
}
