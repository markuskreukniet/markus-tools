import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import {
  eitherLeftResultToErrorString,
  isEitherRightResult
} from '../../../../preload/monads/either'
import SubmittableFileOrFolderInput from '../filePathInput/SubmittableFileOrFolderInput'

export default function LinesOfCode(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [linesOfCodeResult, setLinesOfCodeResult] = createSignal('')

  async function setOutput(filePaths) {
    const result = await window.codeQuality.getLinesOfCode(filePaths)
    if (isEitherRightResult(result)) {
      setLinesOfCodeResult(`Lines of code: ${result.value}`)
    } else {
      setLinesOfCodeResult(eitherLeftResultToErrorString(result))
    }
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
