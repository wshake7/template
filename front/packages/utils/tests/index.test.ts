import { expect, test } from 'vite-plus/test'
import { fnTest1 } from '../src/index.ts'

test('fnTest', () => {
  expect(fnTest1()).toBe('Hello, tsdown!')
})
