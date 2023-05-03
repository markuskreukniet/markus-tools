import { createSignal, Show } from 'solid-js'
import ResultByFilesPage from '../ResultByFilesPage'

export default function DuplicateFiles(props) {
  const [duplicateFiles, setDuplicateFiles] = createSignal('')

  async function setState(filePaths) {
    const duplicateFiles = await window.duplicateFiles.getDuplicateFiles(filePaths)
    const textareaValue = duplicateFiles !== '' ? duplicateFiles : 'No duplicate files found'
    setDuplicateFiles(textareaValue)
  }

  const outputComponent = (
    <Show
      when={duplicateFiles() !== ''}
      fallback={
        <div class="custom-textarea-placeholder">
          Add at least two files or a folder with two files and press 'submit.'
          <br />
          <br />
          Adding a folder also adds the files of its subfolders (its whole folder tree).
          <br />
          <br />
          The more files a folder has, the more time it can take to add the files, which can be
          noticeable. Also, the more files we add, the longer it takes to find duplicate files.
        </div>
      }
    >
      <textarea readonly value={duplicateFiles()} placeholder="" />
    </Show>
  )

  return (
    <ResultByFilesPage
      outputComponent={outputComponent}
      minimumFiles={2}
      handleFilePaths={setState}
      onLoading={props.onLoading}
    />
  )
}
