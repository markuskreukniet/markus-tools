import FileOrFolderInput from '../FileOrFolderInput'
import Page from '../Page'

export default function DuplicateFiles() {
  return (
    <Page title="Duplicate Files Finder">
      <FileOrFolderInput onChange={(e) => console.log('e', e)} />
      <h2>result:</h2>
      <textarea />
    </Page>
  )
}
