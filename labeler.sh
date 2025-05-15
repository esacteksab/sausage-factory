#!/bin/bash

# Set PR title with backticks escaped or replaced
PRTITLE='chore(deps): bump golang from $(00eccd4) to $(915f66a)'
# Or use escaped backticks
# PRTITLE="chore(deps): bump golang from \`00eccd4\` to \`915f66a\`"

REGEX="\(([^)]+)\)"
if [[ ${PRTITLE} =~ ${REGEX} ]]; then
    SCOPE="${BASH_REMATCH[1]}"
    TYPE="$(echo ${PRTITLE} | cut -d'(' -f1)"
else
    TYPE="$(echo "${PRTITLE}" | cut -d':' -f1)"
fi
if [[ ${SCOPE} == "deps" ]]; then
    echo "gh pr edit 127 --add-label "type: ${TYPE},dependencies""
else
    echo "gh pr edit 127 --add-label "type: ${TYPE}""
fi
