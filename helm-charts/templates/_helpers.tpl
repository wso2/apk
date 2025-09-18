{{/*
 Copyright (c) 2022, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.

 WSO2 LLC. licenses this file to you under the Apache License,
 Version 2.0 (the "License"); you may not use this file except
 in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied. See the License for the
 specific language governing permissions and limitations
 under the License.
*/}}

{{/*
Selector labels
*/}}
{{- define "apk-helm.pod.selectorLabels" -}}
app.kubernetes.io/app: {{ .app }}
app.kubernetes.io/release: {{ .root.Release.Name }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "apk-helm.labels" -}}
{{- if .Values.labels }}
{{- if .Values.labels.common }}
{{- range $key, $val := .Values.labels.common -}}
{{ $key }}: {{ $val }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

{{- define "apk-helm.deployment.affinity" -}}
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

{{- define "apk-helm.deployment.nodeSelector" -}}
{{- $value := typeIs "string" .value | ternary .value (.value | toYaml) }}
{{- if contains "{{" (toJson .value) }}
    {{- tpl $value .context }}
{{- else }}
    {{- $value }}
{{- end }}
{{- end -}}

{{- define "apk-helm.deployment.readinessProbe.http" -}}
readinessProbe:
  httpGet:
    path: {{ .readinessProbe.path }}
    port: {{ .readinessProbe.port }}
  initialDelaySeconds: {{ .readinessProbe.initialDelaySeconds }}
  periodSeconds: {{ .readinessProbe.periodSeconds }}
  failureThreshold: {{ .readinessProbe.failureThreshold }}
{{- end }}

{{- define "apk-helm.deployment.livenessProbe.http" -}}
livenessProbe:
  httpGet:
    path: {{ .livenessProbe.path }}
    port: {{ .livenessProbe.port }}
  initialDelaySeconds: {{ .livenessProbe.initialDelaySeconds }}
  periodSeconds: {{ .livenessProbe.periodSeconds }}
  failureThreshold: {{ .livenessProbe.failureThreshold }}
{{- end }}

{{- define "apk-helm.deployment.resources" -}}
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
{{- define "apk-helm.resource.prefix" -}}
{{- printf "%s-wso2-apk" .Release.Name -}}
{{- end -}}




{{- define "apk-helm.deployment.env" -}}
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

{{- define "apk.javaOptions" -}}
  {{- if .Values.wso2.apk.dp.gatewayRuntime.deployment.enforcer.configs.javaOpts }}
    {{- .Values.wso2.apk.dp.gatewayRuntime.deployment.enforcer.configs.javaOpts }}
    {{- if and .Values.wso2.apk.metrics .Values.wso2.apk.metrics.enabled -}}
      {{- " " }}-Dapk.jmx.metrics.enabled=true -javaagent:/home/wso2/lib/jmx_prometheus_javaagent-0.20.0.jar=18006:/tmp/metrics/prometheus-jmx-config-enforcer.yml
    {{- end }}
  {{- else -}}
    -Dhttpclient.hostnameVerifier=AllowAll -Xms512m -Xmx512m -XX:MaxRAMFraction=2
    {{- if and .Values.wso2.apk.metrics .Values.wso2.apk.metrics.enabled }}
      {{- " " }}-Dapk.jmx.metrics.enabled=true -javaagent:/home/wso2/lib/jmx_prometheus_javaagent-0.20.0.jar=18006:/tmp/metrics/prometheus-jmx-config-enforcer.yml
    {{- end }}
  {{- end }}
{{- end }}
