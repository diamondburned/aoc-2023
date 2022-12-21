## Dependencies

- `libqalculate` (for `bin/qalc`)
- `vim` :^)

## Prepare and Run

1. <kbd>:</kbd> `%s/\([a-z]\{4}\)/function("\1()",)/g`
2. <kbd>:</kbd> `%s/function("\(.*\)",):/\1 :=/g`
3. <kbd>:</kbd> `/humn() := .*/d`
4. <kbd>:</kbd> `%s/function("humn()",)/"humn"/g`
5. <kbd>:</kbd> `%s/root() := \(.*\) [+-\*\/] \(.*\)/root() := \1 = \2/`
6. <kbd>G</kbd> <kbd>o</kbd>
7. Add `function("root()",)`
8. Run `./qalc.sh input.qalc`
