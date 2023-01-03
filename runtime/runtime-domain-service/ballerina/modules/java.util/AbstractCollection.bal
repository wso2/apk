import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import runtime_domain_service.java.lang as javalang;
import runtime_domain_service.java.util.'stream as javautilstream;
import runtime_domain_service.java.util.'function as javautilfunction;

# Ballerina class mapping for the Java `java.util.AbstractCollection` class.
@java:Binding {'class: "java.util.AbstractCollection"}
public distinct class AbstractCollection {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `java.util.AbstractCollection` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `java.util.AbstractCollection` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `java.util.AbstractCollection` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `add` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function add(javalang:Object arg0) returns boolean {
        return java_util_AbstractCollection_add(self.jObj, arg0.jObj);
    }

    # The function that maps to the `addAll` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function addAll(Collection arg0) returns boolean {
        return java_util_AbstractCollection_addAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `clear` method of `java.util.AbstractCollection`.
    public function clear() {
        java_util_AbstractCollection_clear(self.jObj);
    }

    # The function that maps to the `contains` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function contains(javalang:Object arg0) returns boolean {
        return java_util_AbstractCollection_contains(self.jObj, arg0.jObj);
    }

    # The function that maps to the `containsAll` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function containsAll(Collection arg0) returns boolean {
        return java_util_AbstractCollection_containsAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `equals` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return java_util_AbstractCollection_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `forEach` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `javautilfunction:Consumer` value required to map with the Java method parameter.
    public function forEach(javautilfunction:Consumer arg0) {
        java_util_AbstractCollection_forEach(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `java.util.AbstractCollection`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = java_util_AbstractCollection_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `java.util.AbstractCollection`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return java_util_AbstractCollection_hashCode(self.jObj);
    }

    # The function that maps to the `isEmpty` method of `java.util.AbstractCollection`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public function isEmpty() returns boolean {
        return java_util_AbstractCollection_isEmpty(self.jObj);
    }

    # The function that maps to the `iterator` method of `java.util.AbstractCollection`.
    #
    # + return - The `Iterator` value returning from the Java mapping.
    public function iterator() returns Iterator {
        handle externalObj = java_util_AbstractCollection_iterator(self.jObj);
        Iterator newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `notify` method of `java.util.AbstractCollection`.
    public function notify() {
        java_util_AbstractCollection_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `java.util.AbstractCollection`.
    public function notifyAll() {
        java_util_AbstractCollection_notifyAll(self.jObj);
    }

    # The function that maps to the `parallelStream` method of `java.util.AbstractCollection`.
    #
    # + return - The `javautilstream:Stream` value returning from the Java mapping.
    public function parallelStream() returns javautilstream:Stream {
        handle externalObj = java_util_AbstractCollection_parallelStream(self.jObj);
        javautilstream:Stream newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `remove` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function remove(javalang:Object arg0) returns boolean {
        return java_util_AbstractCollection_remove(self.jObj, arg0.jObj);
    }

    # The function that maps to the `removeAll` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function removeAll(Collection arg0) returns boolean {
        return java_util_AbstractCollection_removeAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `removeIf` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `javautilfunction:Predicate` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function removeIf(javautilfunction:Predicate arg0) returns boolean {
        return java_util_AbstractCollection_removeIf(self.jObj, arg0.jObj);
    }

    # The function that maps to the `retainAll` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function retainAll(Collection arg0) returns boolean {
        return java_util_AbstractCollection_retainAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `size` method of `java.util.AbstractCollection`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function size() returns int {
        return java_util_AbstractCollection_size(self.jObj);
    }

    # The function that maps to the `spliterator` method of `java.util.AbstractCollection`.
    #
    # + return - The `Spliterator` value returning from the Java mapping.
    public function spliterator() returns Spliterator {
        handle externalObj = java_util_AbstractCollection_spliterator(self.jObj);
        Spliterator newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `stream` method of `java.util.AbstractCollection`.
    #
    # + return - The `javautilstream:Stream` value returning from the Java mapping.
    public function 'stream() returns javautilstream:Stream {
        handle externalObj = java_util_AbstractCollection_stream(self.jObj);
        javautilstream:Stream newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `toArray` method of `java.util.AbstractCollection`.
    #
    # + return - The `javalang:Object[]` value returning from the Java mapping.
    public function toArray() returns javalang:Object[]|error {
        handle externalObj = java_util_AbstractCollection_toArray(self.jObj);
        javalang:Object[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:Object element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `toArray` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `javautilfunction:IntFunction` value required to map with the Java method parameter.
    # + return - The `javalang:Object[]` value returning from the Java mapping.
    public function toArray2(javautilfunction:IntFunction arg0) returns javalang:Object[]|error {
        handle externalObj = java_util_AbstractCollection_toArray2(self.jObj, arg0.jObj);
        javalang:Object[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:Object element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `toArray` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `javalang:Object[]` value required to map with the Java method parameter.
    # + return - The `javalang:Object[]` value returning from the Java mapping.
    public function toArray3(javalang:Object[] arg0) returns javalang:Object[]|error {
        handle externalObj = java_util_AbstractCollection_toArray3(self.jObj, check jarrays:toHandle(arg0, "java.lang.Object"));
        javalang:Object[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:Object element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `wait` method of `java.util.AbstractCollection`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = java_util_AbstractCollection_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = java_util_AbstractCollection_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `java.util.AbstractCollection`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = java_util_AbstractCollection_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

function java_util_AbstractCollection_add(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "add",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_AbstractCollection_addAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "addAll",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_AbstractCollection_clear(handle receiver) = @java:Method {
    name: "clear",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_contains(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "contains",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_AbstractCollection_containsAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "containsAll",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_AbstractCollection_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_AbstractCollection_forEach(handle receiver, handle arg0) = @java:Method {
    name: "forEach",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.util.function.Consumer"]
} external;

function java_util_AbstractCollection_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_isEmpty(handle receiver) returns boolean = @java:Method {
    name: "isEmpty",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_iterator(handle receiver) returns handle = @java:Method {
    name: "iterator",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_parallelStream(handle receiver) returns handle = @java:Method {
    name: "parallelStream",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_remove(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "remove",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_AbstractCollection_removeAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "removeAll",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_AbstractCollection_removeIf(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "removeIf",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.util.function.Predicate"]
} external;

function java_util_AbstractCollection_retainAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "retainAll",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_AbstractCollection_size(handle receiver) returns int = @java:Method {
    name: "size",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_spliterator(handle receiver) returns handle = @java:Method {
    name: "spliterator",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_stream(handle receiver) returns handle = @java:Method {
    name: "stream",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_toArray(handle receiver) returns handle = @java:Method {
    name: "toArray",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_toArray2(handle receiver, handle arg0) returns handle = @java:Method {
    name: "toArray",
    'class: "java.util.AbstractCollection",
    paramTypes: ["java.util.function.IntFunction"]
} external;

function java_util_AbstractCollection_toArray3(handle receiver, handle arg0) returns handle = @java:Method {
    name: "toArray",
    'class: "java.util.AbstractCollection",
    paramTypes: ["[Ljava.lang.Object;"]
} external;

function java_util_AbstractCollection_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "java.util.AbstractCollection",
    paramTypes: []
} external;

function java_util_AbstractCollection_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "java.util.AbstractCollection",
    paramTypes: ["long"]
} external;

function java_util_AbstractCollection_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "java.util.AbstractCollection",
    paramTypes: ["long", "int"]
} external;

