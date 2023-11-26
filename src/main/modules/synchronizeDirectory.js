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
