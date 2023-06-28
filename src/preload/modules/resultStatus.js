// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const resultStatus = Object.freeze({
  ok: 'ok',
  errorSystem: 'errorSystem'
})

export function toResultObject(result, status, message) {
  return { result: result ? result : null, status, message: message ? message : '' }
}
