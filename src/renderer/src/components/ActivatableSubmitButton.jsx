import { createEffect, createSignal } from 'solid-js'

export default function ActivatableSubmitButton(props) {
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
      submit
    </button>
  )
}
