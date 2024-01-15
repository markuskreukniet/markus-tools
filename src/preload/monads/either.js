export const Either = {
  left: (value) => new Left(value),
  right: (value) => new Right(value)
}

class Left {
  constructor(value) {
    this.value = value
  }
  isLeft() {
    return true
  }
  isRight() {
    return false
  }
}

class Right {
  constructor(value) {
    this.value = value
  }
  isLeft() {
    return false
  }
  isRight() {
    return true
  }
}

export const eitherType = Object.freeze({
  left: 'left',
  right: 'right'
})

export function toEitherLeftResult(value) {
  return { type: eitherType.left, value }
}

export function toEitherRightResult(value) {
  return { type: eitherType.right, value }
}

export function isEitherRightResult(eitherResult) {
  return eitherResult.type === eitherType.right
}

export function eitherLeftResultToErrorString(eitherResult) {
  return `error: ${eitherResult.value}`
}

export function getEitherResultValueOrEitherResultToErrorString(eitherResult) {
  if (isEitherRightResult(eitherResult)) {
    return eitherResult.value
  } else {
    return eitherLeftResultToErrorString(eitherResult)
  }
}
