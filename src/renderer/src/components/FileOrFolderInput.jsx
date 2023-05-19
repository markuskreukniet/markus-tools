import { createSignal, For } from 'solid-js'
import ActiveByNumberButton from './ActiveByNumberButton'
import FileSelector from './FileSelector'

// TODO:
// Adding a file could add a duplicate file since there could already be a folder with its whole tree of child folders already containing that file.
// Adding a folder could add a duplicate file since the folder with its whole tree of child folders could contain a duplicate file.

// Checking child folders of a folder is only possible in the main, which is possible by adding such a function in the main.

export default function FileOrFolderInput(props) {
  const [selectedPaths, setSelectedPaths] = createSignal([])
  const [numberOfFilePaths, setNumberOfFilePaths] = createSignal(0)
  let filePaths = []

  function handleSelectedFile(files) {
    if (!selectedPaths().some((path) => path === files[0].path)) {
      setState(files[0].path, files)
    }
  }

  function handleSelectedFolder(files) {
    const folderPath = getSelectedFolderPath(files)

    if (!selectedPaths().some((path) => path === folderPath)) {
      setState(folderPath, files)
    }
  }

  function setState(selectedPath, files) {
    setSelectedPaths([...selectedPaths(), selectedPath])

    // files is a FileList, not an array, so we can't use .map
    for (const file of files) {
      filePaths.push(file.path)
    }

    setNumberOfFilePaths(filePaths.length)
  }

  function resetState() {
    setSelectedPaths([])
    filePaths = []
    setNumberOfFilePaths(0)
  }

  function submit() {
    props.onChange(filePaths)
  }

  return (
    <div>
      <div class="display-flex not-first-child-margin-left-1">
        <FileSelector onChange={handleSelectedFile} />
        <FileSelector onChange={handleSelectedFolder} folder />
      </div>
      <div class="display-flex justify-content-flex-end not-first-child-margin-left-1">
        <ActiveByNumberButton
          minimumNumber={1}
          currentNumber={numberOfFilePaths()}
          onAction={resetState}
          text="reset"
        />
        <ActiveByNumberButton
          minimumNumber={props.minimumFiles}
          currentNumber={numberOfFilePaths()}
          onAction={submit}
          text="submit"
        />
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
