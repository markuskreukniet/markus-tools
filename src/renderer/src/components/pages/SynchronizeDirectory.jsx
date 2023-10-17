import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import FileOrFolderInput from '../filePathInput/FileOrFolderInput'
import { filePathSelectionType } from '../../../../preload/modules/files'

export default function SynchronizeDirectory(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  async function test() {
    const testA = await window.synchronization.synchronizeDirectoryBE('a', 'b')
    setStatus(testA)
  }

  function handleInputDirectoryRO(resultObject) {
    setGetOutput(test)
  }

  const inputComponent = (
    <FileOrFolderInput
      onChange={handleInputDirectoryRO}
      filePathSelectionType={filePathSelectionType.directory}
      maxOneInput
    />
  )

  return (
    <TextResultPage
      title={props.title}
      inputComponent={inputComponent}
      output={status()}
      getOutput={getOutput}
      onLoading={props.onLoading}
    />
  )
}
