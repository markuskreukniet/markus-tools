import { createSignal } from 'solid-js'
import DuplicateFiles from './components/pages/DuplicateFiles'
import ImagesToDateRangeFolder from './components/pages/ImagesToDateRangeFolder'
import LinesOfCode from './components/pages/LinesOfCode'
import ReferencesByUrls from './components/pages/ReferencesByUrls'
import SynchronizeDirectoryTrees from './components/pages/SynchronizeDirectoryTrees'
import Loader from './components/Loader'
import PageNavigator from './components/PageNavigator'
// import logo from './assets/logo.svg'

// TODO: check for error handling, createEffects, useless comments, useless changes (onChange) to parent element
// TODO: use batch, untrack, on (with { defer: true }) from import { batch, on, untrack } from "solid-js";?
// TODO: ImagesToDateRangeFolder: show how many files added/changed/removed

// TODO: rename folder to directory
// TODO: rename files.js to filePath.js
// TODO: No Symlink Handling, which results in bugs. Maybe is the right fix the use of symlinks (LStat)?

// TODO: filePathObjects and filePaths to fileSystemNodes, also the non arrays
// TODO: ResultObject or RO to either

// TODO: bug imagesToDateRangeFolder: selecting same input and output folder , then it does not remove an empty folder
// TODO: bug imagesToDateRangeFolder: Move images out of a date folder and the use the app again top create the same folder, then it wants to create the same folder, which it can't.

function App() {
  const [loading, setLoading] = createSignal(false)
  const navigationBarItemPageCombinations = [
    createNavigationBarItemPageCombination(DuplicateFiles, 'Duplicate Files'),
    createNavigationBarItemPageCombination(LinesOfCode, 'Lines of Code (LOC)'),
    createNavigationBarItemPageCombination(ReferencesByUrls, 'References by URLs'),
    createNavigationBarItemPageCombination(ImagesToDateRangeFolder, 'Images to Date Range Folder'),
    createNavigationBarItemPageCombination(SynchronizeDirectoryTrees, 'Synchronize Directory Trees')
  ]
  const activeNavigationBarItem = navigationBarItemPageCombinations[0].navigationBarItem

  function createNavigationBarItemPageCombination(Component, title) {
    return {
      navigationBarItem: title,
      page: <Component title={title} onLoading={setLoading} />
    }
  }

  return (
    <div class="container">
      <PageNavigator
        navigationBarItemPageCombinations={navigationBarItemPageCombinations}
        activeNavigationBarItem={activeNavigationBarItem}
      />
      <Loader loading={loading()} />

      {/* <img class="hero-logo" src={logo} alt="logo" />
      <h2 class="hero-text">You{"'"}ve successfully created an Electron project with Solid</h2>
      <p class="hero-tagline">
        Please try pressing <code>F12</code> to open the devTool
      </p>

      <div class="links">
        <div class="link-item">
          <a target="_blank" href="https://evite.netlify.app" rel="noopener noreferrer">
            Documentation
          </a>
        </div>
        <div class="link-item link-dot">•</div>
        <div class="link-item">
          <a
            target="_blank"
            href="https://github.com/alex8088/electron-vite"
            rel="noopener noreferrer"
          >
            Getting Help
          </a>
        </div>
        <div class="link-item link-dot">•</div>
        <div class="link-item">
          <a
            target="_blank"
            href="https://github.com/alex8088/quick-start/tree/master/packages/create-electron"
            rel="noopener noreferrer"
          >
            create-electron
          </a>
        </div>
      </div>

      <div class="features">
        <div class="feature-item">
          <article>
            <h2 class="title">Configuring</h2>
            <p class="detail">
              Config with <span>electron.vite.config.ts</span> and refer to the{' '}
              <a target="_blank" href="https://evite.netlify.app/config/" rel="noopener noreferrer">
                config guide
              </a>
              .
            </p>
          </article>
        </div>
        <div class="feature-item">
          <article>
            <h2 class="title">HMR</h2>
            <p class="detail">
              Edit <span>src/renderer</span> files to test HMR. See{' '}
              <a
                target="_blank"
                href="https://evite.netlify.app/guide/hmr-in-renderer.html"
                rel="noopener noreferrer"
              >
                docs
              </a>
              .
            </p>
          </article>
        </div>
        <div class="feature-item">
          <article>
            <h2 class="title">Hot Reloading</h2>
            <p class="detail">
              Run{' '}
              <span>
                {"'"}electron-vite dev --watch{"'"}
              </span>{' '}
              to enable. See{' '}
              <a
                target="_blank"
                href="https://evite.netlify.app/guide/hot-reloading.html"
                rel="noopener noreferrer"
              >
                docs
              </a>
              .
            </p>
          </article>
        </div>
        <div class="feature-item">
          <article>
            <h2 class="title">Debugging</h2>
            <p class="detail">
              Check out <span>.vscode/launch.json</span>. See{' '}
              <a
                target="_blank"
                href="https://evite.netlify.app/guide/debugging.html"
                rel="noopener noreferrer"
              >
                docs
              </a>
              .
            </p>
          </article>
        </div>
        <div class="feature-item">
          <article>
            <h2 class="title">Source Code Protection</h2>
            <p class="detail">
              Supported via built-in plugin <span>bytecodePlugin</span>. See{' '}
              <a
                target="_blank"
                href="https://evite.netlify.app/guide/source-code-protection.html"
                rel="noopener noreferrer"
              >
                docs
              </a>
              .
            </p>
          </article>
        </div>
        <div class="feature-item">
          <article>
            <h2 class="title">Packaging</h2>
            <p class="detail">
              Use{' '}
              <a target="_blank" href="https://www.electron.build" rel="noopener noreferrer">
                electron-builder
              </a>{' '}
              and pre-configured to pack your app.
            </p>
          </article>
        </div>
      </div> */}
    </div>
  )
}

export default App
