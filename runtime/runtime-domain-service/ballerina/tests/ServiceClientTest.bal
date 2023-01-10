import ballerina/test;
@test:Config{dataProvider:getServicesDataProvider}
public function testGetServices(string? query, string sortBy, string sortOrder, int 'limit, int offset,anydata expected){
    ServiceClient serviceClient = new;
    test:assertEquals(serviceClient.getServices(query,sortBy,sortOrder,'limit,offset),expected);
}

public function getServicesDataProvider() returns map<[string|(),string,string,int,int,ServiceList|BadRequestError|InternalServerErrorError]>{
            BadRequestError badRequest = {body: {code: 90912, message: "Invalid Sort By/Sort Order Value "}};
                        BadRequestError badRequest1 = {body: {code: 90912, message: "Invalid KeyWord namespace1"}};

map<[string|(),string,string,int,int,ServiceList|BadRequestError|InternalServerErrorError]> data = {
"1":[(),SORT_BY_SERVICE_NAME,SORT_ORDER_ASC,10,0,{
	"list": [
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe5b",
			"name": "abcde",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T07:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14678abe5b",
			"name": "abcde1",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T10:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe4b",
			"name": "backend",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T08:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14678abe4b",
			"name": "backend-1",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T14:30:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe4be",
			"name": "backend-15",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T13:25:09Z"
		},
		{
			"id": "275b00d1-712c-4df2-b65a-9b14678abe5b",
			"name": "httpbin",
			"namespace": "apk-platform",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T12:25:09Z"
		}
	],
	"pagination": {
		"offset": 0,
		"limit": 10,
		"total": 6
	}
}],
"2":[(),SORT_BY_SERVICE_NAME,SORT_ORDER_DESC,10,0,{
	"list": [
		{
			"id": "275b00d1-712c-4df2-b65a-9b14678abe5b",
			"name": "httpbin",
			"namespace": "apk-platform",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T12:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe4be",
			"name": "backend-15",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T13:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14678abe4b",
			"name": "backend-1",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T14:30:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe4b",
			"name": "backend",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T08:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14678abe5b",
			"name": "abcde1",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T10:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe5b",
			"name": "abcde",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T07:25:09Z"
		}
	],
	"pagination": {
		"offset": 0,
		"limit": 10,
		"total": 6
	}
}],
"3":[(),SORT_BY_SERVICE_CREATED_TIME,SORT_ORDER_ASC,10,0,{
	"list": [
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe5b",
			"name": "abcde",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T07:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe4b",
			"name": "backend",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T08:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14678abe5b",
			"name": "abcde1",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T10:25:09Z"
		},
		{
			"id": "275b00d1-712c-4df2-b65a-9b14678abe5b",
			"name": "httpbin",
			"namespace": "apk-platform",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T12:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe4be",
			"name": "backend-15",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T13:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14678abe4b",
			"name": "backend-1",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T14:30:09Z"
		}
	],
	"pagination": {
		"offset": 0,
		"limit": 10,
		"total": 6
	}
}],
"4":[(),SORT_BY_SERVICE_CREATED_TIME,SORT_ORDER_DESC,10,0,{
	"list": [
		{
			"id": "275b00d1-722c-4df2-b65a-9b14678abe4b",
			"name": "backend-1",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T14:30:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe4be",
			"name": "backend-15",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T13:25:09Z"
		},
		{
			"id": "275b00d1-712c-4df2-b65a-9b14678abe5b",
			"name": "httpbin",
			"namespace": "apk-platform",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T12:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14678abe5b",
			"name": "abcde1",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T10:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe4b",
			"name": "backend",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T08:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe5b",
			"name": "abcde",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T07:25:09Z"
		}
	],
	"pagination": {
		"offset": 0,
		"limit": 10,
		"total": 6
	}
}],
"5":[(),SORT_BY_SERVICE_NAME,SORT_ORDER_DESC,3,0,{
	"list": [
		{
			"id": "275b00d1-712c-4df2-b65a-9b14678abe5b",
			"name": "httpbin",
			"namespace": "apk-platform",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T12:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe4be",
			"name": "backend-15",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T13:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14678abe4b",
			"name": "backend-1",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T14:30:09Z"
		}
	],
	"pagination": {
		"offset": 0,
		"limit": 3,
		"total": 6
	}
}],
"6":[(),SORT_BY_SERVICE_NAME,SORT_ORDER_DESC,3,6,{
	"list": [],
	"pagination": {
		"offset": 6,
		"limit": 3,
		"total": 6
	}
}],
"7":[(),"invalid sort",SORT_ORDER_DESC,10,0,badRequest],
"8":[(),SORT_BY_SERVICE_CREATED_TIME,"invlid order",10,0,badRequest],
"9":["name:httpbin",SORT_BY_SERVICE_NAME,SORT_ORDER_DESC,10,0,{
	"list": [
		{
			"id": "275b00d1-712c-4df2-b65a-9b14678abe5b",
			"name": "httpbin",
			"namespace": "apk-platform",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T12:25:09Z"
		}
	],
	"pagination": {
		"offset": 0,
		"limit": 10,
		"total": 1
	}
}],
"10":["abcde",SORT_BY_SERVICE_NAME,SORT_ORDER_ASC,10,0,{
	"list": [
		{
			"id": "275b00d1-722c-4df2-b65a-9b14677abe5b",
			"name": "abcde",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T07:25:09Z"
		},
		{
			"id": "275b00d1-722c-4df2-b65a-9b14678abe5b",
			"name": "abcde1",
			"namespace": "apk",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T10:25:09Z"
		}
	],
	"pagination": {
		"offset": 0,
		"limit": 10,
		"total": 2
	}
}],
"11":["namespace:apk-platform",SORT_BY_SERVICE_NAME,SORT_ORDER_DESC,10,0,{
	"list": [
		{
			"id": "275b00d1-712c-4df2-b65a-9b14678abe5b",
			"name": "httpbin",
			"namespace": "apk-platform",
			"type": "ClusterIP",
			"portmapping": [
				{
					"name": "http",
					"protocol": "TCP",
					"targetport": 80,
					"port": 80
				}
			],
			"createdTime": "2022-12-13T12:25:09Z"
		}
	],
	"pagination": {
		"offset": 0,
		"limit": 10,
		"total": 1
	}
}],
"12":["namespace1:apk-pl",SORT_BY_SERVICE_NAME,SORT_ORDER_DESC,10,0,badRequest1]
};
return data;
}