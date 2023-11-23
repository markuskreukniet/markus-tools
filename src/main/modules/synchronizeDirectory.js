import { exec } from 'child_process'
import path from 'path'

export default async function synchronizeDirectory(sourceDirectory, destinationDirectory) {
  const goDir = path.join(__dirname, '..', '..', 'go')

  const jsonArguments = JSON.stringify({
    sourceDirectory,
    destinationDirectory
  }).replace(/"/g, '\\"')
  const goProcess = exec(
    `go run . "synchronizeDirectoryTreesToJSON" "${jsonArguments}"`,
    { cwd: goDir },
    (error, stdout) => {
      if (error) {
        console.error(`Error executing Go program: ${error}`)
        return
      }
      console.log(`Go program output: ${stdout}`)
    }
  )

  goProcess.on('close', (code) => {
    console.log(`Go program exited with code ${code}`)
  })

  return `${sourceDirectory} testB ${destinationDirectory}`
}
