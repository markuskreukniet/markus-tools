// We can't use symbols across the Electron IPC (inter-process communication) boundary
const resultStatus = Object.freeze({
  ok: 'ok',
  errorSystem: 'errorSystem'
})

export default resultStatus
