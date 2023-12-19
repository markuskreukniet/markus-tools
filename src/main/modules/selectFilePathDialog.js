import { dialog } from 'electron'
import { toEitherLeftResult, toEitherRightResult } from '../../preload/monads/either'

export default async function selectFilePathDialog(selectDirectory) {
  const properties = [selectDirectory ? 'openDirectory' : 'openFile']
  const filters = selectDirectory ? [{ name: 'All Files', extensions: ['*'] }] : []

  try {
    const result = await dialog.showOpenDialog({
      properties,
      filters
    })
    return result.canceled ? toEitherRightResult('') : toEitherRightResult(result.filePaths[0])
  } catch (error) {
    return toEitherLeftResult(error.message)
  }
}
