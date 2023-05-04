import ResultPage from '../ResultPage'

export default function HtmlTitleWebScraper(props) {
  const inputComponent = (
    <div>
      <div />
    </div>
  )

  const outputComponent = (
    <div>
      <div />
    </div>
  )

  return (
    <ResultPage
      title={props.title}
      inputComponent={inputComponent}
      outputComponent={outputComponent}
      getOutput={function () {}}
      onLoading={props.onLoading}
    />
  )
}
