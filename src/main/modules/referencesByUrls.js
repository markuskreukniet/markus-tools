import http from 'http'
import https from 'https'
import { removeHtmlCssJavaScriptComments } from './modifyString.js'

export default async function referencesByUrls(urlsString) {
  const protocolStrings = ['http://', 'https://']

  // append urls
  const urls = []
  const urlsStringLines = urlsString.split('\n')
  for (const line of urlsStringLines) {
    let startIndex = 0
    while (startIndex < line.length) {
      startIndex = setUrlsAndGetStartIndex(protocolStrings[0], urls, line, startIndex)
      startIndex = setUrlsAndGetStartIndex(protocolStrings[1], urls, line, startIndex)
    }
  }

  // result
  let resultPart = urls.length > 0 ? await extractFormattedReference(urls[0], protocolStrings) : ''
  for (let i = 1; i < urls.length; i++) {
    resultPart += `, ${await extractFormattedReference(urls[i], protocolStrings)}`
  }
  return `(sources: ${resultPart}).`
}

async function extractFormattedReference(url, protocolStrings) {
  let data = ''
  try {
    data = await fetchDataFromUrl(url)
  } catch (error) {
    // TODO:
  }

  // extract first H1 inner HTML
  let innerHtml = ''
  data = removeHtmlCssJavaScriptComments(data)
  const startIndex = data.indexOf('<h1')
  if (startIndex !== -1) {
    const endTag = '</h1>'
    let endIndex = data.indexOf(endTag, startIndex)
    if (endIndex !== -1) {
      endIndex += endTag.length

      // regex: https://css-tricks.com/snippets/javascript/strip-html-tags-in-javascript/
      innerHtml = data
        .substring(startIndex, endIndex)
        .replace(/(<([^>]+)>)/gi, '')
        .trim()
    }
  }

  // result
  // 'let result = url' could be more efficient, but it is cleaner.
  let result = url
  if (innerHtml !== '') {
    for (const protocolString of protocolStrings) {
      if (url.startsWith(protocolString)) {
        result = `"${innerHtml}" by ${url.substr(
          protocolString.length,
          url.indexOf('/', protocolString.length) - protocolString.length
        )}`
      }
    }
  }
  return result
}

function setUrlsAndGetStartIndex(protocolString, urls, line, startIndex) {
  const index = line.indexOf(protocolString, startIndex)
  if (index !== -1) {
    let endIndex = line.indexOf(' ', index)
    if (endIndex === -1) {
      endIndex = line.length
    }
    urls.push(line.substring(index, endIndex))
    startIndex = endIndex
  }
  return startIndex
}

function fetchDataFromUrl(url) {
  const protocol = url.startsWith('https') ? https : http
  return new Promise((resolve, reject) => {
    protocol
      .get(url, (resp) => {
        let data = ''
        resp.on('data', (chunk) => {
          data += chunk
        })
        resp.on('end', () => {
          resolve(data)
        })
      })
      .on('error', (error) => {
        reject(error)
      })
  })
}
