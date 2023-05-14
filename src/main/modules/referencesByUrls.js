export default async function referencesByUrls(urlsString) {
  // AI result
  // const urls = urlsString.split(/\s+/).reduce((acc, val) => {
  //   const urlsString = val.split(/https?:\/\//)
  //   if (urlsString.length > 1) {
  //     urlsString.forEach((url) => {
  //       if (url.trim() !== '') {
  //         acc.push(`https://${url}`)
  //       }
  //     })
  //   } else {
  //     const url = urlsString[0].trim()
  //     if (url !== '') {
  //       acc.push(url)
  //     }
  //   }
  //   return acc
  // }, [])

  const urls = getUrls(urlsString)
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
