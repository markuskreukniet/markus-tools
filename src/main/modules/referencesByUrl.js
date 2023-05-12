export default async function referencesByUrl(urlsString) {
  console.log('urlsString', urlsString)
  console.log('typeof urlsString', typeof urlsString)

  // AI result
  const urls = urlsString.split(/\s+/).reduce((acc, val) => {
    const urlsString = val.split(/https?:\/\//)
    if (urlsString.length > 1) {
      urlsString.forEach((url) => {
        if (url.trim() !== '') {
          acc.push(`https://${url}`)
        }
      })
    } else {
      const url = urlsString[0].trim()
      if (url !== '') {
        acc.push(url)
      }
    }
    return acc
  }, [])

  console.log('urls', urls)
  return 'testing'
}
