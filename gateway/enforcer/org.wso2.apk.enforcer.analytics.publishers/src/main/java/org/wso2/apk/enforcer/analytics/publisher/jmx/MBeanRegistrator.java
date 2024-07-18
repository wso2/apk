/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */
package org.wso2.apk.enforcer.analytics.publisher.jmx;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.enforcer.analytics.publisher.jmx.impl.ExtAuthMetrics;
import org.wso2.apk.enforcer.analytics.publisher.reporter.prometheus.APIInvocationEvent;

import javax.management.*;
import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;

/**
 * The class which is responsible for registering MBeans.
 */
public class MBeanRegistrator {
    private static final Log logger = LogFactory.getLog(MBeanRegistrator.class);
    private static List<ObjectName> mBeans = new ArrayList<>();

    private static final String SERVER_PACKAGE = "org.wso2.apk.enforcer.analytics";

    private MBeanRegistrator() {
    }

    /**
     * Registers an object as an MBean with the MBean server.
     *
     * @param mBeanInstance - The MBean to be registered as an MBean.
     */
    public static void registerMBean(Object mBeanInstance, APIInvocationEvent event) throws RuntimeException, UnsupportedEncodingException {

        if (JMXUtils.isJMXMetricsEnabled()) {
            String className = mBeanInstance.getClass().getName();
            if (className.indexOf('.') != -1) {
                className = className.substring(className.lastIndexOf('.') + 1);
            }

            ExtAuthMetrics extAuthMetrics = (ExtAuthMetrics) mBeanInstance;
            logger.info("extAuthMetrics Object data: " + extAuthMetrics.toString());
            //logger.info("extAuthMetrics Object api name: " + extAuthMetrics.getApiName());

            String objectName = String.format(
                    "%s:type=%s,apiName=%s,apiContext=%s,proxyResponseCode=%d," +
                            "destination=%s,apiCreatorTenantDomain=%s,platform=%s," +
                            "organizationId=%s,apiMethod=%s,apiVersion=%s,environmentId=%s," +
                            "gatewayType=%s,apiCreator=%s,responseCacheHit=%b,backendLatency=%d," +
                            "correlationId=%s,requestMediationLatency=%d,keyType=%s,apiId=%s," +
                            "applicationName=%s,targetResponseCode=%d,applicationOwner=%s," +
                            "userAgent=%s,userName=%s,apiResourceTemplate=%s,regionId=%s,responseLatency=%d," +
                            "responseMediationLatency=%d,userIp=%s,applicationId=%s,apiType=%s,xOriginalGwUrl=%s",
                    SERVER_PACKAGE,
                    className,
                    event.getApiName(),
                    event.getApiContext(),
                    event.getProxyResponseCode(),
                    URLEncoder.encode(event.getDestination(), StandardCharsets.UTF_8.toString()),
                    event.getApiCreatorTenantDomain(),
                    event.getPlatform(),
                    event.getOrganizationId(),
                    event.getApiMethod(),
                    event.getApiVersion(),
                    event.getEnvironmentId(),
                    event.getGatewayType(),
                    event.getApiCreator(),
                    event.isResponseCacheHit(),
                    event.getBackendLatency(),
                    event.getCorrelationId(),
                    event.getRequestMediationLatency(),
                    event.getKeyType(),
                    event.getApiId(),
                    "Resident Key Manager",
                    event.getTargetResponseCode(),
                    event.getApplicationOwner(),
                    event.getUserAgent(),
                    event.getUserName(),
                    event.getApiResourceTemplate(),
                    event.getRegionId(),
                    event.getResponseLatency(),
                    event.getResponseMediationLatency(),
                    event.getUserIp(),
                    event.getApplicationId(),
                    event.getApiType(),
                    URLEncoder.encode(event.getProperties().get("x-original-gw-url"), StandardCharsets.UTF_8.toString())
            );
            logger.info("Registering MBean with object name: " + objectName);
            try {
                MBeanServer mBeanServer = MBeanManagementFactory.getMBeanServer();
                Set set = mBeanServer.queryNames(new ObjectName(objectName), null);
                if (set.isEmpty()) {
                    try {
                        ObjectName name = new ObjectName(objectName);
                        mBeanServer.registerMBean(mBeanInstance, name);
                        mBeans.add(name);
                        logger.info("MBean registered successfully with object name: " + name);
                        logger.info("Mbeans: " + mBeans);
                    } catch (InstanceAlreadyExistsException e) {
                        String msg = "MBean " + objectName + " already exists";
                        logger.error(msg, e);
                        throw new RuntimeException(msg, e);
                    } catch (MBeanRegistrationException | NotCompliantMBeanException e) {
                        String msg = "Execption when registering MBean";
                        logger.error(msg, e);
                        throw new RuntimeException(msg, e);
                    }
                } else {
                    String msg = "MBean " + objectName + " already exists";
                    logger.error(msg);
                    throw new RuntimeException(msg);
                }
            } catch (MalformedObjectNameException e) {
                String msg = "Could not register " + mBeanInstance.getClass() + " MBean";
                logger.error(msg, e);
                throw new RuntimeException(msg, e);
            }
        } else {
            logger.debug("JMX Metrics should be enabled to register MBean instance: {}");
        }
    }
}
