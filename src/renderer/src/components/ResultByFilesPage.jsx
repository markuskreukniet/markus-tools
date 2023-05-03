import FileOrFolderInput from './FileOrFolderInput'
import ResultPage from './ResultPage'

export default function ResultByFilesPage(props) {
  async function handleFilePaths(filePaths) {
    props.onLoading(true)
    await props.handleFilePaths(filePaths)
    props.onLoading(false)
  }

  const inputComponent = (
    <FileOrFolderInput onChange={handleFilePaths} minimumFiles={props.minimumFiles} />
  )

  return (
    <ResultPage
      title={props.title}
      inputComponent={inputComponent}
      outputComponent={props.outputComponent}
    />
  )
}
