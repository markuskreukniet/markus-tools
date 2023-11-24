import { exec } from 'child_process'
import path from 'path'

// TODO: fix errors, on close, and use stringsToGoFunctionCallWithArguments result
export default async function synchronizeDirectory(sourceDirectory, destinationDirectory) {
  const jsonArguments = JSON.stringify({
    sourceDirectory,
    destinationDirectory
  }).replace(/"/g, '\\"')
  await stringsToGoFunctionCallWithArguments('synchronizeDirectoryTreesToJSON', jsonArguments)

  return `${sourceDirectory} testB ${destinationDirectory}`
}

async function stringsToGoFunctionCallWithArguments(functionCall, jsonArguments) {
  await new Promise((resolve, reject) => {
    const goProcess = exec(
      `go run . "${functionCall}" "${jsonArguments}"`,
      { cwd: path.join(__dirname, '..', '..', 'go') },
      (error, stdout) => {
        if (error) {
          console.error(`Error executing Go program: ${error}`)
          reject(error)
          return
        }
        console.log(`Go program output: ${stdout}`)
        resolve(stdout)
      }
    )
    goProcess.on('close', (code) => {
      console.log(`Go program exited with code ${code}`)
    })
  })
}
