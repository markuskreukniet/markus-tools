// import { createSignal } from 'solid-js'
import FileOrFolderInput from '../FileOrFolderInput'
import Page from '../Page'

export default function LinesOfCode(props) {
  // const [linesOfCode, setLinesOfCode] = createSignal(0)

  async function handleFilePaths(filePaths) {
    props.onLoading(true) // duplicate
    console.log(filePaths)
    props.onLoading(false) // duplicate
  }

  return (
    <Page title={props.title}>
      <FileOrFolderInput onChange={handleFilePaths} /> {/* duplicate */}
      <h2>Result:</h2> {/* duplicate */}
      <p>Lines of code:</p>
    </Page>
  )
}
