import { createEffect, createSignal } from 'solid-js'

// This component is not useless. We can use it as a self-closing tag, which reduces some code, and onClick does the same.\
// TODO: Activatable is not English?
export default function ActivatableButton(props) {
  const [active, setActive] = createSignal(false)

  createEffect(() => {
    if (props.active) {
      setActive(true)
    } else {
      setActive(false)
    }
  })

  function getVariantAttribute(variant) {
    if (variant === 'primary') {
      return { id: 'button--primary' }
    } else if (variant === 'secondary') {
      return { class: 'button--secondary' }
    } else if (variant === 'tertiary') {
      return { class: 'button--tertiary' }
    } else {
      return null
    }
  }

  return (
    <button
      onClick={() => props.onAction()}
      disabled={!active()}
      {...getVariantAttribute(props.variant)}
    >
      {props.text}
    </button>
  )
}
