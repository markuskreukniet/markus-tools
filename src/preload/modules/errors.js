import {
  toResultObjectWithNullResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusOk,
  toResultObjectWithNullResultAndResultStatusPartiallyOk
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

  toResultObjectWithNullResult(maxPossibleErrors) {
    if (this.errorCount === 0) {
      return toResultObjectWithNullResultAndResultStatusOk()
    } else if (this.errorCount > 0 && this.errorCount < maxPossibleErrors) {
      return toResultObjectWithNullResultAndResultStatusPartiallyOk(this.errorMessage)
    } else {
      return toResultObjectWithNullResultAndResultStatusErrorSystem(this.errorMessage)
    }
  }
}
