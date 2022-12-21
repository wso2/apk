public type Status record {
    string apiVersion;
    int code;
    string kind;
    string message;
    string reason;
    string status;
    StatusDetails details;
    ListMeta metadata;
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
    string uid;
};

public type ListMeta record {
    string 'continue;
    int remainingItemCount;
    string resourceVersion;
    string selfLink;
};
