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

  // TODO: this button be become a component that we at least also use in ActivatableButton. Or it is a button that is always active
  return (
    <button onClick={clickInput} class="secondary-button">{`add a ${
      props.directory ? 'directory' : 'file'
    }`}</button>
  )
}
