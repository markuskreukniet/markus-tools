import {
  isResultObjectOk,
  toResultObjectWithResultStatusOk
} from '../../../../preload/modules/resultStatus'

export default function FilePathSelector(props) {
  async function clickInput() {
    const openFileDialogRO = await window.dialog.openFileDialogBE(props.directory)

    if (isResultObjectOk(openFileDialogRO) && openFileDialogRO.result !== '') {
      const result = { path: openFileDialogRO.result, isDirectory: false }
      if (props.directory) {
        result.isDirectory = true
      }

      props.onChange(toResultObjectWithResultStatusOk(result))
    } else {
      props.onChange(openFileDialogRO)
    }
  }

  return <button onClick={clickInput}>{`add a ${props.directory ? 'directory' : 'file'}`}</button>
}
