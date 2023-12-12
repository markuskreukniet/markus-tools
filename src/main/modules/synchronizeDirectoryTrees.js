import { toEitherLeftResult, toEitherRightResult } from '../../preload/monads/either'
import { exec } from 'child_process'
import path from 'path'

export default async function synchronizeDirectoryTrees(sourceDirectory, destinationDirectory) {
  // TODO: does the replace work on systems besides Windows?
  const jsonArguments = JSON.stringify({
    sourceDirectory,
    destinationDirectory
  }).replace(/"/g, '\\"')
  const result = JSON.parse(
    await stringsToGoFunctionCallWithArguments('synchronizeDirectoryTreesToJSON', jsonArguments)
  )

  // TODO: make function for this
  if (result.ErrorMessage === '') {
    return toEitherRightResult(result.Result)
  } else {
    return toEitherLeftResult(result.ErrorMessage)
  }
}

async function stringsToGoFunctionCallWithArguments(functionCall, jsonArguments) {
  return new Promise((resolve, reject) => {
    let result = ''
    const goProcess = exec(`go run . "${functionCall}" "${jsonArguments}"`, {
      cwd: path.join(__dirname, '..', '..', 'go')
    })
    goProcess.stdout.on('data', (data) => {
      result += data
    })
    goProcess.on('error', (error) => {
      reject(error)
    })
    goProcess.on('close', (code) => {
      if (code !== 0) {
        reject(new Error(`Go process exited with code ${code}`))
      } else {
        resolve(result)
      }
    })
  })
}
