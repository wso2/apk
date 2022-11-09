package database

const (
	QueryGetApplicationByUUID string = " SELECT " +
		"   APP.UUID," +
		"   APP.NAME," +
		"   APP.SUBSCRIBER_ID," +
		"   APP.ORGANIZATION ORGANIZATION," +
		"   SUB.USER_ID " +
		" FROM " +
		"   AM_SUBSCRIBER SUB," +
		"   AM_APPLICATION APP " +
		" WHERE " +
		"   APP.UUID = $1 " +
		"   AND APP.SUBSCRIBER_ID = SUB.SUBSCRIBER_ID"

	QueryGetAllSubscriptionsForApplication string = "select " +
		"	SUB.uuid as UUID, " +
		"	API.uuid as API_UUID, " +
		"	API.api_version as API_VERSION, " +
		"	SUB.sub_status as SUB_STATUS " +
		"	SUB.organization as ORGANIZATION" +
		"	SUB.created_by as CREATED_BY" +
		"from " +
		"	am_application APP, am_subscription SUB, am_api API " +
		"where 1 = 1 " +
		"	AND APP.application_id = SUB.application_id " +
		"	AND SUB.api_id = API.api_id " +
		"	AND APP.uuid = $1"
)
