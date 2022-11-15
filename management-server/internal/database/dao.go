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
 *
 */

package database

import apkmgt_application "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/apkmgt"

// GetApplicationByUUID returns the Application details from the DB for a given application
func GetApplicationByUUID(uuid string) (*apkmgt_application.Application, error) {
	rows, _ := ExecDBQuery(QueryGetApplicationByUUID, uuid)
	rows.Next()
	values, err := rows.Values()
	if err != nil {
		return nil, err
	} else {
		subs, _ := getSubscriptionsForApplication(uuid)
		keys, _ := getConsumerKeysForApplication(uuid)
		application := &apkmgt_application.Application{
			Uuid:          values[0].(string),
			Name:          values[1].(string),
			Owner:         "",  //ToDo : Check how to get Owner from db
			Attributes:    nil, //ToDo : check the values for Attributes
			Subscriber:    "",
			Organization:  values[3].(string),
			Subscriptions: subs,
			ConsumerKeys:  keys,
		}
		return application, nil
	}
}

// getSubscriptionsForApplication returns all subscriptions from DB, for a given application
func getSubscriptionsForApplication(appUuid string) ([]*apkmgt_application.Subscription, error) {
	rows, err := ExecDBQuery(QueryGetAllSubscriptionsForApplication, appUuid)
	if err != nil {
	}
	var subs []*apkmgt_application.Subscription
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		} else {
			subs = append(subs, &apkmgt_application.Subscription{
				Uuid:               values[0].(string),
				ApiUuid:            values[1].(string),
				PolicyId:           "",
				SubscriptionStatus: values[3].(string),
				Organization:       values[4].(string),
				CreatedBy:          values[5].(string),
			})
		}
	}
	return subs, nil
}

// getConsumerKeysForApplication returns all Consumer Keys from DB, for a given application
func getConsumerKeysForApplication(appUUID string) ([]*apkmgt_application.ConsumerKey, error) {
	rows, err := ExecDBQuery(QueryConsumerKeysForApplication, appUUID)
	if err != nil {
	}
	var keys []*apkmgt_application.ConsumerKey
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		} else {
			keys = append(keys, &apkmgt_application.ConsumerKey{
				Key:        values[0].(string),
				KeyManager: values[1].(string),
			})
		}
	}
	return keys, nil
}
