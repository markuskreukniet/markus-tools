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
    const goProcess = exec(`go run . "${functionCall}" "${jsonArguments}"`, {
      cwd: path.join(__dirname, '..', '..', 'go')
    })

    // build with: go build -o ../out/go/markus-tools-go.exe
    // use this exec:
    // const goProcess = exec(
    //   `"${path.join(
    //     __dirname,
    //     '..',
    //     'go',
    //     'markus-tools-go.exe'
    //   )}" "${functionCall}" "${jsonArguments}"`
    // )

    // add to electron-builder.yml:
    // asarUnpack:
    //   - 'out/go/markus-tools-go.exe'
    // build with: npm run build:win
    // use this exec:
    // const goProcess = exec(
    //   `"${path.join(
    //     process.resourcesPath,
    //     'app.asar.unpacked',
    //     'out',
    //     'go',
    //     'markus-tools-go.exe'
    //   )}" "${functionCall}" "${jsonArguments}"`
    // )

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
