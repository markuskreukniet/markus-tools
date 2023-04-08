import { createSignal } from 'solid-js'
import FileOrFolderInput from '../FileOrFolderInput'
import Page from '../Page'

export default function DuplicateFiles() {
  const [duplicateFiles, setDuplicateFiles] = createSignal('')

  async function handleFilePaths(filePaths) {
    const duplicateFiles = await window.duplicateFiles.getDuplicateFiles(filePaths)
    const textareaValue = duplicateFiles !== '' ? duplicateFiles : 'No duplicate files found'
    setDuplicateFiles(textareaValue)
  }

  return (
    <Page title="Duplicate Files Finder">
      <FileOrFolderInput onChange={handleFilePaths} />
      <h2>result:</h2>
      <textarea
        readonly
        value={duplicateFiles()}
        placeholder="Add at least two files or a folder with two files and press 'submit.' Adding a folder also adds the files of its subfolders (its whole folder tree). The more files a folder has, the more time it can take to add the files, which can be noticeable."
      />
    </Page>
  )
}
