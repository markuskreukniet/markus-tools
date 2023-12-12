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
  duplicateFilesBE: (fileSystemNodes) => ipcRenderer.invoke('duplicateFilesBE', fileSystemNodes)
})

contextBridge.exposeInMainWorld('codeQuality', {
  linesOfCodeBE: (fileSystemNodes) => ipcRenderer.invoke('linesOfCodeBE', fileSystemNodes)
})

contextBridge.exposeInMainWorld('references', {
  referencesByUrlsBE: (urlsString) => ipcRenderer.invoke('referencesByUrlsBE', urlsString)
})

contextBridge.exposeInMainWorld('dateRangeFolder', {
  imagesToDateRangeFolderBE: (fileSystemNodes, path, useDirectoriesTreeInput) =>
    ipcRenderer.invoke('imagesToDateRangeFolderBE', fileSystemNodes, path, useDirectoriesTreeInput)
})

contextBridge.exposeInMainWorld('synchronization', {
  synchronizeDirectoryTreesBE: (originalDirectoryFilePath, destinationDirectoryFilePath) =>
    ipcRenderer.invoke(
      'synchronizeDirectoryTreesBE',
      originalDirectoryFilePath,
      destinationDirectoryFilePath
    )
})

contextBridge.exposeInMainWorld('dialog', {
  openFileDialogBE: (selectFolder) => ipcRenderer.invoke('openFileDialogBE', selectFolder)
})
