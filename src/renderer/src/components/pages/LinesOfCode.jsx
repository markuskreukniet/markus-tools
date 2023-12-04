import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import SubmittableFileOrFolderInput from '../filePathInput/SubmittableFileOrFolderInput'

export default function LinesOfCode(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [linesOfCodeResult, setLinesOfCodeResult] = createSignal(0)

  async function setOutput(filePaths) {
    // TODO: error handling
    setLinesOfCodeResult(`Lines of code: ${await window.codeQuality.getLinesOfCode(filePaths)}`)
  }

  function handleFilePaths(filePaths) {
    // TODO: error handling
    setGetOutput(setOutput(filePaths.result))
  }

  const inputComponent = <SubmittableFileOrFolderInput onChange={handleFilePaths} />

  return (
    <TextResultPage
      title={props.title}
      inputComponent={inputComponent}
      output={linesOfCodeResult()}
      getOutput={getOutput}
      onLoading={props.onLoading}
    />
  )
}
