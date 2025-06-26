package sealsuite

import (
	"github.com/colin-404/logx"
	"github.com/tidwall/gjson"
	"github.com/xid-protocol/xidp/common"
)

// getAllUsersForDepartments 为所有部门获取用户
func (ss *sealsuite) getAllUsersForDepartments(departments *[]Department) *[]gjson.Result {
	logx.Infof("=== 开始获取各部门用户 ===")
	allUsers := make([]gjson.Result, 0)
	for _, dept := range *departments {
		users := ss.getUsersByDepartment(dept.ID)
		if users != nil {
			// 遍历用户数组，添加每个用户而不是整个数组
			allUsers = append(allUsers, users.Array()...)
		}
	}

	return &allUsers
}

// getUsersByDepartment 根据部门ID获取用户列表
func (ss *sealsuite) getUsersByDepartment(departmentID string) *gjson.Result {
	//https://{endpoint}/api/open/v1/department/detail?id=od_rYKNJpXA5Ro1
	url := ss.endpoint + "/api/open/v1/user/list?department_id=" + departmentID
	headers := map[string]string{
		"Authorization": ss.token,
	}

	resp := common.DoHttp("GET", url, nil, headers)
	responseBody := resp.String()

	// 解析用户数据
	result := gjson.Parse(responseBody)
	if !result.Get("data").Exists() {
		logx.Warnf("部门 %s 获取用户数据失败或无数据", departmentID)
		return nil
	}
	if result.Get("data.user_list").Exists() {
		userList := result.Get("data.user_list")
		return &userList
	}

	return nil

}
