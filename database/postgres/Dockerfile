#---------------------------------------------------------------
#
# Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
#
# WSO2 LLC. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
#---------------------------------------------------------------


FROM library/postgres
ENV POSTGRES_USER wso2carbon
ENV POSTGRES_PASSWORD wso2carbon
ENV POSTGRES_DB WSO2AM_DB
ENV PGDATA /tmp/data1
COPY db.sql /docker-entrypoint-initdb.d/
COPY idp.sql /docker-entrypoint-initdb.d/

USER root 

RUN groupmod -g 10001 postgres
RUN usermod -u 10001 postgres

# this 777 will be replaced by 700 at runtime (allows semi-arbitrary "--user" values)
RUN mkdir -p "$PGDATA" && chown -R postgres:postgres "$PGDATA" && chmod 777 "$PGDATA"
USER 10001
