import { createSignal } from 'solid-js'
import DuplicateFiles from './components/pages/DuplicateFiles'
import ImagesToDateRangeFolder from './components/pages/ImagesToDateRangeFolder'
import LinesOfCode from './components/pages/LinesOfCode'
import ReferencesByUrls from './components/pages/ReferencesByUrls'
import SynchronizeDirectory from './components/pages/SynchronizeDirectory'
import Loader from './components/Loader'
import PageNavigator from './components/PageNavigator'
// import logo from './assets/logo.svg'

// TODO: check for good error handling whole app
// TODO: remove some createEffects
// TODO: remove useless comments
// TODO: rename folder to directory
// TODO: rename files.js to filePath.js
// TODO: rename openFileDialog to openFilePathDialog.js with functions
// TODO: imagesToDateRangeFolder check input and output folder is not the same or a child

// TODO: backend functions should return a object for:
// function tagsAndMessageToJsx(openingTag, closingTag, message) {
//   if (openingTag === 'ok-result' && closingTag === '</p>') {
//     return <p class="">{message}</p>
//   }
// }
function App() {
  const [loading, setLoading] = createSignal(false)
  const duplicateFilesTitle = 'Duplicate Files'
  const linesOfCodeTitle = 'Lines of Code (LOC)'
  const referencesByUrlsTitle = 'References by URLs'
  const imagesToDateRangeFolderTitle = 'Images to Date Range Folder'
  const synchronizeDirectoryTitle = 'Synchronize Directory'
  // TODO: make of abstraction of navigationBarItem and page
  const navigationBarItemPageCombinations = [
    {
      navigationBarItem: duplicateFilesTitle,
      page: <DuplicateFiles title={duplicateFilesTitle} onLoading={setLoading} />
    },
    {
      navigationBarItem: linesOfCodeTitle,
      page: <LinesOfCode title={linesOfCodeTitle} onLoading={setLoading} />
    },
    {
      navigationBarItem: referencesByUrlsTitle,
      page: <ReferencesByUrls title={referencesByUrlsTitle} onLoading={setLoading} />
    },
    {
      navigationBarItem: imagesToDateRangeFolderTitle,
      page: <ImagesToDateRangeFolder title={imagesToDateRangeFolderTitle} onLoading={setLoading} />
    },
    {
      navigationBarItem: synchronizeDirectoryTitle,
      page: <SynchronizeDirectory title={synchronizeDirectoryTitle} onLoading={setLoading} />
    }
  ]
  const activeNavigationBarItem = navigationBarItemPageCombinations[0].navigationBarItem

  return (
    <div class="container">
      <PageNavigator
        navigationBarItemPageCombinations={navigationBarItemPageCombinations}
        activeNavigationBarItem={activeNavigationBarItem}
      />
      <Loader loading={loading} />

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
