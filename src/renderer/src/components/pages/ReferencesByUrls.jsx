import { createSignal } from 'solid-js'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import ResultPage from '../page/ResultPage'
import { getEitherResultValueOrEitherResultToErrorString } from '../../../../preload/monads/either'
import TextArea from '../TextArea'

export default function ReferencesByUrls(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [textAreaValue, setTextAreaValue] = createSignal('')
  const [isValid, setIsValid] = createSignal('')
  const [eitherResultOutput, setEitherResultOutput] = createSignal('')

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
    setEitherResultOutput(
      getEitherResultValueOrEitherResultToErrorString(
        await window.goBackend.goFunctionCallBE('referencesByUrlsToJSON', textAreaValue())
      )
    )
  }

  const placeholderContent = (
    <>
      Add one or more website URls and press 'submit.' We can add URLs with or without spaces and
      multiple lines since they will get filtered out.
    </>
  )

  const inputComponent = (
    <div class="references-by-urls__input-component">
      <TextArea
        textAreaValue={textAreaValue}
        onChange={handleChange}
        placeholderContent={placeholderContent}
      />
      <ActivatableSubmitButton onAction={submit} active={isValid()} />
    </div>
  )

  const outputComponent = (
    <textarea readonly value={eitherResultOutput()} class="references-by-urls__output-textarea" />
  )

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
