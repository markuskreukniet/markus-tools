import ResultPage from './ResultPage'

export default function TextResultPage(props) {
  const outputComponent = <p>{props.output}</p>

  return (
    <ResultPage
      title={props.title}
      inputComponent={props.inputComponent}
      outputComponent={outputComponent}
      getOutput={props.getOutput}
      onLoading={props.onLoading}
    />
  )
}
