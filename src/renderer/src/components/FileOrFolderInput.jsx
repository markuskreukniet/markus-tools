import FileSelector from './FileSelector'

import { createSignal, For } from 'solid-js'

// TODO: do not add duplicate selectedPaths or filePaths
// If a selectedPaths is duplicate, don't add its filePaths
// If a selectedPath is a file part of an already selected folder, don't add that file
// If a selectedPath is a folder and a file part of that folder is already added, remove the file and add the folder

// Storing a file path is better than storing a file name since we don't have to combine a folder path, and a file name (combining the strings is less efficient)
export default function FileOrFolderInput(props) {
  const [folderFilePathCombinations, setFolderFilePathCombinations] = createSignal([])
  let filePaths = []

  function handleSelectedFile(files) {
    const folderPath = getSelectedFolderPath(files)
    if (
      !folderFilePathCombinations().some(
        (combination) =>
          (combination.folderPath === folderPath && combination.filePath === null) ||
          combination.filePath === files[0].path
      )
    ) {
      const combination = createCombination(folderPath, files[0].path)
      setFolderFilePathCombinations([...folderFilePathCombinations(), combination])
      handleFilePaths(files)
    }
  }

  // TODO
  function handleSelectedFolder(files) {
    const folderPath = getSelectedFolderPath(files)
    if (
      !folderFilePathCombinations().some(
        (combination) => combination.folderPath === folderPath && combination.filePath === null
      )
    ) {
      const combination = createCombination(folderPath, null)
      setFolderFilePathCombinations([...folderFilePathCombinations(), combination])
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
    setFolderFilePathCombinations([])
    filePaths = []
  }

  function submit() {
    props.onChange(filePaths)
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
        <For each={folderFilePathCombinations()}>
          {(combination) => <li>{combination.filePath}</li>}
        </For>
      </ul>
    </div>
  )
}

function createCombination(folderPath, filePath) {
  return {
    folderPath: folderPath,
    filePath: filePath
  }
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
