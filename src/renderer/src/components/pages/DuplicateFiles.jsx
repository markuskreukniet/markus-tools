import FileOrFolderInput from '../FileOrFolderInput'
import Page from '../Page'

export default function DuplicateFiles() {
  async function getDuplicateFiles(filePaths) {
    const duplicateFiles = await window.duplicateFiles.getDuplicateFiles(filePaths)
    console.log('duplicateFiles', duplicateFiles)
  }

  return (
    <Page title="Duplicate Files Finder">
      <FileOrFolderInput onChange={getDuplicateFiles} />
      <h2>result:</h2>
      <textarea
        readonly
        placeholder="Add at least two files or a folder with two files and press 'submit.'"
      />
    </Page>
  )
}
