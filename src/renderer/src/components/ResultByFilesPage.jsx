import FileOrFolderInput from '../FileOrFolderInput'
import Page from './Page'

export default function ResultByFilesPage(props) {
  return (
    <Page title={props.title}>
      <FileOrFolderInput onChange={props.handleFilePaths} />
      <h2>Result:</h2>
      {props.resultComponent}
    </Page>
  )
}
