import { getUtf8FileContents } from './filePaths.js'

const endOfLine = '\n'

export default async function linesOfCode(filePaths) {
  // TODO: error handling
  const promises = filePaths.map((path) =>
    numberOfFileLinesWithoutCommentsAndEmptyLines(path.value)
  )
  const results = await Promise.all(promises)
  return results.reduce((accumulator, currentValue) => accumulator + currentValue, 0)
}

async function numberOfFileLinesWithoutCommentsAndEmptyLines(filePath) {
  // TODO: error handling
  let code = await getUtf8FileContents(filePath)
  code = removeCommentsAndEmptyLines(code.result)
  const lines = code.split(endOfLine)

  return lines.length
}

function removeCommentsAndEmptyLines(code) {
  // should remove all JavaScript, HTML, and CSS comments
  code = code.replace(/\/\*[\s\S]*?\*\/|([^\\:]|^)\/\/.*$|<!--(.|\s)*?-->/gm, '')

  let lines = code.split(endOfLine)
  lines = lines.filter((line) => line.trim() !== '')

  return lines.join(endOfLine)
}
