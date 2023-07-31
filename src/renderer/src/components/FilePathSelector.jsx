import { isResultObjectOk } from '../../../preload/modules/resultStatus'

export default function FilePathSelector(props) {
  async function clickInput() {
    const openFileDialogRO = await window.dialog.openFileDialogBE(props.directory)

    if (isResultObjectOk(openFileDialogRO)) {
      props.onChange(openFileDialogRO.result)
    } else {
      // TODO:
    }
  }

  return <button onClick={clickInput}>{`add a ${props.directory ? 'directory' : 'file'}`}</button>
}
