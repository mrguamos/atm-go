export function getEnumValues<T extends object>(enumObject: T): T[keyof T][] {
  return Object.values(enumObject).filter((value) => typeof value === 'number') as T[keyof T][]
}

export function getEnumKeys<T extends object>(enumObject: T): T[keyof T][] {
  return Object.keys(enumObject).filter((key) => typeof isNaN(Number(key))) as T[keyof T][]
}


export function generateStan() {
  const randomNumber = Math.floor(Math.random() * 1000000)
  const sixDigitNumber = randomNumber.toString().padStart(6, '0')
  return sixDigitNumber
}

export function generateRrn() {
  const randomNumber = Math.floor(Math.random() * 1000000)
  const sixDigitNumber = randomNumber.toString().padStart(12, '0')
  return sixDigitNumber
}

export function enumFromStringValue<T> (enm: { [s: string]: T}, value: string): T | undefined {
  return (Object.values(enm) as unknown as string[]).includes(value)
    ? value as unknown as T
    : undefined
}