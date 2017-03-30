package refex

// characters that have special regex meaning that we need to process differently
var specialChars = []string{`\`, `(`, `)`}

var patternPrefix = `(`
var patternSuffix = `)`

// regex used to replace arguments
const wildcard = `(.*)`
