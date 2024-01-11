import http from 'http'
import https from 'https'

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
  if (innerHtml !== '') {
    return `"${innerHtml}" by ${getByPart(url, protocolStrings)}`
  } else {
    return ''
  }
}

function extractFirstH1InnerHtml(html) {
  // TODO: remove comments first

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
    .trimStart()
    .trimEnd()
}

// TODO: useless?
function getByPart(url, protocolStrings) {
  for (const protocolString of protocolStrings) {
    if (url.includes(protocolString)) {
      const endPosition = url.indexOf('/', protocolString.length) - protocolString.length
      return url.substr(protocolString.length, endPosition)
    }
  }

  return ''
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
