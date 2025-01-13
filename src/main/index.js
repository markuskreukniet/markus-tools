import { app, shell, BrowserWindow, ipcMain } from 'electron'
import { join } from 'path'
import { electronApp, optimizer, is } from '@electron-toolkit/utils'
import icon from '../../resources/icon.png?asset'

function createWindow() {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    // width: 900,
    // height: 670,
    width: 752, // (34 x 16) + (13 x 16)
    height: 672, // (55 x 16) - (13 x 16)
    show: false,
    autoHideMenuBar: true,
    ...(process.platform === 'linux' ? { icon } : {}),
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false
    }
  })

  mainWindow.on('ready-to-show', () => {
    mainWindow.show()
  })

  mainWindow.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: 'deny' }
  })

  // HMR for renderer base on electron-vite cli.
  // Load the remote URL for development or the local html file for production.
  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    mainWindow.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    mainWindow.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  // Set app user model id for windows
  electronApp.setAppUserModelId('com.electron')

  // Default open or close DevTools by F12 in development
  // and ignore CommandOrControl + R in production.
  // see https://github.com/alex8088/electron-toolkit/tree/master/packages/utils
  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  createWindow()

  app.on('activate', function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// In this file you can include the rest of your app"s specific main process
// code. You can also put them in separate files and require them here.

// self added with ipcMain import
import duplicateFiles from './modules/duplicateFiles'
async function duplicateFilesBE(_, fileSystemNodes) {
  return duplicateFiles(fileSystemNodes)
}
ipcMain.handle('duplicateFilesBE', duplicateFilesBE)

import linesOfCode from './modules/linesOfCode'
async function linesOfCodeBE(_, fileSystemNodes) {
  return linesOfCode(fileSystemNodes)
}
ipcMain.handle('linesOfCodeBE', linesOfCodeBE)

import referencesByUrls from './modules/referencesByUrls'
async function referencesByUrlsBE(_, urlsString) {
  return referencesByUrls(urlsString)
}
ipcMain.handle('referencesByUrlsBE', referencesByUrlsBE)

import imagesToDateRangeFolder from './modules/imagesToDateRangeFolder'
async function imagesToDateRangeFolderBE(_, fileSystemNodes, filePath) {
  return imagesToDateRangeFolder(fileSystemNodes, filePath)
}
ipcMain.handle('imagesToDateRangeFolderBE', imagesToDateRangeFolderBE)

import goFunctionCall from './modules/goFunctionCall'
async function goFunctionCallBE(_, functionName, argumentObject) {
  return goFunctionCall(functionName, argumentObject)
}
ipcMain.handle('goFunctionCallBE', goFunctionCallBE)

import selectFilePathDialog from './modules/selectFilePathDialog'
async function selectFilePathDialogBE(_, selectFolder) {
  return selectFilePathDialog(selectFolder)
}
ipcMain.handle('selectFilePathDialogBE', selectFilePathDialogBE)
