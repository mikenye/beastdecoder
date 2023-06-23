package common

// Number of latitude zones between the equator and a pole.
// In Mode S, NZ is defined to be 15. See: https://mode-s.org/decode/content/ads-b/3-airborne-position.html#core-functions-and-parameters
const Nz = float64(15)

// CPR Format
type CprFormat uint8

const CprFormatEvenFrame = CprFormat(0)
const CprFormatOddFrame = CprFormat(1)
