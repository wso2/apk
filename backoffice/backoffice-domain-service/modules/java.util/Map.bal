import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import backoffice_domain_service.java.lang as javalang;
import backoffice_domain_service.java.util.'function as javautilfunction;

# Ballerina class mapping for the Java `java.util.Map` interface.
@java:Binding {'class: "java.util.Map"}
public distinct class Map {

    *java:JObject;

    # The `handle` field that stores the reference to the `java.util.Map` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `java.util.Map` Java interface.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `java.util.Map` Java interface.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `clear` method of `java.util.Map`.
    public function clear() {
        java_util_Map_clear(self.jObj);
    }

    # The function that maps to the `compute` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg1 - The `javautilfunction:BiFunction` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function compute(javalang:Object arg0, javautilfunction:BiFunction arg1) returns javalang:Object {
        handle externalObj = java_util_Map_compute(self.jObj, arg0.jObj, arg1.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `computeIfAbsent` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg1 - The `javautilfunction:Function` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function computeIfAbsent(javalang:Object arg0, javautilfunction:Function arg1) returns javalang:Object {
        handle externalObj = java_util_Map_computeIfAbsent(self.jObj, arg0.jObj, arg1.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `computeIfPresent` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg1 - The `javautilfunction:BiFunction` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function computeIfPresent(javalang:Object arg0, javautilfunction:BiFunction arg1) returns javalang:Object {
        handle externalObj = java_util_Map_computeIfPresent(self.jObj, arg0.jObj, arg1.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `containsKey` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function containsKey(javalang:Object arg0) returns boolean {
        return java_util_Map_containsKey(self.jObj, arg0.jObj);
    }

    # The function that maps to the `containsValue` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function containsValue(javalang:Object arg0) returns boolean {
        return java_util_Map_containsValue(self.jObj, arg0.jObj);
    }

    # The function that maps to the `entrySet` method of `java.util.Map`.
    #
    # + return - The `Set` value returning from the Java mapping.
    public function entrySet() returns Set {
        handle externalObj = java_util_Map_entrySet(self.jObj);
        Set newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `equals` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return java_util_Map_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `forEach` method of `java.util.Map`.
    #
    # + arg0 - The `javautilfunction:BiConsumer` value required to map with the Java method parameter.
    public function forEach(javautilfunction:BiConsumer arg0) {
        java_util_Map_forEach(self.jObj, arg0.jObj);
    }

    # The function that maps to the `get` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function get(javalang:Object arg0) returns javalang:Object {
        handle externalObj = java_util_Map_get(self.jObj, arg0.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getOrDefault` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg1 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function getOrDefault(javalang:Object arg0, javalang:Object arg1) returns javalang:Object {
        handle externalObj = java_util_Map_getOrDefault(self.jObj, arg0.jObj, arg1.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `java.util.Map`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return java_util_Map_hashCode(self.jObj);
    }

    # The function that maps to the `isEmpty` method of `java.util.Map`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public function isEmpty() returns boolean {
        return java_util_Map_isEmpty(self.jObj);
    }

    # The function that maps to the `keySet` method of `java.util.Map`.
    #
    # + return - The `Set` value returning from the Java mapping.
    public function keySet() returns Set {
        handle externalObj = java_util_Map_keySet(self.jObj);
        Set newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `merge` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg1 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg2 - The `javautilfunction:BiFunction` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function merge(javalang:Object arg0, javalang:Object arg1, javautilfunction:BiFunction arg2) returns javalang:Object {
        handle externalObj = java_util_Map_merge(self.jObj, arg0.jObj, arg1.jObj, arg2.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `put` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg1 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function put(javalang:Object arg0, javalang:Object arg1) returns javalang:Object {
        handle externalObj = java_util_Map_put(self.jObj, arg0.jObj, arg1.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `putAll` method of `java.util.Map`.
    #
    # + arg0 - The `Map` value required to map with the Java method parameter.
    public function putAll(Map arg0) {
        java_util_Map_putAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `putIfAbsent` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg1 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function putIfAbsent(javalang:Object arg0, javalang:Object arg1) returns javalang:Object {
        handle externalObj = java_util_Map_putIfAbsent(self.jObj, arg0.jObj, arg1.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `remove` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function remove(javalang:Object arg0) returns javalang:Object {
        handle externalObj = java_util_Map_remove(self.jObj, arg0.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `remove` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg1 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function remove2(javalang:Object arg0, javalang:Object arg1) returns boolean {
        return java_util_Map_remove2(self.jObj, arg0.jObj, arg1.jObj);
    }

    # The function that maps to the `replace` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg1 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `javalang:Object` value returning from the Java mapping.
    public function replace(javalang:Object arg0, javalang:Object arg1) returns javalang:Object {
        handle externalObj = java_util_Map_replace(self.jObj, arg0.jObj, arg1.jObj);
        javalang:Object newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `replace` method of `java.util.Map`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg1 - The `javalang:Object` value required to map with the Java method parameter.
    # + arg2 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function replace2(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2) returns boolean {
        return java_util_Map_replace2(self.jObj, arg0.jObj, arg1.jObj, arg2.jObj);
    }

    # The function that maps to the `replaceAll` method of `java.util.Map`.
    #
    # + arg0 - The `javautilfunction:BiFunction` value required to map with the Java method parameter.
    public function replaceAll(javautilfunction:BiFunction arg0) {
        java_util_Map_replaceAll(self.jObj, arg0.jObj);
    }

    # The function that maps to the `size` method of `java.util.Map`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function size() returns int {
        return java_util_Map_size(self.jObj);
    }

    # The function that maps to the `values` method of `java.util.Map`.
    #
    # + return - The `Collection` value returning from the Java mapping.
    public function values() returns Collection {
        handle externalObj = java_util_Map_values(self.jObj);
        Collection newObj = new (externalObj);
        return newObj;
    }

}

# The function that maps to the `copyOf` method of `java.util.Map`.
#
# + arg0 - The `Map` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_copyOf(Map arg0) returns Map {
    handle externalObj = java_util_Map_copyOf(arg0.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `entry` method of `java.util.Map`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Entry` value returning from the Java mapping.
public function Map_entry(javalang:Object arg0, javalang:Object arg1) returns Entry {
    handle externalObj = java_util_Map_entry(arg0.jObj, arg1.jObj);
    Entry newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
#
# + return - The `Map` value returning from the Java mapping.
public function Map_of() returns Map {
    handle externalObj = java_util_Map_of();
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
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
# + arg10 - The `javalang:Object` value required to map with the Java method parameter.
# + arg11 - The `javalang:Object` value required to map with the Java method parameter.
# + arg12 - The `javalang:Object` value required to map with the Java method parameter.
# + arg13 - The `javalang:Object` value required to map with the Java method parameter.
# + arg14 - The `javalang:Object` value required to map with the Java method parameter.
# + arg15 - The `javalang:Object` value required to map with the Java method parameter.
# + arg16 - The `javalang:Object` value required to map with the Java method parameter.
# + arg17 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_of10(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6, javalang:Object arg7, javalang:Object arg8, javalang:Object arg9, javalang:Object arg10, javalang:Object arg11, javalang:Object arg12, javalang:Object arg13, javalang:Object arg14, javalang:Object arg15, javalang:Object arg16, javalang:Object arg17) returns Map {
    handle externalObj = java_util_Map_of10(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj, arg7.jObj, arg8.jObj, arg9.jObj, arg10.jObj, arg11.jObj, arg12.jObj, arg13.jObj, arg14.jObj, arg15.jObj, arg16.jObj, arg17.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
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
# + arg10 - The `javalang:Object` value required to map with the Java method parameter.
# + arg11 - The `javalang:Object` value required to map with the Java method parameter.
# + arg12 - The `javalang:Object` value required to map with the Java method parameter.
# + arg13 - The `javalang:Object` value required to map with the Java method parameter.
# + arg14 - The `javalang:Object` value required to map with the Java method parameter.
# + arg15 - The `javalang:Object` value required to map with the Java method parameter.
# + arg16 - The `javalang:Object` value required to map with the Java method parameter.
# + arg17 - The `javalang:Object` value required to map with the Java method parameter.
# + arg18 - The `javalang:Object` value required to map with the Java method parameter.
# + arg19 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_of11(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6, javalang:Object arg7, javalang:Object arg8, javalang:Object arg9, javalang:Object arg10, javalang:Object arg11, javalang:Object arg12, javalang:Object arg13, javalang:Object arg14, javalang:Object arg15, javalang:Object arg16, javalang:Object arg17, javalang:Object arg18, javalang:Object arg19) returns Map {
    handle externalObj = java_util_Map_of11(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj, arg7.jObj, arg8.jObj, arg9.jObj, arg10.jObj, arg11.jObj, arg12.jObj, arg13.jObj, arg14.jObj, arg15.jObj, arg16.jObj, arg17.jObj, arg18.jObj, arg19.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_of2(javalang:Object arg0, javalang:Object arg1) returns Map {
    handle externalObj = java_util_Map_of2(arg0.jObj, arg1.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + arg3 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_of3(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3) returns Map {
    handle externalObj = java_util_Map_of3(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + arg3 - The `javalang:Object` value required to map with the Java method parameter.
# + arg4 - The `javalang:Object` value required to map with the Java method parameter.
# + arg5 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_of4(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5) returns Map {
    handle externalObj = java_util_Map_of4(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
#
# + arg0 - The `javalang:Object` value required to map with the Java method parameter.
# + arg1 - The `javalang:Object` value required to map with the Java method parameter.
# + arg2 - The `javalang:Object` value required to map with the Java method parameter.
# + arg3 - The `javalang:Object` value required to map with the Java method parameter.
# + arg4 - The `javalang:Object` value required to map with the Java method parameter.
# + arg5 - The `javalang:Object` value required to map with the Java method parameter.
# + arg6 - The `javalang:Object` value required to map with the Java method parameter.
# + arg7 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_of5(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6, javalang:Object arg7) returns Map {
    handle externalObj = java_util_Map_of5(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj, arg7.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
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
# + return - The `Map` value returning from the Java mapping.
public function Map_of6(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6, javalang:Object arg7, javalang:Object arg8, javalang:Object arg9) returns Map {
    handle externalObj = java_util_Map_of6(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj, arg7.jObj, arg8.jObj, arg9.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
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
# + arg10 - The `javalang:Object` value required to map with the Java method parameter.
# + arg11 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_of7(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6, javalang:Object arg7, javalang:Object arg8, javalang:Object arg9, javalang:Object arg10, javalang:Object arg11) returns Map {
    handle externalObj = java_util_Map_of7(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj, arg7.jObj, arg8.jObj, arg9.jObj, arg10.jObj, arg11.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
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
# + arg10 - The `javalang:Object` value required to map with the Java method parameter.
# + arg11 - The `javalang:Object` value required to map with the Java method parameter.
# + arg12 - The `javalang:Object` value required to map with the Java method parameter.
# + arg13 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_of8(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6, javalang:Object arg7, javalang:Object arg8, javalang:Object arg9, javalang:Object arg10, javalang:Object arg11, javalang:Object arg12, javalang:Object arg13) returns Map {
    handle externalObj = java_util_Map_of8(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj, arg7.jObj, arg8.jObj, arg9.jObj, arg10.jObj, arg11.jObj, arg12.jObj, arg13.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `of` method of `java.util.Map`.
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
# + arg10 - The `javalang:Object` value required to map with the Java method parameter.
# + arg11 - The `javalang:Object` value required to map with the Java method parameter.
# + arg12 - The `javalang:Object` value required to map with the Java method parameter.
# + arg13 - The `javalang:Object` value required to map with the Java method parameter.
# + arg14 - The `javalang:Object` value required to map with the Java method parameter.
# + arg15 - The `javalang:Object` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_of9(javalang:Object arg0, javalang:Object arg1, javalang:Object arg2, javalang:Object arg3, javalang:Object arg4, javalang:Object arg5, javalang:Object arg6, javalang:Object arg7, javalang:Object arg8, javalang:Object arg9, javalang:Object arg10, javalang:Object arg11, javalang:Object arg12, javalang:Object arg13, javalang:Object arg14, javalang:Object arg15) returns Map {
    handle externalObj = java_util_Map_of9(arg0.jObj, arg1.jObj, arg2.jObj, arg3.jObj, arg4.jObj, arg5.jObj, arg6.jObj, arg7.jObj, arg8.jObj, arg9.jObj, arg10.jObj, arg11.jObj, arg12.jObj, arg13.jObj, arg14.jObj, arg15.jObj);
    Map newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `ofEntries` method of `java.util.Map`.
#
# + arg0 - The `Entry[]` value required to map with the Java method parameter.
# + return - The `Map` value returning from the Java mapping.
public function Map_ofEntries(Entry[] arg0) returns Map|error {
    handle externalObj = java_util_Map_ofEntries(check jarrays:toHandle(arg0, "java.util.Map$Entry"));
    Map newObj = new (externalObj);
    return newObj;
}

function java_util_Map_clear(handle receiver) = @java:Method {
    name: "clear",
    'class: "java.util.Map",
    paramTypes: []
} external;

function java_util_Map_compute(handle receiver, handle arg0, handle arg1) returns handle = @java:Method {
    name: "compute",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.util.function.BiFunction"]
} external;

function java_util_Map_computeIfAbsent(handle receiver, handle arg0, handle arg1) returns handle = @java:Method {
    name: "computeIfAbsent",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.util.function.Function"]
} external;

function java_util_Map_computeIfPresent(handle receiver, handle arg0, handle arg1) returns handle = @java:Method {
    name: "computeIfPresent",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.util.function.BiFunction"]
} external;

function java_util_Map_containsKey(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "containsKey",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_Map_containsValue(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "containsValue",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_Map_copyOf(handle arg0) returns handle = @java:Method {
    name: "copyOf",
    'class: "java.util.Map",
    paramTypes: ["java.util.Map"]
} external;

function java_util_Map_entry(handle arg0, handle arg1) returns handle = @java:Method {
    name: "entry",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_entrySet(handle receiver) returns handle = @java:Method {
    name: "entrySet",
    'class: "java.util.Map",
    paramTypes: []
} external;

function java_util_Map_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_Map_forEach(handle receiver, handle arg0) = @java:Method {
    name: "forEach",
    'class: "java.util.Map",
    paramTypes: ["java.util.function.BiConsumer"]
} external;

function java_util_Map_get(handle receiver, handle arg0) returns handle = @java:Method {
    name: "get",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_Map_getOrDefault(handle receiver, handle arg0, handle arg1) returns handle = @java:Method {
    name: "getOrDefault",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "java.util.Map",
    paramTypes: []
} external;

function java_util_Map_isEmpty(handle receiver) returns boolean = @java:Method {
    name: "isEmpty",
    'class: "java.util.Map",
    paramTypes: []
} external;

function java_util_Map_keySet(handle receiver) returns handle = @java:Method {
    name: "keySet",
    'class: "java.util.Map",
    paramTypes: []
} external;

function java_util_Map_merge(handle receiver, handle arg0, handle arg1, handle arg2) returns handle = @java:Method {
    name: "merge",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.util.function.BiFunction"]
} external;

function java_util_Map_of() returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: []
} external;

function java_util_Map_of10(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7, handle arg8, handle arg9, handle arg10, handle arg11, handle arg12, handle arg13, handle arg14, handle arg15, handle arg16, handle arg17) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_of11(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7, handle arg8, handle arg9, handle arg10, handle arg11, handle arg12, handle arg13, handle arg14, handle arg15, handle arg16, handle arg17, handle arg18, handle arg19) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_of2(handle arg0, handle arg1) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_of3(handle arg0, handle arg1, handle arg2, handle arg3) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_of4(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_of5(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_of6(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7, handle arg8, handle arg9) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_of7(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7, handle arg8, handle arg9, handle arg10, handle arg11) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_of8(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7, handle arg8, handle arg9, handle arg10, handle arg11, handle arg12, handle arg13) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_of9(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7, handle arg8, handle arg9, handle arg10, handle arg11, handle arg12, handle arg13, handle arg14, handle arg15) returns handle = @java:Method {
    name: "of",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_ofEntries(handle arg0) returns handle = @java:Method {
    name: "ofEntries",
    'class: "java.util.Map",
    paramTypes: ["[Ljava.util.Map$Entry;"]
} external;

function java_util_Map_put(handle receiver, handle arg0, handle arg1) returns handle = @java:Method {
    name: "put",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_putAll(handle receiver, handle arg0) = @java:Method {
    name: "putAll",
    'class: "java.util.Map",
    paramTypes: ["java.util.Map"]
} external;

function java_util_Map_putIfAbsent(handle receiver, handle arg0, handle arg1) returns handle = @java:Method {
    name: "putIfAbsent",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_remove(handle receiver, handle arg0) returns handle = @java:Method {
    name: "remove",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object"]
} external;

function java_util_Map_remove2(handle receiver, handle arg0, handle arg1) returns boolean = @java:Method {
    name: "remove",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_replace(handle receiver, handle arg0, handle arg1) returns handle = @java:Method {
    name: "replace",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_replace2(handle receiver, handle arg0, handle arg1, handle arg2) returns boolean = @java:Method {
    name: "replace",
    'class: "java.util.Map",
    paramTypes: ["java.lang.Object", "java.lang.Object", "java.lang.Object"]
} external;

function java_util_Map_replaceAll(handle receiver, handle arg0) = @java:Method {
    name: "replaceAll",
    'class: "java.util.Map",
    paramTypes: ["java.util.function.BiFunction"]
} external;

function java_util_Map_size(handle receiver) returns int = @java:Method {
    name: "size",
    'class: "java.util.Map",
    paramTypes: []
} external;

function java_util_Map_values(handle receiver) returns handle = @java:Method {
    name: "values",
    'class: "java.util.Map",
    paramTypes: []
} external;

