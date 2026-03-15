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

{{- define "golanglab5.configChecksum" -}}
{{- required "set config.configK8sContent via --set-file config.configK8sContent=config/config.k8s.yaml" .Values.config.configK8sContent | sha256sum -}}
{{- end }}

{{- define "golanglab5.otelCollectorConfig" -}}
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318

processors:
  memory_limiter:
    check_interval: 1s
    limit_mib: 256
  k8sattributes:
    auth_type: serviceAccount
    extract:
      metadata:
        - k8s.namespace.name
        - k8s.pod.name
        - k8s.node.name
        - k8s.deployment.name
    pod_association:
      - sources:
          - from: connection
  batch:
    timeout: 1s
    send_batch_size: 1024

exporters:
  otlp/jaeger:
    endpoint: jaeger:4317
    tls:
      insecure: true
  prometheus:
    endpoint: 0.0.0.0:8889
    namespace: app

extensions:
  health_check:
    endpoint: 0.0.0.0:13133

service:
  extensions: [health_check]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, k8sattributes, batch]
      exporters: [otlp/jaeger]
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, k8sattributes, batch]
      exporters: [prometheus]
{{- end }}

{{- define "golanglab5.otelCollectorConfigChecksum" -}}
{{- include "golanglab5.otelCollectorConfig" . | sha256sum -}}
{{- end }}
