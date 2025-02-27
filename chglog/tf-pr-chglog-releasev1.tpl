{{- with (first .Entries)}}
# v{{ .Semver }}

## What's Changed

{{range .Changes -}}
{{ if eq .ConventionalCommit.Category "feat" -}}

### {{.ConventionalCommit.Category}}
- **({{.ConventionalCommit.Category}}):** :feat: {{.ConventionalCommit.Description}}
{{ else -}}
{{ if eq .ConventionalCommit.Category "fix" -}}
- **({{.ConventionalCommit.Category}}):** :fix: {{.ConventionalCommit.Description }}
{{ else -}}
{{ if eq .ConventionalCommit.Category "chore" -}}
- **({{.ConventionalCommit.Category}}):** :chore: {{.ConventionalCommit.Description }}
{{ else -}}
{{ if eq .ConventionalCommit.Category "style" -}}
- **({{.ConventionalCommit.Category}}):** :style: {{.ConventionalCommit.Description }}
{{ else -}}
{{ if eq .ConventionalCommit.Category "docs" -}}
- **({{.ConventionalCommit.Category}}):** :docs: {{.ConventionalCommit.Description }}
{{ else -}}
{{ if eq .ConventionalCommit.Category "build" -}}
- **({{.ConventionalCommit.Category}}):** :build: {{.ConventionalCommit.Description }}
{{ else -}}{{ end }}{{- end }}{{- end}}{{- end}}{{end}}{{end}}{{end}}{{end}}