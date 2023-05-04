import { createSignal } from 'solid-js'
import FileOrFolderInput from './FileOrFolderInput'
import ResultPage from './ResultPage'

export default function ResultByFilesPage(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})

  function handleFilePaths(filePaths) {
    setGetOutput(props.handleFilePaths(filePaths))
  }

  const inputComponent = (
    <FileOrFolderInput onChange={handleFilePaths} minimumFiles={props.minimumFiles} />
  )

  return (
    <ResultPage
      title={props.title}
      inputComponent={inputComponent}
      outputComponent={props.outputComponent}
      getOutput={getOutput}
      onLoading={props.onLoading}
    />
  )
}
