{{/*
Expand the name of the chart.
*/}}
{{- define "golanglab5.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "golanglab5.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "golanglab5.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "golanglab5.labels" -}}
helm.sh/chart: {{ include "golanglab5.chart" . }}
{{ include "golanglab5.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "golanglab5.selectorLabels" -}}
app.kubernetes.io/name: {{ include "golanglab5.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "golanglab5.r2SecretName" -}}
{{- default (printf "%s-r2" (include "golanglab5.fullname" .)) .Values.r2.secretName -}}
{{- end }}

{{- define "golanglab5.r2Env" -}}
- name: R2_Endpoint
  valueFrom:
    secretKeyRef:
      name: {{ include "golanglab5.r2SecretName" . }}
      key: endpoint
- name: R2_ACCESS_KEY_ID
  valueFrom:
    secretKeyRef:
      name: {{ include "golanglab5.r2SecretName" . }}
      key: accessKeyId
- name: R2_SECRET_ACCESS_KEY
  valueFrom:
    secretKeyRef:
      name: {{ include "golanglab5.r2SecretName" . }}
      key: secretAccessKey
{{- end }}
