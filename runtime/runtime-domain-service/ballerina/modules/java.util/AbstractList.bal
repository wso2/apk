import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import runtime_domain_service.java.lang as javalang;
import runtime_domain_service.java.util.'stream as javautilstream;
import runtime_domain_service.java.util.'function as javautilfunction;

# Ballerina class mapping for the Java `java.util.AbstractList` class.
@java:Binding {'class: "java.util.AbstractList"}
public distinct class AbstractList {

    *java:JObject;
    *AbstractCollection;

    # The `handle` field that stores the reference to the `java.util.AbstractList` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `java.util.AbstractList` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `java.util.AbstractList` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `add` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function add(javalang:Object arg0) returns boolean {
        return java_util_AbstractList_add(self.jObj, arg0.jObj);
    }

    # The function that maps to the `add` method of `java.util.AbstractList`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `javalang:Object` value required to map with the Java method parameter.
    public function add2(int arg0, javalang:Object arg1) {
        java_util_AbstractList_add2(self.jObj, arg0, arg1.jObj);
    }

    # The function that maps to the `addAll` method of `java.util.AbstractList`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function addAll(Collection arg0) returns boolean {
        return java_util_AbstractList_addAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `addAll` method of `java.util.AbstractList`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function addAll2(int arg0, Collection arg1) returns boolean {
        return java_util_AbstractList_addAll2(self.jObj, arg0, arg1.jObj);
    }

    # The function that maps to the `clear` method of `java.util.AbstractList`.
    public function clear() {
        java_util_AbstractList_clear(self.jObj);
    }

    # The function that maps to the `contains` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function contains(javalang:Object arg0) returns boolean {
        return java_util_AbstractList_contains(self.jObj, arg0.jObj);
    }

    # The function that maps to the `containsAll` method of `java.util.AbstractList`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function containsAll(Collection arg0) returns boolean {
        return java_util_AbstractList_containsAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `equals` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return java_util_AbstractList_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `forEach` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javautilfunction:Consumer` value required to map with the Java method parameter.
    public function forEach(javautilfunction:Consumer arg0) {
        java_util_AbstractList_forEach(self.jObj, arg0.jObj);
    }

    # The function that maps to the `get` method of `java.util.AbstractList`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function get(int arg0) returns javalang:Object {
        handle externalObj = java_util_AbstractList_get(self.jObj, arg0);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getClass` method of `java.util.AbstractList`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = java_util_AbstractList_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `java.util.AbstractList`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return java_util_AbstractList_hashCode(self.jObj);
    }

    # The function that maps to the `indexOf` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `int` value returning from the Java mapping.
    public function indexOf(javalang:Object arg0) returns int {
        return java_util_AbstractList_indexOf(self.jObj, arg0.jObj);
    }

    # The function that maps to the `isEmpty` method of `java.util.AbstractList`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public function isEmpty() returns boolean {
        return java_util_AbstractList_isEmpty(self.jObj);
    }

    # The function that maps to the `iterator` method of `java.util.AbstractList`.
    #
    # + return - The `Iterator` value returning from the Java mapping.
    public function iterator() returns Iterator {
        handle externalObj = java_util_AbstractList_iterator(self.jObj);
        Iterator newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `lastIndexOf` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `int` value returning from the Java mapping.
    public function lastIndexOf(javalang:Object arg0) returns int {
        return java_util_AbstractList_lastIndexOf(self.jObj, arg0.jObj);
    }

    # The function that maps to the `listIterator` method of `java.util.AbstractList`.
    #
    # + return - The `ListIterator` value returning from the Java mapping.
    public function listIterator() returns ListIterator {
        handle externalObj = java_util_AbstractList_listIterator(self.jObj);
        ListIterator newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `listIterator` method of `java.util.AbstractList`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `ListIterator` value returning from the Java mapping.
    public function listIterator2(int arg0) returns ListIterator {
        handle externalObj = java_util_AbstractList_listIterator2(self.jObj, arg0);
        ListIterator newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `notify` method of `java.util.AbstractList`.
    public function notify() {
        java_util_AbstractList_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `java.util.AbstractList`.
    public function notifyAll() {
        java_util_AbstractList_notifyAll(self.jObj);
    }

    # The function that maps to the `parallelStream` method of `java.util.AbstractList`.
    #
    # + return - The `javautilstream:Stream` value returning from the Java mapping.
    public function parallelStream() returns javautilstream:Stream {
        handle externalObj = java_util_AbstractList_parallelStream(self.jObj);
        javautilstream:Stream newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `remove` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function remove(javalang:Object arg0) returns boolean {
        return java_util_AbstractList_remove(self.jObj, arg0.jObj);
    }

    # The function that maps to the `remove` method of `java.util.AbstractList`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function remove2(int arg0) returns javalang:Object {
        handle externalObj = java_util_AbstractList_remove2(self.jObj, arg0);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `removeAll` method of `java.util.AbstractList`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function removeAll(Collection arg0) returns boolean {
        return java_util_AbstractList_removeAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `removeIf` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javautilfunction:Predicate` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function removeIf(javautilfunction:Predicate arg0) returns boolean {
        return java_util_AbstractList_removeIf(self.jObj, arg0.jObj);
    }

    # The function that maps to the `replaceAll` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javautilfunction:UnaryOperator` value required to map with the Java method parameter.
    public function replaceAll(javautilfunction:UnaryOperator arg0) {
        java_util_AbstractList_replaceAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `retainAll` method of `java.util.AbstractList`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function retainAll(Collection arg0) returns boolean {
        return java_util_AbstractList_retainAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `set` method of `java.util.AbstractList`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function set(int arg0, javalang:Object arg1) returns javalang:Object {
        handle externalObj = java_util_AbstractList_set(self.jObj, arg0, arg1.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `size` method of `java.util.AbstractList`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function size() returns int {
        return java_util_AbstractList_size(self.jObj);
    }

    # The function that maps to the `sort` method of `java.util.AbstractList`.
    #
    # + arg0 - The `Comparator` value required to map with the Java method parameter.
    public function sort(Comparator arg0) {
        java_util_AbstractList_sort(self.jObj, arg0.jObj);
    }

    # The function that maps to the `spliterator` method of `java.util.AbstractList`.
    #
    # + return - The `Spliterator` value returning from the Java mapping.
    public function spliterator() returns Spliterator {
        handle externalObj = java_util_AbstractList_spliterator(self.jObj);
        Spliterator newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `stream` method of `java.util.AbstractList`.
    #
    # + return - The `javautilstream:Stream` value returning from the Java mapping.
    public function 'stream() returns javautilstream:Stream {
        handle externalObj = java_util_AbstractList_stream(self.jObj);
        javautilstream:Stream newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `subList` method of `java.util.AbstractList`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `List` value returning from the Java mapping.
    public function subList(int arg0, int arg1) returns List {
        handle externalObj = java_util_AbstractList_subList(self.jObj, arg0, arg1);
        List newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `toArray` method of `java.util.AbstractList`.
    #
    # + return - The `javalang:Object[]` value returning from the Java mapping.
    public function toArray() returns javalang:Object[]|error {
        handle externalObj = java_util_AbstractList_toArray(self.jObj);
        javalang:Object[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:Object element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `toArray` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javautilfunction:IntFunction` value required to map with the Java method parameter.
    # + return - The `javalang:Object[]` value returning from the Java mapping.
    public function toArray2(javautilfunction:IntFunction arg0) returns javalang:Object[]|error {
        handle externalObj = java_util_AbstractList_toArray2(self.jObj, arg0.jObj);
        javalang:Object[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:Object element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `toArray` method of `java.util.AbstractList`.
    #
    # + arg0 - The `javalang:Object[]` value required to map with the Java method parameter.
    # + return - The `javalang:Object[]` value returning from the Java mapping.
    public function toArray3(javalang:Object[] arg0) returns javalang:Object[]|error {
        handle externalObj = java_util_AbstractList_toArray3(self.jObj, check jarrays:toHandle(arg0, "java.lang.Object"));
        javalang:Object[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:Object element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `wait` method of `java.util.AbstractList`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = java_util_AbstractList_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `java.util.AbstractList`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = java_util_AbstractList_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `java.util.AbstractList`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = java_util_AbstractList_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

function java_util_AbstractList_add(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "add",
    'class: "java.util.AbstractList",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_AbstractList_add2(handle receiver, int arg0, handle arg1) = @java:Method {
    name: "add",
    'class: "java.util.AbstractList",
    paramTypes: ["int", "java.lang.Object"]
} external;

function java_util_AbstractList_addAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "addAll",
    'class: "java.util.AbstractList",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_AbstractList_addAll2(handle receiver, int arg0, handle arg1) returns boolean = @java:Method {
    name: "addAll",
    'class: "java.util.AbstractList",
    paramTypes: ["int", "java.util.Collection"]
} external;

function java_util_AbstractList_clear(handle receiver) = @java:Method {
    name: "clear",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_contains(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "contains",
    'class: "java.util.AbstractList",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_AbstractList_containsAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "containsAll",
    'class: "java.util.AbstractList",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_AbstractList_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "java.util.AbstractList",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_AbstractList_forEach(handle receiver, handle arg0) = @java:Method {
    name: "forEach",
    'class: "java.util.AbstractList",
    paramTypes: ["java.util.function.Consumer"]
} external;

function java_util_AbstractList_get(handle receiver, int arg0) returns handle = @java:Method {
    name: "get",
    'class: "java.util.AbstractList",
    paramTypes: ["int"]
} external;

function java_util_AbstractList_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_indexOf(handle receiver, handle arg0) returns int = @java:Method {
    name: "indexOf",
    'class: "java.util.AbstractList",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_AbstractList_isEmpty(handle receiver) returns boolean = @java:Method {
    name: "isEmpty",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_iterator(handle receiver) returns handle = @java:Method {
    name: "iterator",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_lastIndexOf(handle receiver, handle arg0) returns int = @java:Method {
    name: "lastIndexOf",
    'class: "java.util.AbstractList",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_AbstractList_listIterator(handle receiver) returns handle = @java:Method {
    name: "listIterator",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_listIterator2(handle receiver, int arg0) returns handle = @java:Method {
    name: "listIterator",
    'class: "java.util.AbstractList",
    paramTypes: ["int"]
} external;

function java_util_AbstractList_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_parallelStream(handle receiver) returns handle = @java:Method {
    name: "parallelStream",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_remove(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "remove",
    'class: "java.util.AbstractList",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_AbstractList_remove2(handle receiver, int arg0) returns handle = @java:Method {
    name: "remove",
    'class: "java.util.AbstractList",
    paramTypes: ["int"]
} external;

function java_util_AbstractList_removeAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "removeAll",
    'class: "java.util.AbstractList",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_AbstractList_removeIf(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "removeIf",
    'class: "java.util.AbstractList",
    paramTypes: ["java.util.function.Predicate"]
} external;

function java_util_AbstractList_replaceAll(handle receiver, handle arg0) = @java:Method {
    name: "replaceAll",
    'class: "java.util.AbstractList",
    paramTypes: ["java.util.function.UnaryOperator"]
} external;

function java_util_AbstractList_retainAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "retainAll",
    'class: "java.util.AbstractList",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_AbstractList_set(handle receiver, int arg0, handle arg1) returns handle = @java:Method {
    name: "set",
    'class: "java.util.AbstractList",
    paramTypes: ["int", "java.lang.Object"]
} external;

function java_util_AbstractList_size(handle receiver) returns int = @java:Method {
    name: "size",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_sort(handle receiver, handle arg0) = @java:Method {
    name: "sort",
    'class: "java.util.AbstractList",
    paramTypes: ["java.util.Comparator"]
} external;

function java_util_AbstractList_spliterator(handle receiver) returns handle = @java:Method {
    name: "spliterator",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_stream(handle receiver) returns handle = @java:Method {
    name: "stream",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_subList(handle receiver, int arg0, int arg1) returns handle = @java:Method {
    name: "subList",
    'class: "java.util.AbstractList",
    paramTypes: ["int", "int"]
} external;

function java_util_AbstractList_toArray(handle receiver) returns handle = @java:Method {
    name: "toArray",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_toArray2(handle receiver, handle arg0) returns handle = @java:Method {
    name: "toArray",
    'class: "java.util.AbstractList",
    paramTypes: ["java.util.function.IntFunction"]
} external;

function java_util_AbstractList_toArray3(handle receiver, handle arg0) returns handle = @java:Method {
    name: "toArray",
    'class: "java.util.AbstractList",
    paramTypes: ["[Ljava.lang.Object;"]
} external;

function java_util_AbstractList_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "java.util.AbstractList",
    paramTypes: []
} external;

function java_util_AbstractList_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "java.util.AbstractList",
    paramTypes: ["long"]
} external;

function java_util_AbstractList_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "java.util.AbstractList",
    paramTypes: ["long", "int"]
} external;

