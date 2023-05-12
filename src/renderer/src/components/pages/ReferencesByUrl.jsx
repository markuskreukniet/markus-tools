import { createSignal } from 'solid-js'
import ResultPage from '../ResultPage'
import TextArea from '../TextArea'

export default function ReferencesByUrl(props) {
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
    const result = await window.references.getReferencesByUrl(textAreaValue())
    setReferences(result)
  }

  const placeholderContent = (
    <>
      <div>
        placeholderContent
        <div />
      </div>
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

  const outputComponent = (
    <div>
      <div />
      {references()}
    </div>
  )

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
