public type Endpoint record {|
string url?;
string namespace?;
string name?;
boolean serviceEntry = false;
int weight?;

|};
