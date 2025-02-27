{{- with (first .Entries) -}}
# v{{.Semver}}

## What's changed
{{/* We Have an empty list */}}
{{- $myList := list -}}
{{- range .Changes -}}
{{- $topics := list .ConventionalCommit.Category -}}
{{- range $topic := $topics -}}
{{- $myList = append $myList (first $topics) -}}
{{- end -}}{{- end -}} 
{{/* This is our list of topics */}}
{{ $myList }}###{{0.}}

{{ range .Changes }}
{{- if eq .ConventionalCommit.Category "docs" -}}
{{ $docs := list }}
{{- $items := list .ConventionalCommit.Description -}}
{{ range $item := $items -}}
* **(docs):** {{ $docs = append $docs $item -}}
{{end -}}
{{ first $docs}}
{{end}}{{end}}{{end}}