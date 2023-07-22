// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const inputError = Object.freeze({
  wrongFunctionArguments: 'wrong function arguments',
  couldNotReadOrRemoveADirectory: 'Could not read or remove a directory'
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
}
