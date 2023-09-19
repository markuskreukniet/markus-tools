import { createEffect, createSignal } from 'solid-js'

// TODO: remove this component?
// step 1 is to give FileOrFolderInput to give the onChange also a hasFilePathObject boolean, and use that in SubmittableFileOrFolderInput
export default function ActiveByNumberButton(props) {
  const [active, setActive] = createSignal(false)

  createEffect(() => {
    if (props.currentNumber >= props.minimumNumber) {
      setActive(true)
    } else {
      setActive(false)
    }
  })

  return (
    <button onClick={() => props.onAction()} disabled={!active()}>
      {props.text}
    </button>
  )
}
