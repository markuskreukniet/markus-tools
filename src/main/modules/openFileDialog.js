import dialog from 'electron'

export default async function openFileDialog(selectFolder) {
  const properties = [selectFolder ? 'openDirectory' : 'openFile']
  const filters = selectFolder ? [{ name: 'All Files', extensions: ['*'] }] : []

  try {
    const result = await dialog.showOpenDialog({
      properties,
      filters
    })
    return result.canceled
      ? getFilePathResult([], resultStatus.ok)
      : getFilePathResult(result.filePaths, resultStatus.ok)
  } catch (error) {
    return getFilePathResult([], resultStatus.errorUserAction)
  }
}

function getFilePathResult(filePaths, status) {
  return { result: filePaths, status }
}

const resultStatus = Object.freeze({
  ok: Symbol('ok'),
  errorUserAction: Symbol('errorUserAction')
})