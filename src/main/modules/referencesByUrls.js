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

  // urlsString = urlsString.replaceAll('\n', '') // should happen when add to urls together with trim?

  // let urls = []
  // const httpsSplitted = splitWithSeparatorAsPrefixRecursion(urlsString, 'https://', [])
  // for (const element of httpsSplitted) {
  //   urls = [...urls, ...splitWithSeparatorAsPrefixRecursion(element, 'http://', [])]
  // }

  const urls = getUrls(urlsString)

  console.log('urls', urls)

  return 'testing'
}

// TODO: is separator.length correct? ja
function splitWithSeparatorAsPrefixRecursion(string, separator, array) {
  const separatorIndex = string.indexOf(separator, separator.length)
  if (separatorIndex === -1) {
    array.push(string)
    return array
  }

  const beforeSeparator = string.slice(0, separatorIndex)
  const separatorAndAfterIt = string.slice(separatorIndex)
  array.push(beforeSeparator)

  if (separatorAndAfterIt.includes(separator, separator.length)) {
    return splitWithSeparatorAsPrefixRecursion(separatorAndAfterIt, separator, array)
  } else {
    array.push(separatorAndAfterIt)
    return array
  }
}

function getUrls(urlsString) {
  let urls = []
  const subStrings = ['http://', 'https://']

  while (urlsString.length > subStrings[0].length) {
    const httpIndex = urlsString.indexOf(subStrings[0], subStrings[0].length)
    const httpsIndex = urlsString.indexOf(subStrings[1], subStrings[1].length)

    let firstIndex = httpsIndex
    if (httpIndex === -1 && httpsIndex === -1) {
      urls = possiblyUpdateUrls(urlsString, subStrings, urls)
      return urls
    } else if (httpIndex === -1) {
      firstIndex = httpsIndex
    } else if (httpsIndex === -1) {
      firstIndex = httpIndex
    } else if (httpIndex < httpsIndex) {
      firstIndex = httpIndex
    }

    let beforeIndex = urlsString.slice(0, firstIndex)
    urls = possiblyUpdateUrls(beforeIndex, subStrings, urls)

    urlsString = urlsString.slice(firstIndex)
  }

  return urls
}

function possiblyUpdateUrls(string, subStrings, urls) {
  if (includesOneOfTheSubstrings(string, subStrings)) {
    string = removeAllEndOfLineAndTrim(string)
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

function removeAllEndOfLineAndTrim(string) {
  string = string.replaceAll('\n', '')
  string = string.trim()
  return string
}
