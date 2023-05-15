const http = require('http')
const https = require('https')

export default async function referencesByUrls(urlsString) {
  const urls = getUrls(urlsString)
  for (const url of urls) {
    const httpData = await getData(url)
    const tags = findHtmlTags(httpData, 'h1')
    if (tags?.length === 1) {
      // https://css-tricks.com/snippets/javascript/strip-html-tags-in-javascript/
      const innerHtml = tags[0].replace(/(<([^>]+)>)/gi, '')
    }
  }

  console.log('urls', urls)

  return 'testing'
}

function getUrls(urlsString) {
  urlsString = urlsString.replaceAll('\n', '')
  urlsString = urlsString.replaceAll(' ', '')

  let urls = []
  const subStrings = ['http://', 'https://']

  while (urlsString.length > subStrings[0].length) {
    const httpIndex = urlsString.indexOf(subStrings[0], subStrings[0].length)
    const httpsIndex = urlsString.indexOf(subStrings[1], subStrings[1].length)

    let firstIndex = httpsIndex
    if (httpIndex === -1 && httpsIndex === -1) {
      urls = includesOneOfTheSubstringsAddToUrls(urlsString, subStrings, urls)
      return urls
    } else if (httpIndex === -1) {
      firstIndex = httpsIndex
    } else if (httpsIndex === -1) {
      firstIndex = httpIndex
    } else if (httpIndex < httpsIndex) {
      firstIndex = httpIndex
    }

    let beforeIndex = urlsString.slice(0, firstIndex)
    urls = includesOneOfTheSubstringsAddToUrls(beforeIndex, subStrings, urls)
    urlsString = urlsString.slice(firstIndex)
  }

  return urls
}

// TODO: includesOneOfTheSubstringsAddToUrls and includesOneOfTheSubstrings to one function
function includesOneOfTheSubstringsAddToUrls(string, subStrings, urls) {
  if (includesOneOfTheSubstrings(string, subStrings)) {
    urls.push(string)
  }

  return urls
}

function includesOneOfTheSubstrings(string, substrings) {
  for (const substring of substrings) {
    if (string.includes(substring)) {
      return true
    }
  }
  return false
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
