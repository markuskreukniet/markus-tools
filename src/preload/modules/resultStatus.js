// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const resultStatus = Object.freeze({
  ok: 'ok',
  partiallyOk: 'partiallyOk',
  errorSystem: 'errorSystem'
})

// toResultObject
export function toResultObject(result, status, message) {
  return { result, status, message: message || '' }
}

// toResultObjectWithNullResultAndResultStatus
export function toResultObjectWithNullResultAndResultStatusOk(message) {
  return toResultObject(null, resultStatus.ok, message)
}

export function toResultObjectWithNullResultAndResultStatusErrorSystem(message) {
  return toResultObject(null, resultStatus.errorSystem, message)
}

export function toResultObjectWithNullResultAndResultStatusPartiallyOk(message) {
  return toResultObject(null, resultStatus.partiallyOk, message)
}

// toResultObjectWithResultStatus
export function toResultObjectWithResultStatusOk(result, message) {
  return toResultObject(result, resultStatus.ok, message)
}

export function toResultObjectWithResultStatusErrorSystem(result, message) {
  return toResultObject(result, resultStatus.errorSystem, message)
}

export function toResultObjectWithResultStatusPartiallyOk(result, message) {
  return toResultObject(result, resultStatus.partiallyOk, message)
}

// toResultObjectWithEmptyArrayResultAndResultStatus
export function toResultObjectWithEmptyArrayResultAndResultStatusOk(message) {
  return toResultObject([], resultStatus.ok, message)
}

export function toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(message) {
  return toResultObject([], resultStatus.errorSystem, message)
}

export function toResultObjectWithEmptyArrayResultAndResultStatusPartiallyOk(message) {
  return toResultObject([], resultStatus.partiallyOk, message)
}

// isResultObject
export function isResultObjectOk(resultObject) {
  return resultObject.status === resultStatus.ok
}

export function isResultObjectPartiallyOk(resultObject) {
  return resultObject.status === resultStatus.partiallyOk
}
