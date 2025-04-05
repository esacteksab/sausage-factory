#!/bin/bash
export MESSAGE
export TYPE

MESSAGE="build(deps): update esacteksab/.github requirement to e1f5cdbf94f2feb7be5364f8cabf33550604ff64"

TYPE="$(echo ${MESSAGE} | cut -d':' -f1)"

gh pr edit 7 --add-label "type: ${TYPE}"

echo "Type: ${TYPE}"
