import { contextBridge, ipcRenderer } from 'electron'
import { electronAPI } from '@electron-toolkit/preload'

// Custom APIs for renderer
const api = {}

// Use `contextBridge` APIs to expose Electron APIs to
// renderer only if context isolation is enabled, otherwise
// just add to the DOM global.
if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld('electron', electronAPI)
    contextBridge.exposeInMainWorld('api', api)
  } catch (error) {
    console.error(error)
  }
} else {
  window.electron = electronAPI
  window.api = api
}

// self added with ipcRenderer import
contextBridge.exposeInMainWorld('duplicateFiles', {
  getDuplicateFiles: (filePaths) => ipcRenderer.invoke('getDuplicateFiles', filePaths)
})

contextBridge.exposeInMainWorld('codeQuality', {
  linesOfCodeBE: (filePaths) => ipcRenderer.invoke('linesOfCodeBE', filePaths)
})

contextBridge.exposeInMainWorld('references', {
  getReferencesByUrls: (urlsString) => ipcRenderer.invoke('getReferencesByUrls', urlsString)
})

contextBridge.exposeInMainWorld('dateRangeFolder', {
  imagesToDateRangeFolderBE: (filePathObjects, path, useDirectoriesTreeInput) =>
    ipcRenderer.invoke('imagesToDateRangeFolderBE', filePathObjects, path, useDirectoriesTreeInput)
})

contextBridge.exposeInMainWorld('synchronization', {
  synchronizeDirectoryBE: (originalDirectoryFilePath, destinationDirectoryFilePath) =>
    ipcRenderer.invoke(
      'synchronizeDirectoryBE',
      originalDirectoryFilePath,
      destinationDirectoryFilePath
    )
})

contextBridge.exposeInMainWorld('dialog', {
  openFileDialogBE: (selectFolder) => ipcRenderer.invoke('openFileDialogBE', selectFolder)
})
