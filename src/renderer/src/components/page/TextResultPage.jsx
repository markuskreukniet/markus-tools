import { createSignal, createEffect } from 'solid-js'
import ResultPage from './ResultPage'
import { isEitherRightResult } from '../../../../preload/monads/either'

export default function TextResultPage(props) {
  const [outputComponent, setOutputComponent] = createSignal(<></>)

  createEffect(() => {
    let output = ''
    if (props.eitherResultOutput) {
      if (isEitherRightResult(props.eitherResultOutput)) {
        output = props.eitherResultOutput.value || 'done'
      } else {
        output = `error: ${props.eitherResultOutput.value}`
      }
    }
    setOutputComponent(<p>{output}</p>)
  })

  return (
    <ResultPage
      title={props.title}
      inputComponent={props.inputComponent}
      outputComponent={outputComponent()}
      getOutput={props.getOutput}
      onLoading={props.onLoading}
    />
  )
}
