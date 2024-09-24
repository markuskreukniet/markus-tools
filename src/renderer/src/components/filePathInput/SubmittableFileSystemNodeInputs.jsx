import { For } from 'solid-js'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import FileSystemNodesInput from './FileSystemNodesInput'
import MaxOneDirectoryInput from './MaxOneDirectoryInput'

// TODO: enum: all and maxOneDirectory

export default function SubmittableFileSystemNodeInputs(props) {
  return (
    <div>
      <For each={props.fileSystemNodesInputs}>
        {(input) => {
          switch (input.fileSystemNodesInputType) {
            case 'all':
              return <FileSystemNodesInput onChange={input.onChange} />
            case 'maxOneDirectory':
              return <MaxOneDirectoryInput onChange={input.onChange} />
            default:
              // TODO: error
              return null
          }
        }}
      </For>
      <ActivatableSubmitButton active={props.hasValidInput} onAction={props.onAction} />
    </div>
  )
}
