// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const resultStatus = Object.freeze({
  ok: 'ok',
  partiallyOk: 'partiallyOk',
  errorSystem: 'errorSystem'
})

export function toResultObject(result, status, message) {
  return { result, status, message: message ? message : '' }
}

export function toResultObjectWithNullResult(status, message) {
  return toResultObject(null, status, message)
}

export function toResultObjectWithNullResultByResultObject(resultObject) {
  return toResultObject(null, resultObject.status, resultObject.message)
}

export function isResultObjectOk(resultObject) {
  return resultObject.status === resultStatus.ok
}

export function isResultObjectPartiallyOk(resultObject) {
  return resultObject.status === resultStatus.partiallyOk
}
