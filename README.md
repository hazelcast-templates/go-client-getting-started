# Go Client Getting Started Template

This template creates a Go project that contains the code for the Hazelcast Go Client Getting Started guide.

## Before you Begin

* [Hazelcast Viridian Cloud Account](https://hazelcast.com/products/viridian/)
* [Hazelcast CLC](https://github.com/hazelcast/hazelcast-commandline-client/releases) 

### Start a Hazelcast Viridian Cloud Cluster

1. Sign up for a Hazelcast Viridian Cloud account (free trial is available).
2. Log in to your Hazelcast Viridian Cloud account and start your trial by filling in the welcome questionnaire.
3. A Viridian cluster will be created automatically when you start your trial.
4. Press the Connect Cluster dialog and switch over to the CLC tab for information about how to import the cluster configuration with Hazelcast CLC. 

### Create the Project

```
clc -c YOUR-CLUSTER-NAME project create python-client-getting-started
```