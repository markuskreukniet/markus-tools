import {
  toResultObjectWithNullResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusOk,
  toResultObjectWithNullResultAndResultStatusPartiallyOk,
  toResultObjectWithResultStatusOk,
  toResultObjectWithResultStatusPartiallyOk
} from '../../preload/modules/resultStatus'

// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const inputError = Object.freeze({
  aWrongCombinationOfArguments: 'A wrong combination of arguments.'
})

export function concatErrorMessageOnNewLine(currentErrorMessage, errorMessage) {
  return `${currentErrorMessage}\n${errorMessage}`
}

export class ErrorTracker {
  constructor(maxPossibleErrors) {
    this.errorCount = 0
    this.errorMessage = ''
    this.maxPossibleErrors = maxPossibleErrors || 0
  }

  addNumberOfPossibleErrors(numberOfPossibleErrors) {
    this.maxPossibleErrors = this.maxPossibleErrors + numberOfPossibleErrors
  }

  concatErrorMessageOnNewLineAndIncrementErrorCount(errorMessage) {
    this.errorCount++
    this.errorMessage = concatErrorMessageOnNewLine(this.errorMessage, errorMessage)
  }

  createResultObject(result) {
    if (this.errorCount === 0) {
      return result
        ? toResultObjectWithResultStatusOk(result)
        : toResultObjectWithNullResultAndResultStatusOk()
    } else if (this.errorCount > 0 && this.errorCount < this.maxPossibleErrors) {
      return result
        ? toResultObjectWithResultStatusPartiallyOk(result, this.errorMessage)
        : toResultObjectWithNullResultAndResultStatusPartiallyOk(this.errorMessage)
    } else {
      return toResultObjectWithNullResultAndResultStatusErrorSystem(this.errorMessage)
    }
  }
}
