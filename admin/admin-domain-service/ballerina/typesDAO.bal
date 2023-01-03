
public type ApplicationRatePlanDAO record {
    *Policy;
    string defaulLimitType;
    # Unit of the time. Allowed values are "sec", "min", "hour", "day"
    string timeUnit;
    # Time limit that the throttling limit applies.
    int unitTime;
    int quota;
    # Unit of data allowed to be transfered. Allowed values are "KB", "MB" and "GB"
    string dataUnit?;
};

public type BusinessPlanDAO record {
    *Policy;
    *GraphQLQuery;
    string defaulLimitType;
    # Unit of the time. Allowed values are "sec", "min", "hour", "day"
    string timeUnit;
    # Time limit that the throttling limit applies.
    int unitTime;
    int quota;
    # Unit of data allowed to be transfered. Allowed values are "KB", "MB" and "GB"
    string dataUnit?;
    # Burst control request count
    int rateLimitCount?;
    # Burst control time unit
    string rateLimitTimeUnit?;
    # Number of subscriptions allowed
    int subscriberCount?;
    # Custom attributes added to the Subscription Throttling Policy
    CustomAttribute[] customAttributes?;
    BusinessPlanPermission permissions?;
};