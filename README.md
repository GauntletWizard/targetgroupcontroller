# Introducing TargetGroupController - A way to more efficiently run your kubernetes services

Here at Houseparty, we use AWS ALBs for Load Balancing, SSL Fronting, and more. We find them an effective and reliable tool to manage our traffic and do service routing. What we did not find them effective for was doing healthchecking.

## A primer on how Healthchecking works.

Healthchecking is the process of making requests to your service or backend to assure that it is still up. This can help with problem detection and remediation, and is often used to deal with services with a slow infancy - Instances are not added to the serving pool until they begin passing healthchecks. Kubernetes uses both - LivenessProbes to tell if pods have become unresponsive and automatically restart them, and ReadinessProbes to tell if pods are ready to serve traffic.

# TargetGroupController - Shim between Kubernetes Services and AWS ALB/ELBv2

TargetGroupController is a simple kubernetes styled controller to solve a simple issue. It translates between [Kubernetes Services](https://kubernetes.io/docs/concepts/services-networking/)/Endpoint objects and AWS TargetGroups. For a given kubernetes Service, it watches to rectify the list of IPs in the service ready endpoints with the list of IPs in the targetgroup.

## Problem statement

The built-in kubernetes Service controller can create Amazon ELBs. The routing diagram for these ELBs, however, presents a problem. The service creates a Nodeport, and the registers each node as a backend to the ELB. The healthchecks that the ELB performs check the health of the *nodes*, not of the underlying *pods*. There may be multiple connections to each node, each with a different backing pod. Since it is the nodes that are healthchecked, multiple connections might be considered 'healthy' even if the pods that underly them are not.

To rectify this, we need to improve routing so that the ALB is aware of the state of each individual pod. We do this by using an IP type TargetGroup. Caveat: This requires that the pods be routable within the VPC - You can't be using Calico or other network models that do not expose pod IPs externally.

## Using the TargetGroupController:

Create a ALB, and then create a targetgroup of type 'ip'. Note down it's ARN. Create an IAM user with the template in `aws/policy.json`, and put it's credentials into the secret `aws-targetgroupcontroller`. Then, create the sa, service, and deployment in the k8s directory, substituting in your service name and the targetgroup's ARN.

## Running the TargetGroupController:

TGC Implements a [Prometheus](https://prometheus.io/) monitoring endpoint. Example rules are in the prometheus directory.

# Building:

TargetGroupController is built and pushed by bazel - All that's needed is to run `bazel run :push_dockerimage  "--embed_label=$(git rev-parse HEAD)"`. See the BUILD file for other targets.
