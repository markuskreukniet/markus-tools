// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const resultStatus = Object.freeze({
  ok: 'ok',
  errorSystem: 'errorSystem'
})

export function getResultStatusCombination(filePaths, status) {
  return { result: filePaths, status }
}
