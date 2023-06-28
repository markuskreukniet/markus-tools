import { dialog } from 'electron'
import { resultStatus, toResultObject } from '../../preload/modules/resultStatus'

export default async function openFileDialog(selectFolder) {
  const properties = [selectFolder ? 'openDirectory' : 'openFile']
  const filters = selectFolder ? [{ name: 'All Files', extensions: ['*'] }] : []

  try {
    const result = await dialog.showOpenDialog({
      properties,
      filters
    })
    return result.canceled
      ? toResultObject([], resultStatus.ok)
      : toResultObject(result.filePaths, resultStatus.ok)
  } catch (error) {
    return toResultObject([], resultStatus.errorSystem, error.message)
  }
}
