import { getUtf8FileContents } from './filePaths.js'

const endOfLine = '\n'

export default async function linesOfCode(filePaths) {
  let result = 0
  for (const path of filePaths) {
    // TODO: use promise.all ?
    result += await numberOfFileLinesWithoutCommentsAndEmptyLines(path.value)
  }

  return result
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
