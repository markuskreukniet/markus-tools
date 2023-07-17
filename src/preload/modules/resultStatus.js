// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const resultStatus = Object.freeze({
  ok: 'ok',
  partiallyOk: 'partiallyOk',
  errorSystem: 'errorSystem'
})

// toResultObject
export function toResultObject(result, status, message) {
  return { result, status, message: message ? message : '' }
}

// TODO: remove
export function toResultObjectWithNullResult(status, message) {
  return toResultObject(null, status, message)
}

export function toResultObjectWithNullResultAndResultStatusOk(message) {
  return toResultObject(null, resultStatus.ok, message)
}

export function toResultObjectWithNullResultAndResultStatusErrorSystem(message) {
  return toResultObject(null, resultStatus.errorSystem, message)
}

export function toResultObjectWithNullResultAndResultStatusPartiallyOk(message) {
  return toResultObject(null, resultStatus.partiallyOk, message)
}

export function toResultObjectWithNullResultByResultObject(resultObject) {
  return toResultObject(null, resultObject.status, resultObject.message)
}

export function toResultObjectWithResultStatusOk(result, message) {
  return toResultObject(result, resultStatus.ok, message)
}

// isResultObject
export function isResultObjectOk(resultObject) {
  return resultObject.status === resultStatus.ok
}

export function isResultObjectPartiallyOk(resultObject) {
  return resultObject.status === resultStatus.partiallyOk
}
