import { exec } from 'child_process'
import path from 'path'

// TODO: fix error handling
export default async function synchronizeDirectory(sourceDirectory, destinationDirectory) {
  const jsonArguments = JSON.stringify({
    sourceDirectory,
    destinationDirectory
  }).replace(/"/g, '\\"')
  const result = await stringsToGoFunctionCallWithArguments(
    'synchronizeDirectoryTreesToJSON',
    jsonArguments
  )
  // TODO: this log
  console.log(`Go program output: ${result}`)

  return `${sourceDirectory} testB ${destinationDirectory}`
}

async function stringsToGoFunctionCallWithArguments(functionCall, jsonArguments) {
  return new Promise((resolve, reject) => {
    const goProcess = exec(
      `go run . "${functionCall}" "${jsonArguments}"`,
      { cwd: path.join(__dirname, '..', '..', 'go') },
      (error, stdout) => {
        if (error) {
          console.error(`Error executing Go program: ${error}`)
          reject(error)
          return
        }
        resolve(stdout)
      }
    )
    // TODO: this log
    goProcess.on('close', (code) => {
      console.log(`Go program exited with code ${code}`)
    })
  })
}
