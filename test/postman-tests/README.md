# WSO2 API Postman Based Integration Tests

This folder contains APK postman integration tests that use the test APK product capabilities

## Table of Contents

1. [Overview](#overview)
1. [Using Postman Tests](#using-postman-tests)
    1. [Load Environment into Postman](#adding-an-environment-to-postman)
    1. [Add Collection to Postman](#adding-collection-to-postman)
    1. [Run Tests](#run-tests)
1. [Modifying Tests](#modifying-tests)
1. [Diagnosing Failures](#diagnosing-failures)
1. [Notes](#notes)

## Overview

The goal of the framework is to make it as easy as possible to author and run tests. In the initial release of APK we won't release UI components. So all integrations with the product will happen via APIs. These tests allow you to test and validate each user operation APK support. Also users can use the same to create APIs, applications and other entries programmatically.


## Using Postman Tests

To use these postman tests users need to install postman in their local machines or integrate postman collections via Testkube. In this part of the document we will be discussing how to test product APIs with locally installed postman scripts.

### Adding an Environment to Postman

An environment is a set of variables you can use in your Postman requests. You can use environments to group related sets of values together and manage access to shared Postman data if you are working as part of a team.

* Select Environments from the sidebar.
* Click on the import button and then select file based import.
* Then select the `APK_environment.json` file available under resources/environments directory.
* Here initial values are configured to connect tests with locally installed minikube deployment. If your deployment is different then please change server URLs accordingly.


If you are interested please read more about 
[postman environments](https://learning.postman.com/docs/sending-requests/managing-environments/).


### Adding Collection to Postman

Postman `Collections` are a group of saved requests. Every request you send in Postman appears under the History tab of the sidebar. In APK test cases we have created multiple collections based on the user flows. Collections available under resources/collections directory.


* Select Collections from the sidebar.
* Click on import button and then select file based import. 
* Then select `<Test Name>-collection.json` file available under resources/collections directory.
* Then you can expand the collection and see requests and operations in the collection. These requests are structured based on the order of operations and you will need to run them sequentially.


### Run Tests
Select `Collections` in the sidebar and select the collection you want to run.

In APK tests you will need to run a collection with an environment as it needs some variables from the environment. You can select the environment using the environment selector at the top right of Postman. You can also select Environments in the sidebar, then select the environment you want to use.

On the overview tab, select `Runner` icon Run.

By default, your requests run in the sequence they're listed in the collection. If you need to change the order of execution, select and drag a request to its new location in the order. You can also remove an individual request from the run by clearing the checkbox next to its name.

To access more data about what happened when a request executed, do the following:
* Select `View Results`.
* Select the request's name from the list of results.

The `View Results` page indicates whether each request passed, failed, or was skipped. You can also filter on each using the Passed, Failed, and Skipped tabs at the top. If any tests in a request script fail during the collection run, the whole request fails.


## Modifying Tests
You can install Testkube on your command line to create, manage, and run your API tests or/and just access the Testkube UI, which runs from within your cluster and do it from there.


If you are interested please read more about 
[running APK postman tests using Testkube](https://testkube.io/blog/api-testing-in-kubernetes-with-postman).



## Diagnosing Failures
Each request sent out from API test scripts is logged in the console, so you can view the details of what happened when you sent a request. This means you can use the Postman Console to help debug your requests when an API isn't behaving as you expect. Keeping the console open while you work will increase the visibility of your network calls and log messages while debugging.


## Notes
Before you run postman test scripts please verify your APK services deployed properly in the given namespace/ cluster. Also sometimes you will need to install ingress or configure port forwarding to access domain services.