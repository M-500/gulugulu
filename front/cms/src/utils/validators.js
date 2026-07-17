export function isEmail(value) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)
}

export function isStrongEnoughPassword(value) {
  return typeof value === 'string' && value.length >= 6
}
