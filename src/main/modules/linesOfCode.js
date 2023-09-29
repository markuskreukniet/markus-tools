import fs from 'fs'

const endOfLine = '\n'

export default async function linesOfCode(filePaths) {
  let result = 0
  for (const path of filePaths) {
    result += numberOfFileLinesWithoutCommentsAndEmptyLines(path.value)
  }

  return result
}

function numberOfFileLinesWithoutCommentsAndEmptyLines(path) {
  let code = fs.readFileSync(path, { encoding: 'utf8' })
  code = removeCommentsAndEmptyLines(code)
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
