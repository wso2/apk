    \c "WSO2AM_DB"
        BEGIN TRANSACTION;

        -- Start Creating non-production-idp tables --

        CREATE TABLE CONSUMER_APPS (
            CONSUMER_KEY VARCHAR(255),
            CONSUMER_SECRET VARCHAR(2048),
            APP_NAME VARCHAR(255),
            CALLBACK_URL VARCHAR(2048),
            GRANT_TYPES VARCHAR (1024),
            CONSTRAINT CONSUMER_KEY_CONSTRAINT UNIQUE (CONSUMER_KEY),
            PRIMARY KEY (CONSUMER_KEY)
        );
        commit;        