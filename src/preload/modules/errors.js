import {
  toResultObjectWithNullResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusOk,
  toResultObjectWithNullResultAndResultStatusPartiallyOk,
  toResultObjectWithResultStatusOk,
  toResultObjectWithResultStatusPartiallyOk
} from '../../preload/modules/resultStatus'

// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const inputError = Object.freeze({
  wrongFunctionArguments: 'wrong function arguments'
})

export class ErrorTracker {
  constructor() {
    this.errorCount = 0
    this.errorMessage = ''
  }

  concatErrorMessageOnNewLineAndIncrementErrorCount(errorMessage) {
    this.errorCount++
    this.errorMessage = `${this.errorMessage}\n${errorMessage}`
  }

  createResultObject(maxPossibleErrors, result) {
    if (this.errorCount === 0) {
      return result
        ? toResultObjectWithResultStatusOk(result)
        : toResultObjectWithNullResultAndResultStatusOk()
    } else if (this.errorCount > 0 && this.errorCount < maxPossibleErrors) {
      return result
        ? toResultObjectWithResultStatusPartiallyOk(result, this.errorMessage)
        : toResultObjectWithNullResultAndResultStatusPartiallyOk(this.errorMessage)
    } else {
      return toResultObjectWithNullResultAndResultStatusErrorSystem(this.errorMessage)
    }
  }

  isPartiallyOk(maxPossibleErrors) {
    return this.errorCount > 0 && this.errorCount < maxPossibleErrors
  }
}
