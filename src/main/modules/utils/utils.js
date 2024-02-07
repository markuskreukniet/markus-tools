import { exec } from 'child_process'
import path from 'path'
import { toEitherLeftResult, toEitherRightResult } from '../../../preload/monads/either'

export async function objectToJSONAndStringsToGoFunctionCallWithArguments(
  functionCall,
  nonJSONObject
) {
  const result = JSON.parse(
    await stringsToGoFunctionCallWithArguments(
      functionCall,
      JSON.stringify(nonJSONObject).replace(/"/g, '\\"') // TODO: does the replace work on systems besides Windows?
    )
  )
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

    // build with: go build -o ../out/go/markus-tools.exe
    // use this exec
    // const goProcess = exec(
    //   `"${path.join(
    //     __dirname,
    //     '..',
    //     'go',
    //     'markus-tools.exe'
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
