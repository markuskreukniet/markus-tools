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
  const innerHtml = extractFirstH1InnerHtml(data)

  // result
  let result = ''
  if (innerHtml !== '') {
    let domainName = ''
    for (const protocolString of protocolStrings) {
      if (url.startsWith(protocolString)) {
        domainName = url.substr(
          protocolString.length,
          url.indexOf('/', protocolString.length) - protocolString.length
        )
      }
    }
    result = `"${innerHtml}" by ${domainName}`
  }
  return result
}

// TODO: useless?
function extractFirstH1InnerHtml(html) {
  html = removeHtmlCssJavaScriptComments(html)
  const startIndex = html.indexOf('<h1')
  if (startIndex === -1) {
    return ''
  }
  const endTag = '</h1>'
  let endIndex = html.indexOf(endTag, startIndex)
  if (endIndex === -1) {
    return ''
  }
  endIndex += endTag.length

  // regex: https://css-tricks.com/snippets/javascript/strip-html-tags-in-javascript/
  return html
    .substring(startIndex, endIndex)
    .replace(/(<([^>]+)>)/gi, '')
    .trim()
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
