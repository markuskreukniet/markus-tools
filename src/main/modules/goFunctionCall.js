import { exec } from 'child_process'
import path from 'path'
import { toEitherLeftResult, toEitherRightResult } from '../../preload/monads/either'

export default async function goFunctionCall(functionName, argumentObject) {
  const result = JSON.parse(
    await toGoFunctionCall(
      functionName,
      JSON.stringify(argumentObject).replace(/"/g, '\\"') // TODO: does the replace work on systems besides Windows?
    )
  )
  if (result.ErrorMessage === '') {
    return toEitherRightResult(result.Result)
  } else {
    return toEitherLeftResult(result.ErrorMessage)
  }
}

async function toGoFunctionCall(functionCall, jsonArguments) {
  return new Promise((resolve, reject) => {
    let result = ''
    let goProcess = null
    const commandPart = `"${functionCall}" "${jsonArguments}"`

    // Use 'extraResources' in electron-builder (electron-builder.yml) to include the 'bin' directory (containing binaries such as .exe files) in the 'resources' directory.
    // In this case, we should prefer 'extraResources' instead of 'asarUnpack' or 'extraFiles,' at least for these reasons:
    // - Binaries need to be accessible outside the '.asar' archive. While 'asarUnpack' can unpack specific files from the '.asar' archive, it is unnecessary for executables.
    // - 'extraResources' automatically places files in the 'resources' directory,
    //    which is where Electron expects additional runtime dependencies, and are easily accessed using `process.resourcesPath`, which always points to the 'resources' directory.
    // - Using 'extraFiles' would place the files in the root directory (e.g., `win-unpacked/bin`), cluttering the app's root and making it harder to manage.
    // - Storing binaries in `resources/bin` keeps the app organized and follows Electron's convention for runtime dependencies.

    // These lines should be part of the 'extraResources' part in electron-builder.yml:
    // extraResources:
    //  - from: 'out/bin'
    //    to: 'bin'

    // When we run the 'npm start' command, 'preview' is active, which is scripts.start from package.json.
    // When we run a 'run build:' command, 'production' is active, which is build:win, build:mac, or build:linux from package.json scripts.
    // Adding the parts '--mode preview' and '--mode production' to these scripts was needed.
    switch (import.meta.env.MODE) {
      case 'preview':
        goProcess = exec(`go run . ${commandPart}`, {
          cwd: path.join(__dirname, '..', '..', 'go')
        })
        break
      case 'production':
        goProcess = exec(
          `"${path.join(process.resourcesPath, 'bin', 'markus-tools-go.exe')}" ${commandPart}`
        )
        break
      default:
      // TODO: error
    }

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
