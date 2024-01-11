import http from 'http'
import https from 'https'

export default async function referencesByUrls(urlsString) {
  const protocolStrings = ['http://', 'https://']
  const urls = getUrls(urlsString, protocolStrings)
  let result = urls.length > 0 ? await getReferencePart(urls[0], false, protocolStrings) : ''

  for (let i = 1; i < urls.length; i++) {
    result += await getReferencePart(urls[i], true, protocolStrings)
  }

  return `(sources: ${result}).`
}

async function getReferencePart(url, comma, protocolStrings) {
  let httpData = ''
  try {
    httpData = await fetchDataFromUrl(url)
  } catch (error) {
    //
  }

  const tags = httpData !== '' ? findHtmlTags(httpData, 'h1') : []
  if (tags?.length === 1) {
    // https://css-tricks.com/snippets/javascript/strip-html-tags-in-javascript/
    const innerHtml = tags[0].replace(/(<([^>]+)>)/gi, '')
    const part = `"${innerHtml}" by ${getByPart(url, protocolStrings)}`
    return comma ? `, ${part}` : part
  } else {
    return ''
  }
}

function getByPart(url, protocolStrings) {
  for (const protocolString of protocolStrings) {
    if (url.includes(protocolString)) {
      const endPosition = url.indexOf('/', protocolString.length) - protocolString.length
      return url.substr(protocolString.length, endPosition)
    }
  }

  return ''
}

function getUrls(urlsString, protocolStrings) {
  const urlsStringLines = urlsString.split('\n')
  const urls = []
  for (const line of urlsStringLines) {
    let startIndex = 0
    while (startIndex < line.length) {
      startIndex = testTest(protocolStrings[0], urls, line, startIndex)
      startIndex = testTest(protocolStrings[1], urls, line, startIndex)
    }
  }
  return urls
}

// TODO: better naming
function testTest(protocolString, urls, line, startIndex) {
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

function findHtmlTags(html, tag) {
  // TODO: check https://regex101.com/. It gives now an error. Check if regex is correct
  const regex = new RegExp(`<${tag}[^>]*>(.*?)<\\/${tag}>`, 'gi')
  return html.match(regex)
}
