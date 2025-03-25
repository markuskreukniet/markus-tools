// regex: https://css-tricks.com/snippets/javascript/strip-html-tags-in-javascript/
// innerHtml = data
// .substring(startIndex, endIndex)
// .replace(/(<([^>]+)>)/gi, '')
// .trim()

// function fetchDataFromUrl(url) {
//   const protocol = url.startsWith('https') ? https : http
//   return new Promise((resolve, reject) => {
//     protocol
//       .get(url, (resp) => {
//         let data = ''
//         resp.on('data', (chunk) => {
//           data += chunk
//         })
//         resp.on('end', () => {
//           resolve(data)
//         })
//       })
//       .on('error', (error) => {
//         reject(error)
//       })
//   })
// }
