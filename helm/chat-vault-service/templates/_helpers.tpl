{{- define "chatvaultservice.name" -}}
{{- default "chat-vault-service" .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Helm required labels */}}
{{- define "chatvaultservice.labels" -}}
heritage: {{ .Release.Service }}
release: {{ .Release.Name }}
chart: {{ .Chart.Name }}
app: "{{ template "chatvaultservice.name" . }}"
{{- end -}}

{{/* matchLabels */}}
{{- define "chatvaultservice.matchLabels" -}}
release: {{ .Release.Name }}
app: "{{ template "chatvaultservice.name" . }}"
{{- end -}}