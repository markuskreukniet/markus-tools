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
    httpData = await getData(url)
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
  urlsString = urlsString.replaceAll('\n', '')
  urlsString = urlsString.replaceAll(' ', '')

  let urls = []

  while (urlsString.length > protocolStrings[0].length) {
    const httpIndex = urlsString.indexOf(protocolStrings[0], protocolStrings[0].length)
    const httpsIndex = urlsString.indexOf(protocolStrings[1], protocolStrings[1].length)

    // case httpIndex === -1 is not needed
    let firstIndex = httpsIndex
    if (httpIndex === -1 && httpsIndex === -1) {
      urls = includesOneOfTheSubstringsAddToUrls(urlsString, protocolStrings, urls)
      return urls
    } else if (httpsIndex === -1 || httpIndex < httpsIndex) {
      firstIndex = httpIndex
    }

    const beforeIndex = urlsString.slice(0, firstIndex)
    urls = includesOneOfTheSubstringsAddToUrls(beforeIndex, protocolStrings, urls)
    urlsString = urlsString.slice(firstIndex)
  }

  return urls
}

function includesOneOfTheSubstringsAddToUrls(string, substrings, urls) {
  for (const substring of substrings) {
    if (string.includes(substring)) {
      urls.push(string)
    }
  }
  return urls
}

function getData(url) {
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
      .on('error', (err) => {
        console.log('Error:', err.message)
        reject()
      })
  })
}

function findHtmlTags(html, tag) {
  // TODO: check https://regex101.com/. It gives now an error. Check if regex is correct
  const regex = new RegExp(`<${tag}[^>]*>(.*?)<\\/${tag}>`, 'gi')
  return html.match(regex)
}
