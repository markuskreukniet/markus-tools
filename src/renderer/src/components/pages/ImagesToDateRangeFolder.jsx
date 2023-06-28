import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import FileOrFolderInput from '../FileOrFolderInput'
import { resultStatus } from '../../../../preload/modules/resultStatus'

export default function imagesToDateRangeFolder(props) {
  let resultPath = ''
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  async function setState(filePaths, path) {
    const done = await window.dateRangeFolder.imagesToDateRangeFolderBE(filePaths, path)
    const status = done ? 'done' : 'not done'
    setStatus(status)
  }

  function handleFilePaths(filePaths) {
    setGetOutput(setState(filePaths, resultPath))
  }

  async function test() {
    const done = await window.dialog.openFileDialogBE(true)

    if (done.status === resultStatus.ok) {
      console.log('done', done)
    } else {
      console.log('not done')
    }
  }

  // TODO: placeholder or label? TODO: minimumFiles should be 0 so it can only sort the files in destination path?
  const inputComponent = (
    <div>
      <FileOrFolderInput onChange={handleFilePaths} minimumFiles={2} />
      <input type="text" placeholder="" onChange={(e) => (resultPath = e.target.value)} />
      <button onClick={() => test()}>test</button>
    </div>
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
