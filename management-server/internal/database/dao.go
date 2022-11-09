package database

import apkmgt_application "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/apkmgt"

func GetApplicationByUUID(uuid string) (*apkmgt_application.Application, error) {
	rows, _ := ExecDBQuery(QueryGetApplicationByUUID, uuid)
	rows.Next()
	values, err := rows.Values()
	if err != nil {
		return nil, err
	} else {
		subs, _ := GetSubscriptionsForApplication(uuid)
		application := &apkmgt_application.Application{
			Uuid:          values[0].(string),
			Name:          values[1].(string),
			Owner:         "",
			Attributes:    nil,
			Subscriber:    values[2].(string),
			Organization:  values[3].(string),
			Subscriptions: subs,
			ConsumerKeys:  nil,
		}
		return application, nil
	}
}

func GetSubscriptionsForApplication(appUuid string) ([]*apkmgt_application.Subscription, error) {
	rows, _ := ExecDBQuery(QueryGetApplicationByUUID, appUuid)
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
