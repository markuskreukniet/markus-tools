import { filePathType } from '../../../preload/modules/files'
import {
  isResultObjectOk,
  toResultObjectWithResultStatusOk
} from '../../../preload/modules/resultStatus'

export default function FilePathSelector(props) {
  async function clickInput() {
    const openFileDialogRO = await window.dialog.openFileDialogBE(props.directory)

    if (isResultObjectOk(openFileDialogRO)) {
      const result = { value: openFileDialogRO.result, filePathType: filePathType.file }
      if (props.directory) {
        result.filePathType = filePathType.directory
      }

      props.onChange(toResultObjectWithResultStatusOk(result))
    } else {
      props.onChange(openFileDialogRO)
    }
  }

  return <button onClick={clickInput}>{`add a ${props.directory ? 'directory' : 'file'}`}</button>
}
