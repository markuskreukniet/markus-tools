import { dialog } from 'electron'

export default async function openFileDialog(selectFolder) {
  console.log('selectFolder', selectFolder)

  const properties = [selectFolder ? 'openDirectory' : 'openFile']
  const filters = selectFolder ? [{ name: 'All Files', extensions: ['*'] }] : []

  try {
    const result = await dialog.showOpenDialog({
      properties,
      filters
    })
    return result.canceled ? getFilePathResult([], 'ok') : getFilePathResult(result.filePaths, 'ok')
  } catch (error) {
    return getFilePathResult([], 'errorSystem')
  }
}

function getFilePathResult(filePaths, status) {
  return { result: filePaths, status }
}

// const resultStatus = Object.freeze({
//   ok: Symbol('ok'),
//   errorSystem: Symbol('errorSystem')
// })
