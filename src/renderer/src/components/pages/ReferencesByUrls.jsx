import { createSignal } from 'solid-js'
import ResultPage from '../ResultPage'
import TextArea from '../TextArea'

export default function ReferencesByUrls(props) {
  const [textAreaValue, setTextAreaValue] = createSignal('')
  const [isValid, setIsValid] = createSignal('')
  const [references, setReferences] = createSignal('')

  function setStateInputComponent(textAreaValue) {
    setTextAreaValue(textAreaValue)

    if (textAreaValue === '') {
      setIsValid(false)
    } else {
      setIsValid(true)
    }
  }

  async function submit() {
    const result = await window.references.getReferencesByUrls(textAreaValue())
    setReferences(result)
  }

  const placeholderContent = (
    <>
      Add one or more website URls and press 'submit.' We can add URLs with or without spaces and
      multiple lines since they will get filtered out.
    </>
  )

  const inputComponent = (
    <div>
      <TextArea
        textAreaValue={textAreaValue}
        onChange={setStateInputComponent}
        placeholderContent={placeholderContent}
      />
      <button onClick={submit} disabled={!isValid()}>
        submit
      </button>
    </div>
  )

  const outputComponent = <textarea readonly value={references()} class="textarea-height-5" />

  return (
    <ResultPage
      title={props.title}
      inputComponent={inputComponent}
      outputComponent={outputComponent}
      getOutput={function () {}}
      onLoading={props.onLoading}
    />
  )
}
