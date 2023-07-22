// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const inputError = Object.freeze({
  wrongFunctionArguments: 'wrong function arguments',
  couldNotReadOrRemoveADirectory: 'Could not read or remove a directory'
})

export function ErrorTracker() {
  const that = this
  let errorCount = 0
  let errorMessage = ''

  this.concatErrorMessageOnNewLineAndIncrementErrorCount = function (errorMessage) {
    that.errorCount++
    that.errorMessage = `${that.errorMessage}\n${errorMessage}`
  }

  // TODO: getInstance() needed?
}
