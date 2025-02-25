#!/bin/bash
#
# The changelog formatting script for gh-tf-pr

err() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&2
}

gen_entry () {
    local CHGLOG
    CHGLOG=CHANGELOG.md
    local CHGLOGENTRY
    CHGLOGENTRY=TMP.md
    ${GOPATH}/bin/chglog format --template-file tf-pr-chglog.tpl >> $CHGLOGENTRY

    cat $CHGLOG >> $CHGLOGENTRY
    cat $CHGLOGENTRY > $CHGLOG
    rm $CHGLOGENTRY
    rm changelog.yml
    ${GOPATH}/bin/keep-sorted $CHGLOG
    unset CHGLOGENTRY
    unset CHGLOG
}

gen_entry

echo $GOBIN

# vim: tabstop=2 shiftwidth=2 softtabstop=2 expandtab:
