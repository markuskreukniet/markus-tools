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
  imagesToDateRangeFolderBE: (fileSystemNodes, filePath) =>
    ipcRenderer.invoke('imagesToDateRangeFolderBE', fileSystemNodes, filePath)
})

contextBridge.exposeInMainWorld('goBackend', {
  goFunctionCallBE: (functionName, argumentObject) =>
    ipcRenderer.invoke('goFunctionCallBE', functionName, argumentObject)
})

contextBridge.exposeInMainWorld('dialog', {
  selectFilePathDialogBE: (selectFolder) =>
    ipcRenderer.invoke('selectFilePathDialogBE', selectFolder)
})
