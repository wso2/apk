# Authenticated UsersContext.
#
# + username - username of the user 
# + userId - user Id of the user
# + organization - organization of the user
# + claims - Field Description
public type UserContext record{|
    string username;
    string userId?;
    Organization organization;
    map<anydata> claims = {};
|};

public type User record {|
    string uuid;
    string IDPUserName;
|};