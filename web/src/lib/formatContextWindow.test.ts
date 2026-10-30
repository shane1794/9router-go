import { describe, expect, it } from 'bun:test'
import { formatContextWindow } from './formatContextWindow'

describe('formatContextWindow', () => {
  it('whole millions read without a decimal', () => {
    expect(formatContextWindow(1_000_000)).toBe('1M')
    expect(formatContextWindow(2_000_000)).toBe('2M')
  })

  it('fractional millions keep one decimal', () => {
    expect(formatContextWindow(1_500_000)).toBe('1.5M')
    expect(formatContextWindow(330_000)).toBe('330k')
  })

  it('thousands round to k, small counts stay bare', () => {
    expect(formatContextWindow(256_000)).toBe('256k')
    expect(formatContextWindow(128_000)).toBe('128k')
    expect(formatContextWindow(200)).toBe('200')
  })

  it('unknown or zero reads as a question mark', () => {
    expect(formatContextWindow(0)).toBe('?')
    expect(formatContextWindow(undefined)).toBe('?')
    expect(formatContextWindow(NaN)).toBe('?')
  })
})
