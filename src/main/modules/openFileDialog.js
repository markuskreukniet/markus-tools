import { dialog } from 'electron'
import {
  toResultObjectWithResultStatusErrorSystem,
  toResultObjectWithResultStatusOk
} from '../../preload/modules/resultStatus'

export default async function openFileDialog(selectDirectory) {
  const properties = [selectDirectory ? 'openDirectory' : 'openFile']
  const filters = selectDirectory ? [{ name: 'All Files', extensions: ['*'] }] : []

  try {
    const result = await dialog.showOpenDialog({
      properties,
      filters
    })
    return result.canceled
      ? toResultObjectWithResultStatusOk('')
      : toResultObjectWithResultStatusOk(result.filePaths[0])
  } catch (error) {
    return toResultObjectWithResultStatusErrorSystem('', error.message)
  }
}
