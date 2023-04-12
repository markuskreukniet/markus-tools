export function removeHtmlCssJavaScriptComments(text) {
  return text.replace(/\/\*[\s\S]*?\*\/|([^\\:]|^)\/\/.*$|<!--(.|\s)*?-->/gm, '')
}
