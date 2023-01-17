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
ENV PGDATA /var/lib/postgresql/data1
#COPY db.sql /docker-entrypoint-initdb.d/

USER root 

RUN groupmod -g 10001 postgres
RUN usermod -u 10001 postgres

RUN mkdir -p "$PGDATA"
RUN chown -R postgres "$PGDATA"
RUN chmod 777 "$PGDATA"

RUN mkdir -p /var/run/postgresql1
RUN chown -R postgres /var/run/postgresql1
RUN chmod 775 /var/run/postgresql1

# this 777 will be replaced by 700 at runtime (allows semi-arbitrary "--user" values)
# RUN mkdir -p "$PGDATA" && chown -R postgres:postgres "$PGDATA" && chmod 777 "$PGDATA"

# RUN set -eux; \ 
# groupadd -r postgres1 --gid=10001; \
# # https://salsa.debian.org/postgresql/postgresql-common/blob/997d842ee744687d99a2b2d95c1083a2615c79e8/debian/postgresql-common.postinst#L32-35
# useradd -r -g postgres1 --uid=10001 --home-dir=/var/lib/postgresql --shell=/bin/bash postgres1; \
# # also create the postgres user's home directory with appropriate permissions
# # see https://github.com/docker-library/postgres/issues/274
# chown -R 10001:10001 /var/lib/postgresql 

#RUN chown -R 10001:10001 /var/lib/postgresql/data1

USER 10001
