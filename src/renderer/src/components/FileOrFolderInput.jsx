import FileSelector from './FileSelector'

import { createSignal, For } from 'solid-js'

export default function FileOrFolderInput(props) {
  const [selectedPaths, setSelectedPaths] = createSignal([])

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
    const filePaths = []
    for (const file of files) {
      filePaths.push(file.path)
    }

    props.onChange(filePaths)
  }

  return (
    <div>
      <div>
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
