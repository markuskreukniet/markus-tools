import { createSignal } from 'solid-js'
import ResultByFilesPage from '../ResultByFilesPage'

export default function LinesOfCode(props) {
  const [linesOfCode, setLinesOfCode] = createSignal(0)

  async function setState(filePaths) {
    console.log(filePaths)
    setLinesOfCode(0)
  }

  const resultComponent = <p>Lines of code: {linesOfCode()}</p>

  return (
    <ResultByFilesPage
      resultComponent={resultComponent}
      handleFilePaths={setState}
      onLoading={props.onLoading}
    />
  )
}
