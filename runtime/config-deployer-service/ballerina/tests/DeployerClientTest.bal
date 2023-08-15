import ballerina/test;
import config_deployer_service.model;

@test:Config {dataProvider: RoutesOrderDataProvider}
public isolated function testCreateHttpRoutesOrder(model:Httproute[] httpRoutes, model:Httproute[] expectedSortedhttpRoutes) returns error? {

    DeployerClient deployerClient = new;
    model:Httproute[] sortedHttpRoutes = check deployerClient.createHttpRoutesOrder(httpRoutes);
    test:assertEquals(sortedHttpRoutes, expectedSortedhttpRoutes, "Sorted HttpRoutes are not equal to expected SortedRoutes");
}

public function RoutesOrderDataProvider() returns map<[model:Httproute[], model:Httproute[]]>|error {
    model:Httproute sortedHttpRoute = {
        apiVersion: "gateway.networking.k8s.io/v1beta1",
        kind: "HTTPRoute",
        metadata: {name: "01ee37b2-57b1-12ee-b8c5-e9b11538a0c9"},
        spec: {
            rules: [
                {
                    matches: [
                        {
                            path: {
                                'type: "RegularExpression",
                                value: "/employee/one/two/three"
                            },
                            method: "GET"
                        }
                    ],
                    filters: [
                        {
                            'type: "URLRewrite",
                            urlRewrite: {
                                path: {
                                    'type: "ReplaceFullPath",
                                    replaceFullPath: "/employee/one/two/three"
                                }
                            }
                        }
                    ],
                    backendRefs: [
                        {
                            group: "dp.wso2.com",
                            kind: "Backend",
                            name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                        }
                    ]
                },
                {
                    matches: [
                        {
                            path: {
                                'type: "RegularExpression",
                                value: "/employee/one/two"
                            },
                            method: "GET"
                        }
                    ],
                    filters: [
                        {
                            'type: "URLRewrite",
                            urlRewrite: {
                                path: {
                                    'type: "ReplaceFullPath",
                                    replaceFullPath: "/employee/one/two"
                                }
                            }
                        }
                    ],
                    backendRefs: [
                        {
                            group: "dp.wso2.com",
                            kind: "Backend",
                            name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                        }
                    ]
                },
                {
                    matches: [
                        {
                            path: {
                                'type: "RegularExpression",
                                value: "/employee/one"
                            },
                            method: "GET"
                        }
                    ],
                    filters: [
                        {
                            'type: "URLRewrite",
                            urlRewrite: {
                                path: {
                                    'type: "ReplaceFullPath",
                                    replaceFullPath: "/employee/one"
                                }
                            }
                        }
                    ],
                    backendRefs: [
                        {
                            group: "dp.wso2.com",
                            kind: "Backend",
                            name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                        }
                    ]
                },
                {
                    matches: [
                        {
                            path: {
                                'type: "RegularExpression",
                                value: "/employee"
                            },
                            method: "GET"
                        }
                    ],
                    filters: [
                        {
                            'type: "URLRewrite",
                            urlRewrite: {
                                path: {
                                    'type: "ReplaceFullPath",
                                    replaceFullPath: "/employee"
                                }
                            }
                        }
                    ],
                    backendRefs: [
                        {
                            group: "dp.wso2.com",
                            kind: "Backend",
                            name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                        }
                    ]
                }
            ]
        }
    };

    map<[model:Httproute[], model:Httproute[]]> routesMap = {
        "1": [
            [
                {
                    apiVersion: "gateway.networking.k8s.io/v1beta1",
                    kind: "HTTPRoute",
                    metadata: {name: "01ee37b2-57b1-12ee-b8c5-e9b11538a0c9"},
                    spec: {
                        rules: [
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            },
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee/one"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee/one"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            },
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee/one/two"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee/one/two"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            },
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee/one/two/three"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee/one/two/three"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            }
                        ]
                    }
                }
            ],
            [
                sortedHttpRoute
            ]
        ],
        "2": [
            [
                {
                    apiVersion: "gateway.networking.k8s.io/v1beta1",
                    kind: "HTTPRoute",
                    metadata: {name: "01ee37b2-57b1-12ee-b8c5-e9b11538a0c9"},
                    spec: {
                        rules: [
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee/one"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee/one"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            },
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            },
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee/one/two"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee/one/two"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            },
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee/one/two/three"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee/one/two/three"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            }
                        ]
                    }
                }
            ],
            [
                sortedHttpRoute
            ]
        ],
        "3": [
            [
                {
                    apiVersion: "gateway.networking.k8s.io/v1beta1",
                    kind: "HTTPRoute",
                    metadata: {name: "01ee37b2-57b1-12ee-b8c5-e9b11538a0c9"},
                    spec: {
                        rules: [
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee/one/two/three"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee/one/two/three"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            },
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee/one/two"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee/one/two"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            },
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee/one"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee/one"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            },
                            {
                                matches: [
                                    {
                                        path: {
                                            'type: "RegularExpression",
                                            value: "/employee"
                                        },
                                        method: "GET"
                                    }
                                ],
                                filters: [
                                    {
                                        'type: "URLRewrite",
                                        urlRewrite: {
                                            path: {
                                                'type: "ReplaceFullPath",
                                                replaceFullPath: "/employee"
                                            }
                                        }
                                    }
                                ],
                                backendRefs: [
                                    {
                                        group: "dp.wso2.com",
                                        kind: "Backend",
                                        name: "backend-71eb6a67aab2b394c78202cc4ada0962c64c819b-api"
                                    }
                                ]
                            }
                        ]
                    }
                }
            ],
            [
                sortedHttpRoute
            ]
        ]
    };
    return routesMap;
}
