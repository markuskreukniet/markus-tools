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
