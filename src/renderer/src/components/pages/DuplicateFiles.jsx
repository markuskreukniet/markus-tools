import { createSignal } from 'solid-js'
import ResultPage from '../page/ResultPage'
import { isResultObjectOk } from '../../../../preload/modules/resultStatus'
import SubmittableFileOrFolderInput from '../filePathInput/SubmittableFileOrFolderInput'
import TextArea from '../TextArea'

export default function DuplicateFiles(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [result, setResult] = createSignal('')

  async function setStateOutputComponent(filePathObjects) {
    const duplicateFiles = await window.duplicateFiles.getDuplicateFiles(filePathObjects)
    const textareaValue = duplicateFiles !== '' ? duplicateFiles : 'No duplicate files found'
    setResult(textareaValue)
  }

  // TODO: looks a lot like imagesToDateRangeFolder handleInputFilePathsRO
  function handleFilePathsRO(resultObject) {
    if (isResultObjectOk(resultObject)) {
      setGetOutput(setStateOutputComponent(resultObject.result))
    } else {
      setResult(resultObject.message)
    }
  }

  const inputComponent = <SubmittableFileOrFolderInput onChange={handleFilePathsRO} />

  const placeholderContent = (
    <>
      Add at least two files or a folder with two files and press 'submit.'
      <br />
      <br />
      Adding a folder also adds the files of its subfolders (its whole folder tree).
      <br />
      <br />
      The more files a folder has, the more time it can take to add the files, which can be
      noticeable. Also, the more files we add, the longer it takes to find duplicate files.
    </>
  )

  // TODO: should be result()?
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
