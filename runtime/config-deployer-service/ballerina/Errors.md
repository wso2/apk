# Errors in Runtime Domain Service

These are the runtime domain service errors and their respective error codes.

| Error Code | Status Code | Error Message |
|---|---|---|
| 909000  | --- | Common code for other error |
| 909001  | 404 | **ID** not found |
| 909002  | 404 | Context/Name doesn't exist |
| 909003  | 404 | apiId not found in request |
| 909004  | 406 | Invalid property id in Request |
| 909005  | 404 | type field unavailable |
| 909006  | 406 | Unsupported API type |
| 909007  | 406 | Multiple fields of url, file, inlineAPIDefinition given |
| 909008  | 406 | Atleast one of the field required |
| 909009  | 406 | Additional properties not provided |
| 909010  | 406 | Invalid operation policy name |
| 909011  | 409 | API Name **apiName** already exist |
| 909012  | 409 | API Base Path **apiBasePath** already exist |
| 909013  | 406 | Sandbox endpoint not specified |
| 909014  | 406 | Production endpoint not specified |
| 909015  | 406 | API Base Path **apiBasePath** invalid |
| 909016  | 406 | API name **apiName** invalid |
| 909017  | 406 | Invalid API request |
| 909018  | 500 | Error while generating token |
| 909019  | 406 | Invalid keyword **keyWord** |
| 909020  | 406 | Invalid Sort By/Sort Order value |
| 909021  | 406 | Atleast one operation need to specified |
| 909022  | 500 | Internal server error |
| 909023  | 500 | Internal error occured while retrieving definition |
| 909024  | 406 | Invalid parameters provided for policy **policyName** |
| 909025  | 406 | Presence of both resource level and API level operation policies is not allowed |
| 909026  | 406 | Presence of both resource level and API level rate limits is not allowed |
| 909027  | 500 | Error while retrieving API |
| 909028  | 500 | Internal error occured while deploying API |
| 909029  | 500 | Error while retrieving Mediation policy |
| 909030  | 400 | Certificate is expired |
| 909031  | 500 | Error while adding certificate |
| 909032  | 500 | Host/Certificte is empty in payload |
| 909033  | 500 | Error while retrieving endpoint certificate request |
| 909034  | 404 | Certificate **certificateId** not found |
| 909035  | 500 | Error while deleting endpoint certificate |
| 909036  | 500 | Error while getting endpoint certificate content |
| 909037  | 500 | Error while getting endpoint certificate by id |
| 909038  | 500 | Error while updating endpoint certificate |
| 909039  | 406 | Invalid value for offset |
| 909040  | 500 | Internal Error Occured while converting json to yaml |
| 909041  | 406 | Accept header should be application/json or application/yaml |
| 909042  | 400 | Unsupported API type |
| 909043  | 500 | Error occured while generating openapi definition |
| 909044  | 406 | Retrieve definition from Url result |
| 909045  | 404 | New Version/APIID does not exist |
| 909046  | 409 | New version **newVersion** already exist |
| 909047  | 404 | **serviceId** service does not exist |
| 909048  | 406 | API Type change not supported from update |
| 909049  | 406 | Context change not supported from update |
| 909050  | 406 | Version change not supported from update |
