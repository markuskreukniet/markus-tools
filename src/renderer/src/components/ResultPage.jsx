import Page from './Page'

export default function ResultPage(props) {
  return (
    <Page title={props.title}>
      {props.inputComponent}
      <h2>Result:</h2>
      {props.outputComponent}
    </Page>
  )
}
