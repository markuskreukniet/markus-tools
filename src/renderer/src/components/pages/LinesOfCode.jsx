import { createSignal } from 'solid-js'
import ResultByFilesPage from '../ResultByFilesPage'

export default function LinesOfCode(props) {
  const [linesOfCode, setLinesOfCode] = createSignal(0)

  async function setState(filePaths) {
    const linesOfCode = await window.codeQuality.getLinesOfCode(filePaths)
    setLinesOfCode(linesOfCode)
  }

  const resultComponent = <p>Lines of code: {linesOfCode()}</p>

  return (
    <ResultByFilesPage
      resultComponent={resultComponent}
      minimumFiles={1}
      handleFilePaths={setState}
      onLoading={props.onLoading}
    />
  )
}
