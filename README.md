# k8s-grafana-controller

A grafana multi-tenancy management tool, written in Go with the [client-go](https://github.com/kubernetes/client-go) library.

With the controller, a grafana organization is bond with a Kubernetes tenant. Each organization has a viewer and the viewer is only added to that specific organization. A server admin controls all the organizations.

![](user.svg)

In one organization, there are Pod, Deployment and StatefulSet dashboards. For now, tenants are managed using kubernetes namespaces, which means  dashboards in an organization only shows data of the related namespace. Additionally, the default organization shows all the namespaces and has dashboards showing the cluster status. Only the server admin is in that organization. 
![namespace](namespace.svg)



## Usage  
Server admin access is needed to use the controller.
