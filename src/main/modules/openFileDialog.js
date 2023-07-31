import { dialog } from 'electron'
import {
  toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem,
  toResultObjectWithEmptyArrayResultAndResultStatusOk,
  toResultObjectWithResultStatusOk
} from '../../preload/modules/resultStatus'

export default async function openFileDialog(selectFolder) {
  const properties = [selectFolder ? 'openDirectory' : 'openFile', 'multiSelections']
  const filters = selectFolder ? [{ name: 'All Files', extensions: ['*'] }] : []

  try {
    const result = await dialog.showOpenDialog({
      properties,
      filters
    })
    return result.canceled
      ? toResultObjectWithEmptyArrayResultAndResultStatusOk()
      : toResultObjectWithResultStatusOk(result.filePaths)
  } catch (error) {
    return toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(error.message)
  }
}
