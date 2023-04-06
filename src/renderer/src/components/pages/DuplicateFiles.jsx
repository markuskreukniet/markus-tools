import FileOrFolderInput from '../FileOrFolderInput'

export default function DuplicateFiles() {
  return (
    <div>
      <h1>Duplicate Files Finder</h1>

      <FileOrFolderInput onChange={(e) => console.log('e', e)} />

      <h2>result:</h2>
      <textarea rows="8" cols="55" />
    </div>
  )
}
