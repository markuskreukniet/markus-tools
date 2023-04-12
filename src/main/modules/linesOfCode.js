import { readFileSync } from 'node:fs'
import { removeHtmlCssJavaScriptComments } from './modifyString.js'
import { toEitherLeftResult, toEitherRightResult } from '../../preload/monads/either'

export default function linesOfCode(fileSystemNodes) {
  let numberOfLines = 0
  for (const node of fileSystemNodes) {
    let content = ''
    try {
      content = readFileSync(node.path, { encoding: 'utf8' })
    } catch (error) {
      return toEitherLeftResult(error.message)
    }
    content = removeHtmlCssJavaScriptComments(content)
    numberOfLines += content.split('\n').filter((line) => line.trim() !== '').length
  }
  return toEitherRightResult(numberOfLines)
}
