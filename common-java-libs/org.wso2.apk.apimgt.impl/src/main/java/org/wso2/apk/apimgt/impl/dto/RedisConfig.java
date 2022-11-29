/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.apimgt.impl.dto;

/**
 * RedisConfig Model class for connection properties of a Redis Server
 */
public class RedisConfig {

    private boolean isRedisEnabled;
    private String host;
    private int port;
    private String user;
    private char[] password;
    private int databaseId;
    private int connectionTimeout;
    private boolean isSslEnabled;
    private int maxTotal = 8;
    private int maxIdle = 8;
    private int minIdle = 0;
    private boolean testOnBorrow = false;
    private boolean testOnReturn = false;
    private boolean testWhileIdle = true;

    private boolean blockWhenExhausted = true;
    private long minEvictableIdleTimeMillis = 60000L;
    private long timeBetweenEvictionRunsMillis = 30000L;
    private int numTestsPerEvictionRun = -1;

    public int getMaxTotal() {

        return maxTotal;
    }

    public void setMaxTotal(int maxTotal) {

        this.maxTotal = maxTotal;
    }

    public int getMaxIdle() {

        return maxIdle;
    }

    public void setMaxIdle(int maxIdle) {

        this.maxIdle = maxIdle;
    }

    public int getMinIdle() {

        return minIdle;
    }

    public void setMinIdle(int minIdle) {

        this.minIdle = minIdle;
    }

    public boolean isTestOnBorrow() {

        return testOnBorrow;
    }

    public void setTestOnBorrow(boolean testOnBorrow) {

        this.testOnBorrow = testOnBorrow;
    }

    public boolean isTestOnReturn() {

        return testOnReturn;
    }

    public void setTestOnReturn(boolean testOnReturn) {

        this.testOnReturn = testOnReturn;
    }

    public boolean isTestWhileIdle() {

        return testWhileIdle;
    }

    public void setTestWhileIdle(boolean testWhileIdle) {

        this.testWhileIdle = testWhileIdle;
    }

    public boolean isBlockWhenExhausted() {

        return blockWhenExhausted;
    }

    public void setBlockWhenExhausted(boolean blockWhenExhausted) {

        this.blockWhenExhausted = blockWhenExhausted;
    }

    public long getMinEvictableIdleTimeMillis() {

        return minEvictableIdleTimeMillis;
    }

    public void setMinEvictableIdleTimeMillis(long minEvictableIdleTimeMillis) {

        this.minEvictableIdleTimeMillis = minEvictableIdleTimeMillis;
    }

    public long getTimeBetweenEvictionRunsMillis() {

        return timeBetweenEvictionRunsMillis;
    }

    public void setTimeBetweenEvictionRunsMillis(long timeBetweenEvictionRunsMillis) {

        this.timeBetweenEvictionRunsMillis = timeBetweenEvictionRunsMillis;
    }

    public int getNumTestsPerEvictionRun() {

        return numTestsPerEvictionRun;
    }

    public void setNumTestsPerEvictionRun(int numTestsPerEvictionRun) {

        this.numTestsPerEvictionRun = numTestsPerEvictionRun;
    }

    /**
     * Public default constructor
     */
    public RedisConfig() {

    }

    public boolean isRedisEnabled() {

        return isRedisEnabled;
    }

    public void setRedisEnabled(boolean redisEnabled) {

        isRedisEnabled = redisEnabled;
    }

    public String getHost() {

        return host;
    }

    public void setHost(String host) {

        this.host = host;
    }

    public int getPort() {

        return port;
    }

    public void setPort(int port) {

        this.port = port;
    }

    public String getUser() {

        return user;
    }

    public void setUser(String user) {

        this.user = user;
    }

    public char[] getPassword() {

        return password;
    }

    public void setPassword(char[] password) {

        this.password = password;
    }

    public int getDatabaseId() {

        return databaseId;
    }

    public void setDatabaseId(int databaseId) {

        this.databaseId = databaseId;
    }

    public int getConnectionTimeout() {

        return connectionTimeout;
    }

    public void setConnectionTimeout(int connectionTimeout) {

        this.connectionTimeout = connectionTimeout;
    }

    public boolean isSslEnabled() {

        return isSslEnabled;
    }

    public void setSslEnabled(boolean sslEnabled) {

        isSslEnabled = sslEnabled;
    }
}
