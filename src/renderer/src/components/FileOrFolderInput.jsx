import { createSignal, For } from 'solid-js'
import FileSelector from './FileSelector'

// TODO:
// Adding a file could add a duplicate file since there could already be a folder with possible child folders already containing that file.
// Adding a folder could add a duplicate file since the folder with possible child folders could contain a duplicate file.
export default function FileOrFolderInput(props) {
  const [selectedPaths, setSelectedPaths] = createSignal([])
  let filePaths = []

  function handleSelectedFile(files) {
    if (!selectedPaths().some((path) => path === files[0].path)) {
      setSelectedPaths([...selectedPaths(), files[0].path])
      handleFilePaths(files)
    }
  }

  function handleSelectedFolder(files) {
    const folderPath = getSelectedFolderPath(files)

    if (!selectedPaths().some((path) => path === folderPath)) {
      setSelectedPaths([...selectedPaths(), folderPath])
      handleFilePaths(files)
    }
  }

  function handleFilePaths(files) {
    // files is a FileList, not an array, so we can't use .map
    for (const file of files) {
      filePaths.push(file.path)
    }
  }

  function reset() {
    setSelectedPaths([])
    filePaths = []
  }

  function submit() {
    props.onChange(filePaths)
  }

  return (
    <div>
      <div class="display-flex not-first-child-margin-left-1">
        <button onClick={reset}>reset</button>
        <button onClick={submit}>submit</button>
      </div>
      <div class="display-flex not-first-child-margin-left-1">
        <FileSelector onChange={handleSelectedFile} />
        <FileSelector onChange={handleSelectedFolder} folder />
      </div>
      <ul>
        <For each={selectedPaths()}>{(path) => <li>{path}</li>}</For>
      </ul>
    </div>
  )
}

function getSelectedFolderPath(files) {
  const firstFolderPath = getFolderPath(files[0])
  const lastFolderPath = getFolderPath(files[files.length - 1])

  let prefix = ''

  for (let i = 0; i < firstFolderPath.length; i++) {
    if (firstFolderPath[i] === lastFolderPath[i]) {
      prefix += firstFolderPath[i]
    } else {
      break
    }
  }

  return prefix
}

function getFolderPath(file) {
  return file.path.replace(`\\${file.name}`, '')
}
