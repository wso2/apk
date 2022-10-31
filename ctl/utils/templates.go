package utils

var configMapTemplate = `apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Name}}
  {{if ne .Namespace "" -}}
  namespace: {{.Namespace}}
  {{- end}}
data:
	{{if ne .File "" -}}
    swagger.yaml: |
		{{.SwaggerContent}}
	{{- else -}}
	swagger: |
		{{.DefaultSwagger}}
	{{- end -}}`

const DefaultSwaggerFile = `openapi: 3.0.1
info:
  title: Default
  version: '1.0.0'
servers:
  - url: /
security:
  - default: []
paths:
  /*:
    get:
      responses:
        '200':
          description: OK
      security:
        - default: []
      x-auth-type: Application & Application User
      x-throttling-tier: Unlimited
      x-wso2-application-security:
        security-types:
          - oauth2
        optional: false
    put:
      responses:
        '200':
          description: OK
      security:
        - default: []
      x-auth-type: Application & Application User
      x-throttling-tier: Unlimited
      x-wso2-application-security:
        security-types:
          - oauth2
        optional: false
    post:
      responses:
        '200':
          description: OK
      security:
        - default: []
      x-auth-type: Application & Application User
      x-throttling-tier: Unlimited
      x-wso2-application-security:
        security-types:
          - oauth2
        optional: false
    delete:
      responses:
        '200':
          description: OK
      security:
        - default: []
      x-auth-type: Application & Application User
      x-throttling-tier: Unlimited
      x-wso2-application-security:
        security-types:
          - oauth2
        optional: false
    patch:
      responses:
        '200':
          description: OK
      security:
        - default: []
      x-auth-type: Application & Application User
      x-throttling-tier: Unlimited
      x-wso2-application-security:
        security-types:
          - oauth2
        optional: false
components:
  securitySchemes:
    default:
      type: oauth2
      flows:
        implicit:
          authorizationUrl: 'https://test.com'
          scopes: {}
x-wso2-auth-header: Authorization
x-wso2-cors:
  corsConfigurationEnabled: false
  accessControlAllowOrigins:
    - '*'
  accessControlAllowCredentials: false
  accessControlAllowHeaders:
    - authorization
    - Access-Control-Allow-Origin
    - Content-Type
    - SOAPAction
    - apikey
    - Internal-Key
  accessControlAllowMethods:
    - GET
    - PUT
    - POST
    - DELETE
    - PATCH
    - OPTIONS
x-wso2-production-endpoints:
  urls:
    - 'https://run.mocky.io/v3/ffd5ada6-fca6-4c63-ab74-771331d5c913'
  type: http
x-wso2-sandbox-endpoints:
  urls:
    - 'https://run.mocky.io/v3/ffd5ada6-fca6-4c63-ab74-771331d5c913'
  type: http
x-wso2-basePath: /north/1
x-wso2-transports:
  - http
  - https
x-wso2-response-cache:
  enabled: false
  cacheTimeoutInSeconds: 300`
