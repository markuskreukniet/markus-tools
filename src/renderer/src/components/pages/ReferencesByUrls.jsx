import { createSignal } from 'solid-js'
import ResultPage from '../page/ResultPage'
import TextArea from '../TextArea'

export default function ReferencesByUrls(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [textAreaValue, setTextAreaValue] = createSignal('')
  const [isValid, setIsValid] = createSignal('')
  const [references, setReferences] = createSignal('')

  function handleChange(textAreaValue) {
    setTextAreaValue(textAreaValue)

    if (textAreaValue === '') {
      setIsValid(false)
    } else {
      setIsValid(true)
    }
  }

  function submit() {
    setGetOutput(setStateWithBE())
  }

  async function setStateWithBE() {
    const result = await window.references.referencesByUrlsBE(textAreaValue())
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
        onChange={handleChange}
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
      getOutput={getOutput}
      onLoading={props.onLoading}
    />
  )
}
