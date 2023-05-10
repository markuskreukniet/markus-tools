import { createEffect, createSignal, Show } from 'solid-js'

// TODO: eslint-disable-next-line solid/reactivity
export default function TextArea(props) {
  const [showTextArea, setShowTextArea] = createSignal(false)

  createEffect(() => {
    setState(props.readOnly, props.textAreaValue())
  })

  function setState(readOnly, textAreaValue) {
    if (readOnly) {
      if (textAreaValue === '') {
        setShowTextArea(false)
      } else {
        setShowTextArea(true)
      }
    }
  }

  function handleClick() {
    //
  }

  function handleBlur() {
    //
  }

  function handleChange() {
    //
  }

  function handleFunctionOrNull(readOnly, handleFunction) {
    return !readOnly ? handleFunction : null
  }

  return (
    <Show
      when={showTextArea()}
      fallback={
        <div
          class="custom-textarea-placeholder"
          onClick={handleFunctionOrNull(props.readOnly, handleClick)}
        >
          {props.placeholderContent}
        </div>
      }
    >
      <textarea
        readonly
        value={props.readOnly ? props.textAreaValue() : null}
        // eslint-disable-next-line solid/reactivity
        onChange={handleFunctionOrNull(props.readOnly, handleChange)}
        // eslint-disable-next-line solid/reactivity
        onBlur={handleFunctionOrNull(props.readOnly, handleBlur)}
      />
    </Show>
  )
}
