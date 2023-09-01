import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import SubmittableFileOrFolderInput from '../filePathInput/SubmittableFileOrFolderInput'

export default function LinesOfCode(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [linesOfCode, setLinesOfCode] = createSignal(0)

  async function setOutput(filePaths) {
    const linesOfCode = await window.codeQuality.getLinesOfCode(filePaths)
    setLinesOfCode(linesOfCode)
  }

  function handleFilePaths(filePaths) {
    setGetOutput(setOutput(filePaths))
  }

  const inputComponent = (
    <SubmittableFileOrFolderInput onChange={handleFilePaths} minimumFiles={1} />
  )

  return (
    <TextResultPage
      title={props.title}
      inputComponent={inputComponent}
      output={`Lines of code: ${linesOfCode()}`}
      getOutput={getOutput}
      onLoading={props.onLoading}
    />
  )
}
