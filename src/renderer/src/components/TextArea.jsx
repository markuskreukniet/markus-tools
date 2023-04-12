import { createEffect, createSignal, Show } from 'solid-js'

export default function TextArea(props) {
  const [showTextArea, setShowTextArea] = createSignal(false)

  createEffect(() => {
    setState(props.textAreaValue())
  })

  function setState(textAreaValue) {
    if (textAreaValue === '') {
      setShowTextArea(false)
    } else {
      setShowTextArea(true)
    }
  }

  function handleClick(e) {
    setShowTextArea(true)
    document.elementFromPoint(e.clientX, e.clientY).focus()
  }

  function handleBlur() {
    setState(props.textAreaValue())
  }

  function handleChange(e) {
    props.onChange(e.target.value)
  }

  function handleFunctionOrNull(readOnly, handleFunction) {
    return !readOnly ? handleFunction : null
  }

  function getPlaceholderContent() {
    const placeholderContent = (
      <>
        Select at least two files or a directory with two files and press 'submit.'
        <br />
        <br />
        Selecting a directory selects the files of the directory and its subdirectories (its whole
        directory tree).
      </>
    )

    if (
      (props.placeholderContent && props.addToDefaultPlaceholderContent) ||
      (!props.placeholderContent && !props.addToDefaultPlaceholderContent)
    ) {
      return placeholderContent
    } else if (props.placeholderContent) {
      return props.placeholderContent
    } else {
      return (
        <>
          {placeholderContent}
          {props.addToDefaultPlaceholderContent}
        </>
      )
    }
  }

  const placeholderClass = 'text-area__custom-placeholder'

  return (
    <Show
      when={showTextArea()}
      fallback={
        <Show
          when={props.readOnly}
          fallback={
            <button class={placeholderClass} onClick={handleClick}>
              {getPlaceholderContent()}
            </button>
          }
        >
          <p class={placeholderClass}>{getPlaceholderContent()}</p>
        </Show>
      }
    >
      <textarea
        readonly={props.readOnly}
        value={props.readOnly ? props.textAreaValue() : null}
        onInput={handleFunctionOrNull(props.readOnly, handleChange)}
        onBlur={handleFunctionOrNull(props.readOnly, handleBlur)}
        class="text-area__input"
      />
    </Show>
  )
}
