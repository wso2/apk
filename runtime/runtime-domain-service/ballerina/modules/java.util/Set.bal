import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import runtime_domain_service.java.lang as javalang;
import runtime_domain_service.java.util.'stream as javautilstream;
import runtime_domain_service.java.util.'function as javautilfunction;

# Ballerina class mapping for the Java `java.util.Set` interface.
@java:Binding {'class: "java.util.Set"}
public distinct class Set {

    *java:JObject;

    # The `handle` field that stores the reference to the `java.util.Set` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `java.util.Set` Java interface.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `java.util.Set` Java interface.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `add` method of `java.util.Set`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public isolated function add(javalang:Object arg0) returns boolean {
        return java_util_Set_add(self.jObj, arg0.jObj);
    }

    # The function that maps to the `addAll` method of `java.util.Set`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function addAll(Collection arg0) returns boolean {
        return java_util_Set_addAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `clear` method of `java.util.Set`.
    public function clear() {
        java_util_Set_clear(self.jObj);
    }

    # The function that maps to the `contains` method of `java.util.Set`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function contains(javalang:Object arg0) returns boolean {
        return java_util_Set_contains(self.jObj, arg0.jObj);
    }

    # The function that maps to the `containsAll` method of `java.util.Set`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function containsAll(Collection arg0) returns boolean {
        return java_util_Set_containsAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `equals` method of `java.util.Set`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return java_util_Set_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `forEach` method of `java.util.Set`.
    #
    # + arg0 - The `javautilfunction:Consumer` value required to map with the Java method parameter.
    public function forEach(javautilfunction:Consumer arg0) {
        java_util_Set_forEach(self.jObj, arg0.jObj);
    }

    # The function that maps to the `hashCode` method of `java.util.Set`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return java_util_Set_hashCode(self.jObj);
    }

    # The function that maps to the `isEmpty` method of `java.util.Set`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public function isEmpty() returns boolean {
        return java_util_Set_isEmpty(self.jObj);
    }

    # The function that maps to the `iterator` method of `java.util.Set`.
    #
    # + return - The `Iterator` value returning from the Java mapping.
    public function iterator() returns Iterator {
        handle externalObj = java_util_Set_iterator(self.jObj);
        Iterator newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `parallelStream` method of `java.util.Set`.
    #
    # + return - The `javautilstream:Stream` value returning from the Java mapping.
    public function parallelStream() returns javautilstream:Stream {
        handle externalObj = java_util_Set_parallelStream(self.jObj);
        javautilstream:Stream newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `remove` method of `java.util.Set`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function remove(javalang:Object arg0) returns boolean {
        return java_util_Set_remove(self.jObj, arg0.jObj);
    }

    # The function that maps to the `removeAll` method of `java.util.Set`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function removeAll(Collection arg0) returns boolean {
        return java_util_Set_removeAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `removeIf` method of `java.util.Set`.
    #
    # + arg0 - The `javautilfunction:Predicate` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function removeIf(javautilfunction:Predicate arg0) returns boolean {
        return java_util_Set_removeIf(self.jObj, arg0.jObj);
    }

    # The function that maps to the `retainAll` method of `java.util.Set`.
    #
    # + arg0 - The `Collection` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function retainAll(Collection arg0) returns boolean {
        return java_util_Set_retainAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `size` method of `java.util.Set`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function size() returns int {
        return java_util_Set_size(self.jObj);
    }

    # The function that maps to the `spliterator` method of `java.util.Set`.
    #
    # + return - The `Spliterator` value returning from the Java mapping.
    public function spliterator() returns Spliterator {
        handle externalObj = java_util_Set_spliterator(self.jObj);
        Spliterator newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `stream` method of `java.util.Set`.
    #
    # + return - The `javautilstream:Stream` value returning from the Java mapping.
    public function 'stream() returns javautilstream:Stream {
        handle externalObj = java_util_Set_stream(self.jObj);
        javautilstream:Stream newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `toArray` method of `java.util.Set`.
    #
    # + return - The `javalang:Object[]` value returning from the Java mapping.
    public function toArray() returns javalang:Object[]|error {
        handle externalObj = java_util_Set_toArray(self.jObj);
        javalang:Object[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:Object element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `toArray` method of `java.util.Set`.
    #
    # + arg0 - The `javautilfunction:IntFunction` value required to map with the Java method parameter.
    # + return - The `javalang:Object[]` value returning from the Java mapping.
    public function toArray2(javautilfunction:IntFunction arg0) returns javalang:Object[]|error {
        handle externalObj = java_util_Set_toArray2(self.jObj, arg0.jObj);
        javalang:Object[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:Object element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `toArray` method of `java.util.Set`.
    #
    # + arg0 - The `javalang:Object[]` value required to map with the Java method parameter.
    # + return - The `javalang:Object[]` value returning from the Java mapping.
    public function toArray3(javalang:Object[] arg0) returns javalang:Object[]|error {
        handle externalObj = java_util_Set_toArray3(self.jObj, check jarrays:toHandle(arg0, "java.lang.Object"));
        javalang:Object[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:Object element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

}

# The function that maps to the `copyOf` method of `java.util.Set`.
#
# + arg0 - The `Collection` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_copyOf(Collection arg0) returns Set {
    handle externalObj = java_util_Set_copyOf(arg0.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + return - The `Set` value returning from the Java mapping.
public function Set_of() returns Set {
    handle externalObj = java_util_Set_of();
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + arg3 - The `javalang:Object` value required to map with the Java method parameter.
# + arg4 - The `javalang:Object` value required to map with the Java method parameter.
# + arg5 - The `javalang:Object` value required to map with the Java method parameter.
# + arg6 - The `javalang:Object` value required to map with the Java method parameter.
# + arg7 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of10(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6, javalang:Object arg7) returns Set {
    handle externalObj = java_util_Set_of10(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj, arg7.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + arg3 - The `javalang:Object` value required to map with the Java method parameter.
# + arg4 - The `javalang:Object` value required to map with the Java method parameter.
# + arg5 - The `javalang:Object` value required to map with the Java method parameter.
# + arg6 - The `javalang:Object` value required to map with the Java method parameter.
# + arg7 - The `javalang:Object` value required to map with the Java method parameter.
# + arg8 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of11(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6, javalang:Object arg7, javalang:Object arg8) returns Set {
    handle externalObj = java_util_Set_of11(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj, arg7.jObj, arg8.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + arg3 - The `javalang:Object` value required to map with the Java method parameter.
# + arg4 - The `javalang:Object` value required to map with the Java method parameter.
# + arg5 - The `javalang:Object` value required to map with the Java method parameter.
# + arg6 - The `javalang:Object` value required to map with the Java method parameter.
# + arg7 - The `javalang:Object` value required to map with the Java method parameter.
# + arg8 - The `javalang:Object` value required to map with the Java method parameter.
# + arg9 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of12(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6, javalang:Object arg7, javalang:Object arg8, javalang:Object arg9) returns Set {
    handle externalObj = java_util_Set_of12(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj, arg7.jObj, arg8.jObj, arg9.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of2(javalang:Object arg0) returns Set {
    handle externalObj = java_util_Set_of2(arg0.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object[]` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of3(javalang:Object[] arg0) returns Set|error {
    handle externalObj = java_util_Set_of3(check jarrays:toHandle(arg0, "java.lang.Object"));
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of4(javalang:Object arg0, javalang:Object arg1) returns Set {
    handle externalObj = java_util_Set_of4(arg0.jObj, arg1.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of5(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2) returns Set {
    handle externalObj = java_util_Set_of5(arg0.jObj, arg1.jObj, arg2.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + arg3 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of6(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3) returns Set {
    handle externalObj = java_util_Set_of6(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + arg3 - The `javalang:Object` value required to map with the Java method parameter.
# + arg4 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of7(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4) returns Set {
    handle externalObj = java_util_Set_of7(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + arg3 - The `javalang:Object` value required to map with the Java method parameter.
# + arg4 - The `javalang:Object` value required to map with the Java method parameter.
# + arg5 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of8(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5) returns Set {
    handle externalObj = java_util_Set_of8(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Set`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + arg3 - The `javalang:Object` value required to map with the Java method parameter.
# + arg4 - The `javalang:Object` value required to map with the Java method parameter.
# + arg5 - The `javalang:Object` value required to map with the Java method parameter.
# + arg6 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Set` value returning from the Java mapping.
public function Set_of9(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6) returns Set {
    handle externalObj = java_util_Set_of9(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj);
    Set newObj = new (externalObj);
    return newObj;
}

isolated function java_util_Set_add(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "add",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_Set_addAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "addAll",
    'class: "java.util.Set",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_Set_clear(handle receiver) = @java:Method {
    name: "clear",
    'class: "java.util.Set",
    paramTypes: []
} external;

function java_util_Set_contains(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "contains",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_Set_containsAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "containsAll",
    'class: "java.util.Set",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_Set_copyOf(handle arg0) returns handle = @java:Method {
    name: "copyOf",
    'class: "java.util.Set",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_Set_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_Set_forEach(handle receiver, handle arg0) = @java:Method {
    name: "forEach",
    'class: "java.util.Set",
    paramTypes: ["java.util.function.Consumer"]
} external;

function java_util_Set_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "java.util.Set",
    paramTypes: []
} external;

function java_util_Set_isEmpty(handle receiver) returns boolean = @java:Method {
    name: "isEmpty",
    'class: "java.util.Set",
    paramTypes: []
} external;

function java_util_Set_iterator(handle receiver) returns handle = @java:Method {
    name: "iterator",
    'class: "java.util.Set",
    paramTypes: []
} external;

function java_util_Set_of() returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: []
} external;

function java_util_Set_of10(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Set_of11(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7, handle arg8) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Set_of12(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7, handle arg8, handle arg9) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Set_of2(handle arg0) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_Set_of3(handle arg0) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["[Ljava.lang.Object;"]
} external;

function java_util_Set_of4(handle arg0, handle arg1) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object", "java.lang.Object"]
} external;

function java_util_Set_of5(handle arg0, handle arg1, handle arg2) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Set_of6(handle arg0, handle arg1, handle arg2, handle arg3) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Set_of7(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Set_of8(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Set_of9(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Set_parallelStream(handle receiver) returns handle = @java:Method {
    name: "parallelStream",
    'class: "java.util.Set",
    paramTypes: []
} external;

function java_util_Set_remove(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "remove",
    'class: "java.util.Set",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_Set_removeAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "removeAll",
    'class: "java.util.Set",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_Set_removeIf(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "removeIf",
    'class: "java.util.Set",
    paramTypes: ["java.util.function.Predicate"]
} external;

function java_util_Set_retainAll(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "retainAll",
    'class: "java.util.Set",
    paramTypes: ["java.util.Collection"]
} external;

function java_util_Set_size(handle receiver) returns int = @java:Method {
    name: "size",
    'class: "java.util.Set",
    paramTypes: []
} external;

function java_util_Set_spliterator(handle receiver) returns handle = @java:Method {
    name: "spliterator",
    'class: "java.util.Set",
    paramTypes: []
} external;

function java_util_Set_stream(handle receiver) returns handle = @java:Method {
    name: "stream",
    'class: "java.util.Set",
    paramTypes: []
} external;

function java_util_Set_toArray(handle receiver) returns handle = @java:Method {
    name: "toArray",
    'class: "java.util.Set",
    paramTypes: []
} external;

function java_util_Set_toArray2(handle receiver, handle arg0) returns handle = @java:Method {
    name: "toArray",
    'class: "java.util.Set",
    paramTypes: ["java.util.function.IntFunction"]
} external;

function java_util_Set_toArray3(handle receiver, handle arg0) returns handle = @java:Method {
    name: "toArray",
    'class: "java.util.Set",
    paramTypes: ["[Ljava.lang.Object;"]
} external;

