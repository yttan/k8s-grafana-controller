package grafana

import (
	"encoding/json"
	"errors"

	"github.com/golang/glog"
)

type GrafanaClient struct {
	grafanaIP string
	user      string
	password  string
}

// NewGrafanaClient creates a new client to control grafana pod
func NewGrafanaClient(grafanaIP string, user string, password string) (*GrafanaClient, error) {
	if grafanaIP == "" {
		return nil, errors.New("grafanaIP is empty string")
	}
	if user == "" {
		return nil, errors.New("user is empty string")
	}
	if password == "" {
		return nil, errors.New("password is empty string")
	}
	return &GrafanaClient{
		grafanaIP: grafanaIP,
		user:      user,
		password:  password,
	}, nil
}

// PostTenant posts a new tenant to grafana.
func (c *GrafanaClient) PostTenant(namespace string, dbList []map[string]interface{}) {
	c.PostOrg(namespace, c.grafanaIP)
	orgID := c.GetOrgID(namespace, c.grafanaIP)
	if orgID != 0 {
		c.SwitchOrg(orgID, c.grafanaIP)
		c.PostPrometheusDataSource(c.grafanaIP)
		for _, db := range dbList {
			dashboardStr := selectDashboard(db, namespace)
			if dashboardStr != "" {
				c.PostDashboard(dashboardStr, c.grafanaIP)
			}
		}
		c.PostUser(namespace, c.grafanaIP)
		c.PostUserToOrg(namespace, orgID, c.grafanaIP)
		userID := c.GetUserID(namespace, c.grafanaIP)
		c.SwitchUserContext(userID, orgID, c.grafanaIP)
		c.DeleteUserInOrg(userID, 1, c.grafanaIP)
	}
	glog.Flush()
}

// DeleteTenant deletes a tenant in Grafana
func (c *GrafanaClient) DeleteTenant(namespace string) {
	orgID := c.GetOrgID(namespace, c.grafanaIP)
	userID := c.GetUserID(namespace, c.grafanaIP)
	c.DeleteOrg(orgID, c.grafanaIP)
	c.DeleteUser(userID, c.grafanaIP)
	glog.Flush()
}

// GetDashboardList gets all the dashboards in an organization
func (c *GrafanaClient) GetDashboardList() []map[string]interface{} {
	var dbList []map[string]interface{}
	c.SwitchOrg(1, c.grafanaIP)
	allDbs := c.GetAllDashboards(c.grafanaIP)
	if allDbs == nil {
		return nil
	}
	dbList = c.processDashboardList(allDbs)
	glog.Flush()
	return dbList
}

// processDashboardList converts the dashboardlist string received to interface
func (c *GrafanaClient) processDashboardList(dsData []byte) []map[string]interface{} {
	var dbList []map[string]interface{}
	var ds interface{}
	err := json.Unmarshal(dsData, &ds)
	if err != nil {
		glog.Error(err)
		return nil
	}
	switch r := ds.(type) {
	case []interface{}:
		for _, d := range r {
			uid := d.(map[string]interface{})["uid"]
			uidStr := uid.(string)
			dashboard := c.GetDashboardByUID(uidStr, c.grafanaIP)
			if dashboard != nil {
				dbList = append(dbList, dashboard)
			}
		}
	}
	glog.Flush()
	return dbList
}

// select Deployment, Pods and StatefulSet, then modify the selected dashboards
func selectDashboard(dashboard map[string]interface{}, namespace string) string {
	switch dashboard["title"] {
	case "Deployment":
		dbstr := processDashboard(dashboard, namespace)
		return dbstr
	case "Pods":
		dbstr := processDashboard(dashboard, namespace)
		return dbstr
	case "StatefulSet":
		dbstr := processDashboard(dashboard, namespace)
		return dbstr
	default:
		return ""
	}
}

// modify dashboard before post them to grafana
func processDashboard(dashboard map[string]interface{}, namespace string) string {
	var nullString *string
	dashboard["id"] = nullString
	dashboard["uid"] = nullString
	dashboard["version"] = 0
	templates := dashboard["templating"].(map[string]interface{})
	temp := processTemplate(templates, namespace)
	dashboard["templating"] = temp
	db, err := json.Marshal(dashboard)
	if err != nil {
		glog.Warningln("unable to marshal dashboard")
	}
	var dbstr string
	dbstr = "{\"dashboard\":" + string(db) + ", \"overwrite\": false}"
	glog.Flush()
	return dbstr
}

// add namespace to dashboard template regex
func processTemplate(template map[string]interface{}, namespace string) map[string]interface{} {
	tempList := template["list"].([]interface{})
	for _, t := range tempList {
		allValue := t.(map[string]interface{})["allValue"]
		if allValue != nil {
			var nullString *string
			t.(map[string]interface{})["allValue"] = nullString
		}
		label := t.(map[string]interface{})["label"]
		if label != nil {
			labelStr := label.(string)
			if labelStr == "Namespace" {
				t.(map[string]interface{})["regex"] = namespace
				t.(map[string]interface{})["hide"] = 2
			}
		}
	}
	result := make(map[string]interface{})
	result["list"] = tempList
	return result
}
