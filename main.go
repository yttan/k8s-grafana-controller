package main

import (
	"k8s-grafana-controller/controller"

	"github.com/golang/glog"
)

func main() {
	clientset, err := controller.InitClientSet()
	if err != nil {
		glog.Fatal(err)
	}
	grafanaClient, err := controller.InitGrafanaClient()
	if err != nil {
		glog.Fatal(err)
	}
	controllerClient, err := controller.InitControllerClient(grafanaClient)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Flush()
	go controller.WatchTenants(clientset, controllerClient)
	go controller.WatchGrafana(clientset, controllerClient)
	select {}
}
