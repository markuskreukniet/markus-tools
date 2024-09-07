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

  function getVariantAttribute(variant) {
    if (variant === 'primary') {
      return { id: 'primary-button' }
    } else if (variant === 'secondary') {
      return { class: 'secondary-button' }
    } else if (variant === 'tertiary') {
      return { class: 'tertiary-button' }
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
