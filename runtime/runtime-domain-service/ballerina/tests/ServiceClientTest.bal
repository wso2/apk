import ballerina/test;

@test:Config {dataProvider: getServicesDataProvider}
public function testGetServices(string? query, string sortBy, string sortOrder, int 'limit, int offset, anydata expected) {
    ServiceClient serviceClient = new;
    test:assertEquals(serviceClient.getServices(query, sortBy, sortOrder, 'limit, offset,organiztion1), expected);
}

public function getServicesDataProvider() returns map<[string|(), string, string, int, int, ServiceList|BadRequestError|InternalServerErrorError]> {
    BadRequestError badRequest = {body: {code: 90912, message: "Invalid Sort By/Sort Order Value "}};
    BadRequestError badRequest1 = {body: {code: 90912, message: "Invalid KeyWord namespace1"}};

    map<[string|(), string, string, int, int, ServiceList|BadRequestError|InternalServerErrorError]> data = {
        "1": [
            (),
            SORT_BY_SERVICE_NAME,
            SORT_ORDER_ASC,
            10,
            0,
            {
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
            }
        ],
        "2": [
            (),
            SORT_BY_SERVICE_NAME,
            SORT_ORDER_DESC,
            10,
            0,
            {
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
            }
        ],
        "3": [
            (),
            SORT_BY_SERVICE_CREATED_TIME,
            SORT_ORDER_ASC,
            10,
            0,
            {
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
            }
        ],
        "4": [
            (),
            SORT_BY_SERVICE_CREATED_TIME,
            SORT_ORDER_DESC,
            10,
            0,
            {
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
            }
        ],
        "5": [
            (),
            SORT_BY_SERVICE_NAME,
            SORT_ORDER_DESC,
            3,
            0,
            {
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
            }
        ],
        "6": [
            (),
            SORT_BY_SERVICE_NAME,
            SORT_ORDER_DESC,
            3,
            6,
            {
                "list": [],
                "pagination": {
                    "offset": 6,
                    "limit": 3,
                    "total": 6
                }
            }
        ],
        "7": [(), "invalid sort", SORT_ORDER_DESC, 10, 0, badRequest],
        "8": [(), SORT_BY_SERVICE_CREATED_TIME, "invlid order", 10, 0, badRequest],
        "9": [
            "name:httpbin",
            SORT_BY_SERVICE_NAME,
            SORT_ORDER_DESC,
            10,
            0,
            {
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
            }
        ],
        "10": [
            "abcde",
            SORT_BY_SERVICE_NAME,
            SORT_ORDER_ASC,
            10,
            0,
            {
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
            }
        ],
        "11": [
            "namespace:apk-platform",
            SORT_BY_SERVICE_NAME,
            SORT_ORDER_DESC,
            10,
            0,
            {
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
            }
        ],
        "12": ["namespace1:apk-pl", SORT_BY_SERVICE_NAME, SORT_ORDER_DESC, 10, 0, badRequest1]
    };
    return data;
}

@test:Config {dataProvider: serviceByIdDataProvider}
public function testGetServiceByID(string id, anydata expected) {
    ServiceClient serviceClient = new;
    test:assertEquals(serviceClient.getServiceById(id,organiztion1), expected);
}

public function serviceByIdDataProvider() returns map<[string, Service|BadRequestError|NotFoundError|InternalServerErrorError]> {
    NotFoundError notfound = {body: {code: 90914, message: "Service abcd-efght not found"}};
    map<[string, Service|BadRequestError|NotFoundError|InternalServerErrorError]> data = {
        "1": [
            "275b00d1-712c-4df2-b65a-9b14678abe5b",
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
        "2": ["abcd-efght", notfound]
    };
    return data;
}

@test:Config {dataProvider: serviceUsageDataProvider}
public function testServiceUsageByID(string serviceId, anydata expected) {
    ServiceClient serviceClient = new;
    test:assertEquals(serviceClient.getServiceUsageByServiceId(serviceId,organiztion1), expected);
}

function serviceUsageDataProvider() returns map<[string, APIList|BadRequestError|NotFoundError|InternalServerErrorError]> {
    NotFoundError notfound = {body: {code: 90914, message: "Service 275b00d1-722c-4df2-b65a-9b14678abe6b not found"}};
    map<[string, APIList|BadRequestError|NotFoundError|InternalServerErrorError]> data = {
        "1": [
            "275b00d1-722c-4df2-b65a-9b14677abe4b",
            {
                "count": 2,
                "list": [
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    }
                ],
                "pagination": {
                    "total": 2
                }
            }
        ],
        "2": ["275b00d1-722c-4df2-b65a-9b14677abe5b", {"count": 0, "list": [], "pagination": {"total": 0}}],
        "3": ["275b00d1-722c-4df2-b65a-9b14678abe6b", notfound]

    };
    return data;
}

