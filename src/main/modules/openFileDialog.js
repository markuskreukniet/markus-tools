import { dialog } from 'electron'
import { resultStatus, getResultStatusCombination } from '../../preload/modules/resultStatus'

export default async function openFileDialog(selectFolder) {
  const properties = [selectFolder ? 'openDirectory' : 'openFile']
  const filters = selectFolder ? [{ name: 'All Files', extensions: ['*'] }] : []

  try {
    const result = await dialog.showOpenDialog({
      properties,
      filters
    })
    return result.canceled
      ? getResultStatusCombination([], resultStatus.ok)
      : getResultStatusCombination(result.filePaths, resultStatus.ok)
  } catch (error) {
    return getResultStatusCombination([], resultStatus.errorSystem)
  }
}
