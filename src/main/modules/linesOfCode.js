import { readFileSync } from 'node:fs'
import { removeHtmlCssJavaScriptComments } from './modifyString.js'
import { toEitherLeftResult, toEitherRightResult } from '../../preload/monads/either'

export default function linesOfCode(filePaths) {
  let numberOfLines = 0
  // TODO: are not filePaths?
  for (const path of filePaths) {
    let content = ''
    try {
      content = readFileSync(path.path, { encoding: 'utf8' })
    } catch (error) {
      return toEitherLeftResult(error.message)
    }
    content = removeHtmlCssJavaScriptComments(content)
    numberOfLines += content.split('\n').filter((line) => line.trim() !== '').length
  }
  return toEitherRightResult(numberOfLines)
}
