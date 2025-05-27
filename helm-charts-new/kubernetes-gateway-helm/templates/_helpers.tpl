{{/*
Expand the name of the chart.
*/}}
{{- define "kubernetes-gateway-helm.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kubernetes-gateway-helm.fullname" -}}
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
{{- define "kubernetes-gateway-helm.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kubernetes-gateway-helm.labels" -}}
helm.sh/chart: {{ include "kubernetes-gateway-helm.chart" . }}
{{ include "kubernetes-gateway-helm.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kubernetes-gateway-helm.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kubernetes-gateway-helm.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "kubernetes-gateway-helm.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "kubernetes-gateway-helm.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kubernetes-gateway-helm.pod.selectorLabels" -}}
app.kubernetes.io/app: {{ .app }}
app.kubernetes.io/release: {{ .root.Release.Name }}
{{- end }}

{{- define "kubernetes-gateway-helm.deployment.affinity" -}}
{{- $value := typeIs "string" .value | ternary .value (.value | toYaml) }}
{{- if  (not .value) -}}
podAntiAffinity:
  preferredDuringSchedulingIgnoredDuringExecution:
  - podAffinityTerm:
      labelSelector:
        matchExpressions:
        - key: app.kubernetes.io/app
          operator: In
          values:
          - {{ .app }}
      topologyKey: "topology.kubernetes.io/zone"
    weight: 100
{{- else if contains "{{" (toJson .value) }}
    {{- tpl $value .context }}
{{- else }}
    {{- $value }}
{{- end }}
{{- end -}}

{{- define "kubernetes-gateway-helm.deployment.nodeSelector" -}}
{{- $value := typeIs "string" .value | ternary .value (.value | toYaml) }}
{{- if contains "{{" (toJson .value) }}
    {{- tpl $value .context }}
{{- else }}
    {{- $value }}
{{- end }}
{{- end -}}

{{- define "kubernetes-gateway-helm.deployment.readinessProbe.http" -}}
readinessProbe:
  httpGet:
    path: {{ .readinessProbe.path }}
    port: {{ .readinessProbe.port }}
  initialDelaySeconds: {{ .readinessProbe.initialDelaySeconds }}
  periodSeconds: {{ .readinessProbe.periodSeconds }}
  failureThreshold: {{ .readinessProbe.failureThreshold }}
{{- end }}

{{- define "kubernetes-gateway-helm.deployment.livenessProbe.http" -}}
livenessProbe:
  httpGet:
    path: {{ .livenessProbe.path }}
    port: {{ .livenessProbe.port }}
  initialDelaySeconds: {{ .livenessProbe.initialDelaySeconds }}
  periodSeconds: {{ .livenessProbe.periodSeconds }}
  failureThreshold: {{ .livenessProbe.failureThreshold }}
{{- end }}

{{- define "kubernetes-gateway-helm.deployment.resources" -}}
resources:
  requests:
    memory: {{ .requests.memory }}
    cpu: {{ .requests.cpu }}
  limits:
    memory: {{ .limits.memory }}
    cpu: {{ .limits.cpu }}
{{- end }}


{{/*
Common prefix prepended to Kubernetes resources of this chart
*/}}
{{- define "kubernetes-gateway-helm.resource.prefix" -}}
{{- printf "%s-wso2-kgw" .Release.Name -}}
{{- end -}}




{{- define "kubernetes-gateway-helm.deployment.env" -}}
env:
{{- if . -}}
{{- range $key, $val := . }}
  - name: {{ $key }}
    value: {{ quote  $val }}
{{- end }}
{{- end -}}
{{- end -}}

{{- define "commaJoinedQuotedList" -}}
{{- $list := list }}
{{- range .}}
{{- $list = append $list (printf "\"%s\"" .) }}
{{- end }}
{{- join ", " $list }}
{{- end }}

{{- define "generateVhosts" -}}
{{- if . -}}
{{- $vhosts := . -}}
{{- range $vhost := $vhosts }}
{{- printf "[[vhosts]]\n" -}}
{{- printf "  name = \"%s\"\n" $vhost.name -}}
{{- print "  hosts = [" -}}
{{- $len := len $vhost.hosts -}}
{{- range  $i, $host := $vhost.hosts }}
{{- printf "\"%s\"" $host -}}
{{- if lt $i (sub $len 1) }},{{ end }}
{{- end }}
{{- print "]\n" -}}
{{- printf "  type = \"%s\"\n" $vhost.type -}}
{{- end }}
{{- end }}
{{- end }}

{{- define "createYamlList" -}}
{{- if . -}}
{{ range $val := . }}
- "{{ $val }}"
{{- end -}}
{{- else -}}
- "*.gw.wso2.com"
- "*.sandbox.gw.wso2.com"
- "prod.gw.wso2.com"
{{- end -}}
{{- end -}}

{{- define "kgw.javaOptions" -}}
  {{- if .Values.wso2.kgw.dp.gatewayRuntime.deployment.enforcer.configs.javaOpts }}
    {{- .Values.wso2.kgw.dp.gatewayRuntime.deployment.enforcer.configs.javaOpts }}
    {{- if and .Values.wso2.kgw.metrics .Values.wso2.kgw.metrics.enabled -}}
      {{- " " }}-Dkgw.jmx.metrics.enabled=true -javaagent:/home/wso2/lib/jmx_prometheus_javaagent-0.20.0.jar=18006:/tmp/metrics/prometheus-jmx-config-enforcer.yml
    {{- end }}
  {{- else -}}
    -Dhttpclient.hostnameVerifier=AllowAll -Xms512m -Xmx512m -XX:MaxRAMFraction=2
    {{- if and .Values.wso2.kgw.metrics .Values.wso2.kgw.metrics.enabled }}
      {{- " " }}-Dkgw.jmx.metrics.enabled=true -javaagent:/home/wso2/lib/jmx_prometheus_javaagent-0.20.0.jar=18006:/tmp/metrics/prometheus-jmx-config-enforcer.yml
    {{- end }}
  {{- end }}
{{- end }}
