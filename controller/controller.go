package controller

import (
	"flag"
	"k8s-grafana-controller/grafana"
	"os"
	"path/filepath"
	"regexp"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Initiate clientset
func InitClientSet() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	configPath := pathToConfig()
	kubeconfig = flag.String("kubeconfig", filepath.Join(configPath, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func InitGrafanaClient() (*grafana.GrafanaClient, error) {
	grafanaClient, err := grafana.NewGrafanaClient(grafanaIP())
	if err != nil {
		return nil, err
	}
	return grafanaClient, nil
}

// WatchGrafana watches the grafana pod. If the pod is deleted, post existing tenants to ensure correctness
func WatchGrafana(clientset *kubernetes.Clientset, grafanaClient *grafana.GrafanaClient) {

	watchGrafana, err := clientset.CoreV1().Pods("monitoring").Watch(metav1.ListOptions{Watch: true})
	if err != nil {
		glog.Fatal(err)
	} else {
		dbList := grafanaClient.GetDashboardList()
		eventChan := watchGrafana.ResultChan()
		for event := range eventChan {
			pod, ok := event.Object.(*v1.Pod)
			if !ok {
				glog.Errorln("unexpected type when watching pods")
			} else {
				switch event.Type {
				case watch.Deleted:
					glog.Infoln("Pod " + pod.Name + " deleted")
					pattern := "^kube-prometheus-grafana[0-9a-z-]+"
					match, _ := regexp.Match(pattern, []byte(pod.Name))
					if match {
						namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
						if err != nil {
							glog.Error(err)
						}
						if len(namespaces.Items) == 0 {
							glog.Warning("No namespaces found")
						} else {
							for _, namespace := range namespaces.Items {
								grafanaClient.PostTenant(namespace.Name, dbList)
								glog.Infoln("namespace " + namespace.Name + " added")
							}
						}
					}
				case watch.Added:
					glog.Infoln("Pod " + pod.Name + " added")
				case watch.Error:
					glog.Infoln("Pod " + pod.Name + " has an error")
				}
			}
		}
	}
	glog.Flush()
}

// WatchTenants watches namespaces of K8S. If a new namespace is created, add tenant accordingly.
func WatchTenants(clientset *kubernetes.Clientset, grafanaClient *grafana.GrafanaClient) {
	var watchns watch.Interface
	watchns, err := clientset.CoreV1().Namespaces().Watch(metav1.ListOptions{Watch: true})
	if err != nil {
		glog.Fatal(err)
	} else {
		dbList := grafanaClient.GetDashboardList()
		eventChan := watchns.ResultChan()
		for event := range eventChan {
			ns, ok := event.Object.(*v1.Namespace)
			if !ok {
				glog.Errorln("unexpected type when watching namespaces")
			} else {
				switch event.Type {
				case watch.Added:
					grafanaClient.PostTenant(ns.Name, dbList)
					glog.Infoln("namespace " + ns.Name + " added")
				case watch.Deleted:
					grafanaClient.DeleteTenant(ns.Name)
					glog.Infoln("namespace " + ns.Name + " deleted")
				case watch.Error:
					glog.Infoln("namespace " + ns.Name + " has an error")
				}
			}
		}
	}
	glog.Flush()
}

func grafanaIP() string {
	ip := os.Getenv("GRAFANA_IP")
	return ip
}

func pathToConfig() string {
	path := os.Getenv("CONFIG_PATH")
	return path
}
