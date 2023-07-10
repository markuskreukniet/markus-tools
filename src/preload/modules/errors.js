// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const inputError = Object.freeze({
  wrongFunctionArguments: 'wrong function arguments'
})
