{{/*
# -------------------------------------------------------------------------------------
#
# Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com). All Rights Reserved.
#
# This software is the property of WSO2 LLC. and its suppliers, if any.
# Dissemination of any information or reproduction of any material contained 
# herein is strictly forbidden, unless permitted by WSO2 in accordance with the 
# WSO2 Commercial License available at https://wso2.com/licenses/eula/3.2
#
# --------------------------------------------------------------------------------------s
*/}}

{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "apim-helm-cp.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "apim-helm-cp.fullname" -}}
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
{{- define "apim-helm-cp.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "apim-helm-cp.labels" -}}
app.kubernetes.io/name: {{ include "apim-helm-cp.name" . }}
helm.sh/chart: {{ include "apim-helm-cp.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Common prefix prepended to Kubernetes resources of this chart
*/}}
{{- define "apim-helm-cp.resource.prefix" -}}
{{- if eq .Values.kubernetes.resourceSuffix "${RESOURCE_SUFFIX}" }}
{{- "wso2am-cp" }}
{{- else }}
{{- "wso2am-cp-" }}{{ .Values.kubernetes.resourceSuffix }}
{{- end -}}
{{- end -}}