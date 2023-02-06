# DefaultApi

All URIs are relative to *http://api.example.com/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**usersGet**](DefaultApi.md#usersGet) | **GET** /users | Returns a list of users. |


<a name="usersGet"></a>
# **usersGet**
> List&lt;String&gt; usersGet()

Returns a list of users.

Optional extended description in CommonMark or HTML.

### Example
```java
// Import classes:
import org.wso2.client.api.ApiClient;
import org.wso2.client.api.ApiException;
import org.wso2.client.api.Configuration;
import org.wso2.client.api.models.*;
import org.wso2.client.api.PizzaShask.DefaultApi;

public class Example {
  public static void main(String[] args) {
    ApiClient defaultClient = Configuration.getDefaultApiClient();
    defaultClient.setBasePath("http://api.example.com/v1");

    DefaultApi apiInstance = new DefaultApi(defaultClient);
    try {
      List<String> result = apiInstance.usersGet();
      System.out.println(result);
    } catch (ApiException e) {
      System.err.println("Exception when calling DefaultApi#usersGet");
      System.err.println("Status code: " + e.getCode());
      System.err.println("Reason: " + e.getResponseBody());
      System.err.println("Response headers: " + e.getResponseHeaders());
      e.printStackTrace();
    }
  }
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

**List&lt;String&gt;**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A JSON array of user names |  -  |

