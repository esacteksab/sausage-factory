{{range .Entries}}
# v{{ .Semver }}

### What's Changed

<!-- keep-sorted start by_regex=\w+ prefix_order=feat,fix,chore,style,docs,build -->
{{range .Changes -}}{{$note := splitList "\n" .Note}}
- **({{ .ConventionalCommit.Category}}):** {{ .ConventionalCommit.Description -}}
{{end -}}
{{- end}}
<!-- keep-sorted end -->