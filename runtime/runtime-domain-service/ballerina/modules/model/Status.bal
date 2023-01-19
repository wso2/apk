public type Status record {
    string apiVersion = "v1";
    int code;
    string kind = "Status";
    string message?;
    string reason?;
    string status?;
    StatusDetails details?;
    ListMeta metadata?;
};

public type StatusCause record {
    string 'field;
    string message;
    string reason;
};

public type StatusDetails record {
    StatusCause[] 'causes;
    string group;
    string kind;
    string name;
    int retryAfterSeconds?;
    string uid?;
};

public type ListMeta record {
    string 'continue?;
    int remainingItemCount?;
    string resourceVersion?;
    string selfLink?;
};
