export default async function referencesByUrls(urlsString) {
  console.log('urlsString', urlsString)
  console.log('typeof urlsString', typeof urlsString)

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

  urlsString = 'https://lol2.com'

  const httpsSplitted = splitWithSeparatorAsPrefixRecursion(urlsString, 'http://', [])

  //let urls = []
  //const httpsSplitted = splitWithSeparatorAsPrefixRecursion(urlsString, 'https://', [])
  console.log('httpsSplitted', httpsSplitted)
  // for (const element of httpsSplitted) {
  //   urls = [...urls, ...splitWithSeparatorAsPrefixRecursion(element, 'http://', [])]
  // }

  // console.log('urls', urls)

  return 'testing'
}

function splitWithSeparatorAsPrefixRecursion(string, separator, array) {
  const separatorIndex = string.indexOf(separator, separator.length)
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
