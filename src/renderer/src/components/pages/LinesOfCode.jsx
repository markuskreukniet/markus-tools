import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import FileOrFolderInput from '../FileOrFolderInput'

export default function LinesOfCode(props) {
  const [linesOfCode, setLinesOfCode] = createSignal(0)

  async function handleFilePaths(filePaths) {
    const linesOfCode = await window.codeQuality.getLinesOfCode(filePaths)
    setLinesOfCode(linesOfCode)
  }

  const inputComponent = <FileOrFolderInput onChange={handleFilePaths} minimumFiles={1} />

  return (
    <TextResultPage
      title={props.title}
      inputComponent={inputComponent}
      output={`Lines of code: ${linesOfCode()}`}
      minimumFiles={1}
      handleFilePaths={handleFilePaths}
      onLoading={props.onLoading}
    />
  )
}
