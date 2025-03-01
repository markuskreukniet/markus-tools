import { createSignal } from 'solid-js'
import DuplicateFiles from './components/pages/DuplicateFiles'
import ImagesToDateRangeDirectory from './components/pages/ImagesToDateRangeDirectory'
import LinesOfCode from './components/pages/LinesOfCode'
import PlainTextFilesToText from './components/pages/PlainTextFilesToText'
import ReferencesByUrls from './components/pages/ReferencesByUrls'
import SynchronizeDirectoryTrees from './components/pages/SynchronizeDirectoryTrees'
import LoadingSpinner from './components/LoadingSpinner'
import PageNavigator from './components/PageNavigator'
// import logo from './assets/logo.svg'

// TODO: check for error handling, createEffects, useless comments, useless changes (onChange) to parent element
// TODO: In JavaScript, when a function modifies an array, it is also modified outside the function, refactor code with this logic.
// TODO: use batch, untrack, on (with { defer: true }) from import { batch, on, untrack } from "solid-js";?
// TODO: ImagesToDateRangeDirectory: show how many files added/changed/removed
// TODO: use min-width: 0? https://www.youtube.com/watch?v=cH8VbLM1958&t=4s
// TODO: fix font-size? https://www.youtube.com/watch?v=rg3zgQ3xBRc
// TODO: when to use rem and when px(see description) https://www.youtube.com/watch?v=xCSw6bPXZks
// TODO: use onMouseDown? https://www.youtube.com/watch?v=yaMGtiPckAQ

// TODO: rename folder to directory
// TODO: rename files.js to filePath.js
// TODO: No Symlink Handling, which results in bugs. Maybe is the right fix the use of symlinks (LStat)?

// TODO: filePathObjects and filePaths to fileSystemNodes, also the non arrays
// TODO: ResultObject or RO to either

// TODO: bug imagesToDateRangeDirectory: selecting same input and output directory, then it does not remove an empty directory
// TODO: bug imagesToDateRangeDirectory: Move images out of a date directory and the use the app again top create the same directory, then it wants to create the same directory, which it can't.

function App() {
  const [loading, setLoading] = createSignal(false)
  const navigationBarItemPageCombinations = [
    createNavigationBarItemPageCombination(DuplicateFiles, 'Duplicate Files'),
    createNavigationBarItemPageCombination(LinesOfCode, 'Lines of Code (LOC)'),
    createNavigationBarItemPageCombination(ReferencesByUrls, 'References by URLs'),
    createNavigationBarItemPageCombination(
      ImagesToDateRangeDirectory,
      'Images to Date Range Directory'
    ),
    createNavigationBarItemPageCombination(
      SynchronizeDirectoryTrees,
      'Synchronize Directory Trees'
    ),
    createNavigationBarItemPageCombination(PlainTextFilesToText, 'Plain Text Files to Text')
  ]
  const activeNavigationBarItem = navigationBarItemPageCombinations[0].navigationBarItem

  function createNavigationBarItemPageCombination(Component, title) {
    return {
      navigationBarItem: title,
      page: <Component title={title} onLoading={setLoading} />
    }
  }

  return (
    <div id="app">
      <PageNavigator
        navigationBarItemPageCombinations={navigationBarItemPageCombinations}
        activeNavigationBarItem={activeNavigationBarItem}
      />
      <LoadingSpinner loading={loading()} />
    </div>
  )
}

export default App
