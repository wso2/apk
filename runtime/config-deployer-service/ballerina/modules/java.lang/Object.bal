import ballerina/jballerina.java;

# Ballerina class mapping for the Java `java.lang.Object` class.
@java:Binding {'class: "java.lang.Object"}
public distinct class Object {

    *java:JObject;

    # The `handle` field that stores the reference to the `java.lang.Object` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `java.lang.Object` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `java.lang.Object` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "";
    }

    # The function that maps to the `equals` method of `java.lang.Object`.
    #
    # + arg0 - The `Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(Object arg0) returns boolean {
        return java_lang_Object_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `java.lang.Object`.
    #
    # + return - The `Class` value returning from the Java mapping.
    public function getClass() returns Class {
        handle externalObj = java_lang_Object_getClass(self.jObj);
        Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `java.lang.Object`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return java_lang_Object_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `java.lang.Object`.
    public function notify() {
        java_lang_Object_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `java.lang.Object`.
    public function notifyAll() {
        java_lang_Object_notifyAll(self.jObj);
    }

    # The function that maps to the `wait` method of `java.lang.Object`.
    #
    # + return - The `InterruptedException` value returning from the Java mapping.
    public function 'wait() returns InterruptedException? {
        error|() externalObj = java_lang_Object_wait(self.jObj);
        if (externalObj is error) {
            InterruptedException e = error InterruptedException(INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `java.lang.Object`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns InterruptedException? {
        error|() externalObj = java_lang_Object_wait2(self.jObj, arg0);
        if (externalObj is error) {
            InterruptedException e = error InterruptedException(INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `java.lang.Object`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns InterruptedException? {
        error|() externalObj = java_lang_Object_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            InterruptedException e = error InterruptedException(INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `java.lang.Object`.
#
# + return - The new `Object` class generated.
public function newObject1() returns Object {
    handle externalObj = java_lang_Object_newObject1();
    Object newObj = new (externalObj);
    return newObj;
}

function java_lang_Object_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "java.lang.Object",
    paramTypes: ["java.lang.Object"]
} external;

function java_lang_Object_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "java.lang.Object",
    paramTypes: []
} external;

function java_lang_Object_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "java.lang.Object",
    paramTypes: []
} external;

function java_lang_Object_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "java.lang.Object",
    paramTypes: []
} external;

function java_lang_Object_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "java.lang.Object",
    paramTypes: []
} external;

function java_lang_Object_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "java.lang.Object",
    paramTypes: []
} external;

function java_lang_Object_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "java.lang.Object",
    paramTypes: ["long"]
} external;

function java_lang_Object_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "java.lang.Object",
    paramTypes: ["long", "int"]
} external;

function java_lang_Object_newObject1() returns handle = @java:Constructor {
    'class: "java.lang.Object",
    paramTypes: []
} external;

