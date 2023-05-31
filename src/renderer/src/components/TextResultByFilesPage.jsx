import ResultByFilesPage from './ResultByFilesPage'

export default function TextResultByFilesPage(props) {
  const outputComponent = <p>{props.output}</p>

  return (
    <ResultByFilesPage
      title={props.title}
      outputComponent={outputComponent}
      minimumFiles={props.minimumFiles}
      handleFilePaths={props.handleFilePaths}
      onLoading={props.onLoading}
    />
  )
}
