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

{{- define "apk-helm.cnpg.appUserPassword" -}}
{{ .Values.wso2.apk.cp.cnpg.appUserPassword | b64enc}}
{{- end }}

{{- define "apk-helm.cnpg.superUserPassword" -}}
{{ .Values.wso2.apk.cp.cnpg.superUserPassword | b64enc}}
{{- end }}


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