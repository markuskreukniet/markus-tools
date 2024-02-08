import { createSignal } from 'solid-js'
import ResultPage from '../page/ResultPage'
import { eitherLeftResultToErrorString } from '../../../../preload/monads/either'
import SubmittableFileSystemNodesInput from '../filePathInput/SubmittableFileSystemNodesInput'
import TextArea from '../TextArea'

export default function DuplicateFiles(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [result, setResult] = createSignal('')

  async function setStateWithBE(filePathObjects) {
    const duplicateFiles = await window.duplicateFiles.duplicateFilesBE(filePathObjects)
    const textareaValue = duplicateFiles !== '' ? duplicateFiles : 'No duplicate files found'
    setResult(textareaValue)
  }

  // TODO: looks a lot like imagesToDateRangeFolder handleInputFilePathsRO
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
      Add at least two files or a directory with two files and press 'submit.'
      <br />
      <br />
      Adding a directory also adds the files of its subdirectories (its whole directory tree).
      <br />
      <br />
      The more files a directory has, or the more files we add, the longer it takes to find
      duplicate files.
    </>
  )

  const outputComponent = (
    <TextArea readOnly textAreaValue={result} placeholderContent={placeholderContent} />
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
