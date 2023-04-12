import { For } from 'solid-js'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import { filePathSelectionType } from '../../../../preload/modules/files'
import FileSystemNodesInput from './FileSystemNodesInput'

// TODO: enum: all and maxOneDirectory

// TODO: is valid check should happen here

export default function SubmittableFileSystemNodeInputs(props) {
  function handleChange(result, handler) {
    if (result.isRight()) {
      if (result.value.hasFileSystemNode) {
        result.value = result.value.selectedFileSystemNodes[0].path
      } else {
        result.value = ''
      }
    }
    handler(result)
  }

  return (
    <div class="submittable-file-system-node-inputs">
      <For each={props.fileSystemNodesInputs}>
        {(input) => {
          switch (input.fileSystemNodesInputType) {
            case 'all':
              return <FileSystemNodesInput onChange={input.onChange} />
            case 'maxOneDirectory':
              return (
                <FileSystemNodesInput
                  onChange={(result) => handleChange(result, input.onChange)}
                  filePathSelectionType={filePathSelectionType.directory}
                  maxOneInput
                />
              )
            default:
              // TODO: error
              return null
          }
        }}
      </For>
      <div class="submittable-file-system-node-inputs__activatable-submit-button-wrapper">
        <ActivatableSubmitButton active={props.hasValidInput} onAction={props.onAction} />
      </div>
    </div>
  )

  // TODO: rename onAction to onChange???
}
