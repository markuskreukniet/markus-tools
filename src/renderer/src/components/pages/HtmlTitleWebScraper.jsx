import { createSignal } from 'solid-js'
import ResultPage from '../ResultPage'
import TextArea from '../TextArea'

export default function HtmlTitleWebScraper(props) {
  const [textAreaValue, setTextAreaValue] = createSignal('')

  const placeholderContent = (
    <>
      <div>
        content
        <div />
      </div>
    </>
  )

  const inputComponent = (
    <TextArea
      textAreaValue={textAreaValue}
      onChange={setTextAreaValue}
      placeholderContent={placeholderContent}
    />
  )

  const outputComponent = (
    <div>
      <div />
      {textAreaValue()}
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
