#!/bin/bash
#
# The release script for gh-tf-pr

export VERSION
export CHGLOG
export PATTERN

CHGLOG=CHANGELOG.md
PATTERN="v[\d].[\d].[\d](-\w+)?$"


TAG=$(git describe --tags HEAD 2> /dev/null)

does_match() {
    VERSION=$(grep -P $PATTERN $CHGLOG | head -n 1 |cut -d' ' -f2)
    # shellcheck disable=SC2053
    if [[ $TAG == $VERSION ]]; then
        echo "They Match"
    else
        echo "Something isn't right"
    fi
}

does_match

echo "Tag: $TAG"
echo "Version: $VERSION"
