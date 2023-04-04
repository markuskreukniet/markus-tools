import { createSignal, For } from 'solid-js'

export default function FileSelector(props) {
  const [selectedPaths, setSelectedPaths] = createSignal([])

  function handleSelectedFileOrFolder(files) {
    // add selected file or folder to the ul
    const pathToAdd = files.length === 1 ? files[0].path : getSelectedFolderPath(files)
    setSelectedPaths([...selectedPaths(), pathToAdd])

    // convert files to file paths
    // files is a FileList, not an array, so we can't use .map
    const filePaths = []
    for (const file of files) {
      filePaths.push(file.path)
    }

    props.onChange(filePaths)
  }

  return (
    <div>
      <input
        type="file"
        webkitdirectory={props.folder}
        onClick={
          (e) => (e.target.value = '') /* makes selecting the same file or folder possible */
        }
        onChange={(e) => handleSelectedFileOrFolder(e.target.files)}
      />
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
