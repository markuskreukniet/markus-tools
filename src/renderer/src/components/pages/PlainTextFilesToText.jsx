import { createSignal } from 'solid-js'
import ResultPage from '../page/ResultPage'
import {
  eitherLeftResultToErrorString,
  isEitherRightResult
} from '../../../../preload/monads/either'
import SubmittableFileSystemNodesInput from '../filePathInput/SubmittableFileSystemNodesInput'
import TextArea from '../TextArea'

// TODO: PlainTextFilesToText is almost a copy of DuplicateFiles
export default function PlainTextFilesToText(props) {
  const [eitherResultOutput, setEitherResultOutput] = createSignal('')
  const [getOutput, setGetOutput] = createSignal(function () {})

  async function setStateWithBE(uniqueFileSystemNodes) {
    const result = await window.goBackend.goFunctionCallBE('plainTextFilesToTextToJSON', {
      uniqueFileSystemNodes
    })
    if (isEitherRightResult(result)) {
      setEitherResultOutput(result.value !== '' ? result.value : 'No text found')
    } else {
      setEitherResultOutput(eitherLeftResultToErrorString(result))
    }
  }

  function handleChange(result) {
    if (result.isRight()) {
      setGetOutput(setStateWithBE(result.value))
    } else {
      setEitherResultOutput(eitherLeftResultToErrorString(result))
    }
  }

  const inputComponent = <SubmittableFileSystemNodesInput onChange={handleChange} />

  const placeholderContent = (
    <>
      <br />
      <br />
      The more files we select, the longer it takes to show the text.
    </>
  )

  const outputComponent = (
    <TextArea
      readOnly
      textAreaValue={eitherResultOutput}
      addToDefaultPlaceholderContent={placeholderContent}
    />
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
