import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import { isEitherRightResult } from '../../../../preload/monads/either'
import SubmittableFileSystemNodesInput from '../filePathInput/SubmittableFileSystemNodesInput'

export default function LinesOfCode(props) {
  const [eitherResultOutput, setEitherResultOutput] = createSignal(null)
  const [getOutput, setGetOutput] = createSignal(function () {})

  async function setStateWithBE(fileSystemNodes) {
    const result = await window.codeQuality.linesOfCodeBE(fileSystemNodes)
    if (isEitherRightResult(result)) {
      result.value = `Lines of code: ${result.value}`
    }
    setEitherResultOutput(result)
  }

  function handleChange(result) {
    if (result.isRight()) {
      setGetOutput(setStateWithBE(result.value))
    } else {
      setEitherResultOutput(result)
    }
  }

  const inputComponent = <SubmittableFileSystemNodesInput onChange={handleChange} />

  return (
    <TextResultPage
      title={props.title}
      inputComponent={inputComponent}
      eitherResultOutput={eitherResultOutput()}
      getOutput={getOutput}
      onLoading={props.onLoading}
    />
  )
}
