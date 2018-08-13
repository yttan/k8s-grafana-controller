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

func (c *GrafanaClient) PostDashboard(dashboardStr string, namespace string, grafanaIP string) {
	endpoint := "/api/dashboards/db"
	var requestBody = []byte(dashboardStr)
	url := "http://admin:admin@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to upload dashboard")
	}
}

func (c *GrafanaClient) GetAllDashboards(grafanaIP string) []byte {
	endpoint := "/api/search?type=dash-db&query=&starred=false"
	url := "http://admin:admin@" + c.grafanaIP + endpoint
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

func (c *GrafanaClient) GetDashboardByUID(uid string, grafanaIP string) map[string]interface{} {
	endpoint := "/api/dashboards/uid/" + uid
	url := "http://admin:admin@" + c.grafanaIP + endpoint
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
	url := "http://admin:admin@" + c.grafanaIP + endpoint
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
	url := "http://admin:admin@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to delete org")
	}
}

func (c *GrafanaClient) PostOrg(name string, grafanaIP string) {
	var requestBody = []byte(`{"name":"` + name + `"}`)
	endpoint := "/api/orgs"
	url := "http://admin:admin@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to post org")
	}
}

func (c *GrafanaClient) PostDataSource(grafanaIP string) {
	endpoint := "/api/datasources"
	prometheus := prometheusIP() + ":9090"
	var requestBody = []byte(`{"name":"prometheus","type":"prometheus","url":"http://` + prometheus + `","access":"proxy"}`)

	url := "http://admin:admin@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to add data source")
	}
}

func (c *GrafanaClient) PostUser(name string, grafanaIP string) {
	endpoint := "/api/admin/users"
	var requestBody = []byte(`{"name":"` + name + `","login":"` + name + `","password":"password"}`)
	url := "http://admin:admin@" + c.grafanaIP + endpoint
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
	url := "http://admin:admin@" + c.grafanaIP + endpoint
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
	url := "http://admin:admin@" + c.grafanaIP + endpoint
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

func (c *GrafanaClient) AdminSwitchOrg(orgID int, grafanaIP string) {
	endpoint := "/api/user/using/" + strconv.Itoa(orgID)
	url := "http://admin:admin@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to switch org")
	}
}

func (c *GrafanaClient) PostUserToOrg(name string, orgID int, grafanaIP string) {
	endpoint := "/api/orgs/" + strconv.Itoa(orgID) + "/users"
	var requestBody = []byte(`{"loginOrEmail":"` + name + `","role":"Viewer"}`)
	url := "http://admin:admin@" + c.grafanaIP + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		glog.Error(err)
	}
	status, respBody := tryRequest(req, 3)
	if !reqSuccess(status, respBody) {
		glog.Warningln("fail to add user to org")
	}
}

func (c *GrafanaClient) SwitchUserContext(userID int, orgID int, grafanaIP string) {
	endpoint := "/api/users/" + strconv.Itoa(userID) + "/using/" + strconv.Itoa(orgID)
	url := "http://admin:admin@" + c.grafanaIP + endpoint
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
	url := "http://admin:admin@" + c.grafanaIP + endpoint
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

func tryRequest(req *http.Request, num int) (string, []byte) {
	var status string
	var body []byte
	var err error
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
