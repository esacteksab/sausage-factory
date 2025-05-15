# Matches [\w]@?(v[\d])

- @v4

## Does Not Match

- uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 #v4.2.2

# Matches @[[\w+]+$

- @v4
- @master
- @SHA$

## Does Not Match

- uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 #v4.2.2

# Matches @[a-zA-Z].+$

- @v4
- @master
- @SHA$

## Does Not Match

- uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 #v4.2.2

# Matches @[[:alnum:]].\w+

- @d91bce49a530db16a0d74697709e451f7f9e0648
- @11bd71901bbe5b1630ceea73d27597364c9af683
- @master

## Does Not Match

- @v4

## Get Tags

gh api repos/crate-ci/committed/releases/latest

gh api repos/crate-ci/committed/releases/tags/v1.1.7
