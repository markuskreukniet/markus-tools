import { createEffect } from 'solid-js'
import Page from './Page'

export default function ResultPage(props) {
  createEffect(() => {
    loadGetOutput(props.getOutput).catch(() => {})
  })

  async function loadGetOutput(getOutput) {
    props.onLoading(true)
    await getOutput()
    props.onLoading(false)
  }

  return (
    <Page title={props.title}>
      {props.inputComponent}
      <h2>Result:</h2>
      {props.outputComponent}
    </Page>
  )
}
