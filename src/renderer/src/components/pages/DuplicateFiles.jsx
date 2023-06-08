import { createSignal } from 'solid-js'
import ResultPage from '../page/ResultPage'
import FileOrFolderInput from '../FileOrFolderInput'
import TextArea from '../TextArea'

export default function DuplicateFiles(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [duplicateFiles, setDuplicateFiles] = createSignal('')

  async function setStateOutputComponent(filePaths) {
    const duplicateFiles = await window.duplicateFiles.getDuplicateFiles(filePaths)
    const textareaValue = duplicateFiles !== '' ? duplicateFiles : 'No duplicate files found'
    setDuplicateFiles(textareaValue)
  }

  function handleFilePaths(filePaths) {
    setGetOutput(setStateOutputComponent(filePaths))
  }

  const inputComponent = <FileOrFolderInput onChange={handleFilePaths} minimumFiles={2} />

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

  const outputComponent = (
    <TextArea readOnly textAreaValue={duplicateFiles} placeholderContent={placeholderContent} />
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
