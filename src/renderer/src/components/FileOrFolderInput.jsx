import FileSelector from './FileSelector'

import { createSignal, For } from 'solid-js'

// TODO: do not add duplicate selectedPaths or filePaths
// If a selectedPaths is duplicate, don't add its filePaths
// If a selectedPath is a file part of an already selected folder, don't add that file
// If a selectedPath is a folder and a file part of that folder is already added, remove the file and add the folder
export default function FileOrFolderInput(props) {
  const [selectedPaths, setSelectedPaths] = createSignal([])
  const [filePaths, setFilePaths] = createSignal([])

  function handleSelectedFile(files) {
    setSelectedPaths([...selectedPaths(), files[0].path])
    handleFilePaths(files)
  }

  function handleSelectedFolder(files) {
    setSelectedPaths([...selectedPaths(), getSelectedFolderPath(files)])
    handleFilePaths(files)
  }

  function handleFilePaths(files) {
    // files is a FileList, not an array, so we can't use .map
    const newFilePaths = []
    for (const file of files) {
      newFilePaths.push(file.path)
    }

    setFilePaths([...filePaths(), ...newFilePaths])
  }

  function reset() {
    setSelectedPaths([])
    setFilePaths([])
  }

  function submit() {
    props.onChange(filePaths())
  }

  return (
    <div>
      <div>
        <FileSelector onChange={handleSelectedFile} />
        <FileSelector onChange={handleSelectedFolder} folder />
      </div>
      <div>
        <button onClick={reset}>reset</button>
        <button onClick={submit}>submit</button>
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
