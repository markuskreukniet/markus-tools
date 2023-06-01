import { createSignal } from 'solid-js'
import TextResultByFilesPage from '../page/TextResultByFilesPage'

export default function LinesOfCode(props) {
  const [linesOfCode, setLinesOfCode] = createSignal(0)

  async function setState(filePaths) {
    const linesOfCode = await window.codeQuality.getLinesOfCode(filePaths)
    setLinesOfCode(linesOfCode)
  }

  return (
    <TextResultByFilesPage
      title={props.title}
      output={`Lines of code: ${linesOfCode()}`}
      minimumFiles={1}
      handleFilePaths={setState}
      onLoading={props.onLoading}
    />
  )
}
