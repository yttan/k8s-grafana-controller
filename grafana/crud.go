package grafana

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/golang/glog"
)

// PostDashboard posts a dashboard to grafana
func (c *GrafanaClient) PostDashboard(dashboardStr string, grafanaIP string) {
	endpoint := "/api/dashboards/db"
	var requestBody = []byte(dashboardStr)
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to upload dashboard")
	}
}

// GetAllDashboards gets all the dashboards in an organization. The received json does not include dashboard json data. To get dashboard json, use the uids you got and call GetDashboardByUID.
func (c *GrafanaClient) GetAllDashboards(grafanaIP string) []byte {
	endpoint := "/api/search?type=dash-db&query=&starred=false"
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to get all dashboards")
		return nil
	}
	return respBody
}

// GetDashboardByUID gets the dashboard json from grafana and marshals the received string
func (c *GrafanaClient) GetDashboardByUID(uid string, grafanaIP string) map[string]interface{} {
	endpoint := "/api/dashboards/uid/" + uid
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to get dashboard by uid " + uid)
		return nil
	}
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		glog.Error(err)
		return nil
	}
	dashboard := result["dashboard"].(map[string]interface{})
	return dashboard
}

func (c *GrafanaClient) DeleteUser(userID int, grafanaIP string) {
	endpoint := "/api/admin/users/" + strconv.Itoa(userID)
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("unable to delete user")
	}
}

func (c *GrafanaClient) DeleteOrg(orgID int, grafanaIP string) {
	endpoint := "/api/orgs/" + strconv.Itoa(orgID)
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to delete org")
	}
}

// PostOrg adds a new organization
func (c *GrafanaClient) PostOrg(name string, grafanaIP string) {
	var requestBody = []byte(`{"name":"` + name + `"}`)
	endpoint := "/api/orgs"
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to post org")
	}
}

// PostPrometheusDataSource adds a prometheus data source to the current organization. The prometheus ip address is read through environment variable PROMETHEUS_IP. If you have a DNS, you can change the url in request body to "http://{your prometheus service name}:9090" and delete the environment variable.
func (c *GrafanaClient) PostPrometheusDataSource(grafanaIP string) {
	endpoint := "/api/datasources"
	prometheus := prometheusIP() + ":9090"
	var requestBody = []byte(`{"name":"prometheus","type":"prometheus","url":"http://` + prometheus + `","access":"proxy"}`)

	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to add data source")
	}
}

// PostUser adds a new user
func (c *GrafanaClient) PostUser(name string, grafanaIP string) {
	endpoint := "/api/admin/users"
	var requestBody = []byte(`{"name":"` + name + `","login":"` + name + `","password":"password"}`)
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to add user")
	}
}

func (c *GrafanaClient) GetUserID(name string, grafanaIP string) int {
	var id int = 0
	endpoint := "/api/users/lookup?loginOrEmail=" + name
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, body := tryRequest(req, 3)
	if !reqSuccess(status, body) {
		glog.Warningln("fail to get user id")
	} else {
		m := make(map[string]int)
		_ = json.Unmarshal(body, &m)
		id = m["id"]
	}
	return id
}

func (c *GrafanaClient) GetOrgID(name string, grafanaIP string) int {
	var id int = 0
	endpoint := "/api/orgs/name/" + name
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, body := tryRequest(req, 3)
	if !reqSuccess(status, body) {
		glog.Warningln("fail to get org id")
	} else {
		m := make(map[string]int)
		_ = json.Unmarshal(body, &m)
		id = m["id"]
	}
	return id
}

// SwitchOrg changes the current organization of the client user.
func (c *GrafanaClient) SwitchOrg(orgID int, grafanaIP string) {
	endpoint := "/api/user/using/" + strconv.Itoa(orgID)
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to switch org")
	}
}

// PostUserToOrg adds a user to an organization.
func (c *GrafanaClient) PostUserToOrg(name string, orgID int, grafanaIP string) {
	endpoint := "/api/orgs/" + strconv.Itoa(orgID) + "/users"
	var requestBody = []byte(`{"loginOrEmail":"` + name + `","role":"Viewer"}`)
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to add user to org")
	}
}

// SwitchUserContext changes the current organization of a user. To call this, the client user must be server admin.
func (c *GrafanaClient) SwitchUserContext(userID int, orgID int, grafanaIP string) {
	endpoint := "/api/users/" + strconv.Itoa(userID) + "/using/" + strconv.Itoa(orgID)
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to switch user context")
	}
}

func (c *GrafanaClient) DeleteUserInOrg(userID int, orgID int, grafanaIP string) {
	endpoint := "/api/orgs/" + strconv.Itoa(orgID) + "/users/" + strconv.Itoa(userID)
	url := "http://" + c.user + ":" + c.password + "@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to delete user in org")
	}
}

func reqSuccess(status string, reqBody []byte) bool {
	if status != "200 OK" {
		glog.Warningln(string(reqBody))
		return false
	}
	return true
}

// try to send request. If the request fails, retry the request until the total number of requests meets the given number. num must be more than 1.
func tryRequest(req *http.Request, num int) (string, []byte) {
	var status string
	var body []byte
	var err error
	if num < 1 {
		glog.Error("Given request num less than 1")
		return "", []byte(err.Error())
	}
	for i := 1; i <= num; i++ {
		status, body, err = request(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", []byte(err.Error())
	}
	return status, body
}

func request(req *http.Request) (string, []byte, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return "", nil, err
	}
	status := resp.Status
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
	}
	resp.Body.Close()
	return status, body, nil
}

func prometheusIP() string {
	ip := os.Getenv("PROMETHEUS_IP")
	return ip
}
