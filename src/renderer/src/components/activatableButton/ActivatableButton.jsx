import { createEffect, createSignal } from 'solid-js'

// This component is not useless. We can use it as a self-closing tag, which reduces some code, and onClick does the same.
export default function ActivatableButton(props) {
  const [active, setActive] = createSignal(false)

  createEffect(() => {
    if (props.active) {
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
