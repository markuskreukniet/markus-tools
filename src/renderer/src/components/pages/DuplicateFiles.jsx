import { createSignal } from 'solid-js'
import ResultPage from '../page/ResultPage'
import { eitherLeftResultToErrorString } from '../../../../preload/monads/either'
import SubmittableFileSystemNodesInput from '../filePathInput/SubmittableFileSystemNodesInput'
import TextArea from '../TextArea'

export default function DuplicateFiles(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [result, setResult] = createSignal('')

  async function setStateWithBE(uniqueFileSystemNodes) {
    const duplicateFiles = await window.goBackend.goFunctionCallBE(
      'getDuplicateFilesAsNewlineSeparatedStringToJSON',
      { uniqueFileSystemNodes }
    )
    const textareaValue =
      duplicateFiles.value !== '' ? duplicateFiles.value : 'No duplicate files found'
    setResult(textareaValue)
  }

  // TODO: looks a lot like filesToDateRangeDirectory handleInputFilePathsRO
  function handleChange(result) {
    if (result.isRight()) {
      setGetOutput(setStateWithBE(result.value))
    } else {
      setResult(eitherLeftResultToErrorString(result))
    }
  }

  const inputComponent = <SubmittableFileSystemNodesInput onChange={handleChange} />

  const placeholderContent = (
    <>
      <br />
      <br />
      The more files we select, the longer it takes to find duplicate files.
    </>
  )

  const outputComponent = (
    <TextArea readOnly textAreaValue={result} addToDefaultPlaceholderContent={placeholderContent} />
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
