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

  urlsString = urlsString.replaceAll('\n', '') // should happen when add to urls together with trim?

  let urls = []
  const httpsSplitted = splitWithSeparatorAsPrefixRecursion(urlsString, 'https://', [])
  for (const element of httpsSplitted) {
    urls = [...urls, ...splitWithSeparatorAsPrefixRecursion(element, 'http://', [])]
  }

  console.log('urls', urls)

  return 'testing'
}

// TODO: is separator.length correct?
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

// does not work
// function splitWithSeparatorAsPrefix(stringWithSeparator, separator) {
//   const result = []

//   while (stringWithSeparator.includes(separator, separator.length)) {
//     const separatorIndex = stringWithSeparator.indexOf(separator, separator.length)
//     const beforeSeparator = stringWithSeparator.slice(0, separatorIndex)
//     stringWithSeparator = stringWithSeparator.slice(separatorIndex)
//     result.push(beforeSeparator)
//   }

//   if (result.length === 0) {
//     result.push(stringWithSeparator)
//   }

//   return result
// }
