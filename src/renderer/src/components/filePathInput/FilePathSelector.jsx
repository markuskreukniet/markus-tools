import ActivatableButton from '../activatableButton/ActivatableButton'
import { Either, isEitherRightResult } from '../../../../preload/monads/either'

export default function FilePathSelector(props) {
  async function clickInput() {
    const result = await window.dialog.selectFilePathDialogBE(props.directory)
    if (isEitherRightResult(result)) {
      const fileSystemNode = { path: result.value, isDirectory: false }
      if (props.directory && result.value !== '') {
        fileSystemNode.isDirectory = true
      }
      props.onChange(Either.right(fileSystemNode))
    } else {
      props.onChange(Either.left(result.value))
    }
  }

  return (
    <ActivatableButton
      text={`add a ${props.directory ? 'directory' : 'file'}`}
      active
      onAction={clickInput}
      variant={'secondary'}
    />
  )
}
