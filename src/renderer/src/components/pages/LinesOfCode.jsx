import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import SubmittableFileOrFolderInput from '../filePathInput/SubmittableFileOrFolderInput'

export default function LinesOfCode(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [linesOfCode, setLinesOfCode] = createSignal(0)

  async function setOutput(filePaths) {
    setLinesOfCode(await window.codeQuality.getLinesOfCode(filePaths))
  }

  function handleFilePaths(filePaths) {
    setGetOutput(setOutput(filePaths.result.selectedFilePathObjects))
  }

  const inputComponent = <SubmittableFileOrFolderInput onChange={handleFilePaths} />

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
