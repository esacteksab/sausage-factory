## Tools

<!-- mdformat-toc start --slug=github --no-anchors --maxlevel=6 --minlevel=1 -->

- [Tools](#tools)
  - [Looking at (in no particular order)](#looking-at-in-no-particular-order)
  - [Adjacent](#adjacent)
    - [Conventional Commits Tooling](#conventional-commits-tooling)
    - [Release Tools](#release-tools)
    - [Changelog](#changelog)
    - [Inspiration](#inspiration)

<!-- mdformat-toc end -->

### Looking at (in no particular order)

- [gh-extenions-precompile](https://github.com/cli/gh-extension-precompile) GitHub action from the `cli/cli` team themselves.
  - This is referenced in their Creating GitHub CLI extensions [Tips for writing precompiled GitHub CLI extensions](https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions#tips-for-writing-precompiled-github-cli-extensions).
- Hugo's [hugoreleaser](https://github.com/gohugoio/hugoreleaser).
  - Their reason for [Why another Go release tool](https://github.com/gohugoio/hugoreleaser?tab=readme-ov-file#why-another-go-release-tool) is compelling.
- I'm not sure I can talk about releasing a Go binary in 2025 and not mention [GoReleaser](https://goreleaser.com/).


## CLI "Things"

- https://devcenter.heroku.com/articles/cli-style-guide
- https://medium.com/jit-team/guidelines-for-creating-your-own-cli-tool-c95d4af62919
- https://evilmartians.com/chronicles/cli-ux-best-practices-3-patterns-for-improving-progress-displays
- https://www.atlassian.com/blog/it-teams/10-design-principles-for-delightful-clis
- (a bit of irony) https://github.com/cli-guidelines/cli-guidelines
- https://smallstep.com/blog/the-poetics-of-cli-command-names/
- https://kddnewton.com/2016/10/11/exploring-cli-best-practices.html
- https://blog.prototypr.io/ux-for-command-line-tools-4630eb0b3c9b
- https://www.thoughtworks.com/en-us/insights/blog/engineering-effectiveness/elevate-developer-experiences-cli-design-guidelines (most of the links above came from this article)


### Adjacent

#### Conventional Commits Tooling

Conventional Commits [tooling for conventional commits](https://www.conventionalcommits.org/en/about/#tooling-for-conventional-commits).

More specifically (in no particular order)

- [leodido/go-conventionalcommits](https://github.com/leodido/go-conventionalcommits).
- [joselitofilho/go-conventional-commits](https://github.com/joselitofilho/go-conventional-commits).
- [Conventional Commits Linter](https://gitlab.com/DeveloperC/conventional_commits_linter) looks to be written in Rust. Actively maintained but with a ["Pre 1.0.0 breaking changes _may_ be introduced without increasing the major version"](https://gitlab.com/DeveloperC/conventional_commits_linter#pre-100-breaking-changes-may-be-introduced-without-increasing-the-major-version) warning gives me pause in the short-term.
- From the same folks as the previously mentioned tool, a tool called [Conventional Commits Next Version](https://gitlab.com/DeveloperC/conventional_commits_next_version) which uses Convetional Commits to calculate the semantic version for a release.

#### Release Tools

- [release-please from Google](https://github.com/googleapis/release-please) Not a release tool like the tools above, but could be useful in orchestrating a release.

#### Changelog

- [hashicorp/go-changelog](https://github.com/hashicorp/go-changelog). Mixed feelings about HashiCorp and the future of their projects with the pending [IBM acquisition](https://www.hashicorp.com/en/blog/hashicorp-joins-ibm).
- Seems GoReleaser has some capabilities to [customize a changelog](https://goreleaser.com/customization/changelog/). I haven't dug in enough to have an opinion.
  - They also have a Go tool called [chlog](https://github.com/goreleaser/chglog) available on their GitHub. It seems to be actively developed, but the [Status](https://github.com/goreleaser/chglog?tab=readme-ov-file#status) in their [README.md](https://github.com/goreleaser/chglog/blob/main/README.md) says _alpha_.

#### Inspiration

- [paultyng/changelog-gen](https://github.com/paultyng/changelog-gen/tree/master).
- [twisted/towncrier](https://github.com/twisted/towncrier).
- [Turbogit (tug)](https://b4nst.github.io/turbogit/).
