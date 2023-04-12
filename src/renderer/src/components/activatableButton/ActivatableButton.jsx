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
    switch (variant) {
      case 'primary':
        return { id: 'button--primary' }
      case 'secondary':
        return { class: 'button--secondary' }
      case 'tertiary':
        return { class: 'button--tertiary' }
      default:
        return null
    }
  }

  return (
    <button
      onMouseDown={() => props.onAction()}
      disabled={!active()}
      {...getVariantAttribute(props.variant)}
    >
      {props.text}
    </button>
  )
}
