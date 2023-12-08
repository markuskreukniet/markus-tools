import { isEitherRightResult } from '../../../../preload/monads/either'
import { toResultObjectWithResultStatusOk } from '../../../../preload/modules/resultStatus'

export default function FilePathSelector(props) {
  async function clickInput() {
    const result = await window.dialog.openFileDialogBE(props.directory)
    if (isEitherRightResult(result)) {
      const fileSystemNode = { path: result.value, isDirectory: false }
      if (props.directory) {
        fileSystemNode.isDirectory = true
      }
      // TODO: use either
      props.onChange(toResultObjectWithResultStatusOk(fileSystemNode))
    } else {
      // TODO: is wrong and should use either
      props.onChange(toResultObjectWithResultStatusOk(''))
    }
  }

  return <button onClick={clickInput}>{`add a ${props.directory ? 'directory' : 'file'}`}</button>
}
