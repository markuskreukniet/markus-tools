import { exec } from 'child_process'

export default async function synchronizeDirectory(sourceDirectory, destinationDirectory) {
  const jsonArguments = JSON.stringify({
    sourceDirectory,
    destinationDirectory
  }).replace(/"/g, '\\"')
  const goProcess = exec(
    `go run ./go/main.go ./go/json_function_result.go ./go/file_utils.go ./go/synchronize_directory_trees.go "synchronizeDirectoryTreesToJSON" "${jsonArguments}"`,
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
