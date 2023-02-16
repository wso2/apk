# Authenticated UsersContext.
#
# + username - Field Description  
# + organization - Field Description  
# + claims - Field Description
public type UserContext record{|
    string username;
    Organization organization?;
    map<anydata> claims = {};
|};
