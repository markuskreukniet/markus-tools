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

  async function setStateWithBE(fileSystemNodes) {
    const result = await window.codeQuality.linesOfCodeBE(fileSystemNodes)
    if (isEitherRightResult(result)) {
      setLinesOfCodeResult(`Lines of code: ${result.value}`)
    } else {
      setLinesOfCodeResultWithEitherLeftResultToErrorString(result)
    }
  }

  function handleChange(result) {
    if (result.isRight()) {
      setGetOutput(setStateWithBE(result.value))
    } else {
      setLinesOfCodeResultWithEitherLeftResultToErrorString(result)
    }
  }

  function setLinesOfCodeResultWithEitherLeftResultToErrorString(result) {
    setLinesOfCodeResult(eitherLeftResultToErrorString(result))
  }

  const inputComponent = <SubmittableFileOrFolderInput onChange={handleChange} />

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
