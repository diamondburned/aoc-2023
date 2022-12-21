## Dependencies

- `libqalculate` (for `bin/qalc`)
- `vim` :^)

## Prepare and Run

1. `:%s/\([a-z]\{4}\)/function("\1()",)/g`
2. Remove all `function()`s off the left-hand side.
3. Delete the `humn() :=` definition.
4. Replace `function("humn()",)` with `"humn"` (variable).
5. Find `root()` and replace the operator with `=`.
6. Append `function("root(x)",)` to end of file.
7. `./qalc.sh input.qalc`
