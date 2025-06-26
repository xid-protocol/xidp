package sealsuite

import (
	"github.com/colin-404/logx"
	"github.com/tidwall/gjson"
	"github.com/xid-protocol/xidp/common"
)

// Department 部门结构体
type Department struct {
	ID                     string       `json:"id"`
	Name                   string       `json:"name"`
	Type                   int          `json:"type"`
	ParentID               string       `json:"parent_id"`
	Seq                    int          `json:"seq"`
	ThirdPartyType         int          `json:"third_party_type,omitempty"`
	ThirdPartyDepartmentID string       `json:"third_party_department_id,omitempty"`
	Path                   string       `json:"path"`  // 完整路径，如 /LiquidityTech/Prime/Tech
	Level                  int          `json:"level"` // 层级深度
	SubDepartments         []Department `json:"sub_departments"`
}

func (ss *sealsuite) getDeparment() *[]Department {
	url := ss.endpoint + "/api/open/v1/department/list"
	headers := map[string]string{
		"Authorization": ss.token,
	}
	resp := common.DoHttp("GET", url, nil, headers)
	dataArray := gjson.Get(resp.String(), "data")

	//根部门
	rootDepartment := dataArray.Array()[0]
	logx.Infof("根部门: %v", rootDepartment.Get("name").String())

	//子部门
	subDepartments := rootDepartment.Get("sub_departments")
	logx.Infof("subDepartments: %v", subDepartments)

	if !subDepartments.Exists() {
		logx.Warnf("没有找到子部门")
		return nil
	}

	// 开始递归处理所有部门
	allDepartments := ss.getAllDepartmentsWithPath(resp.String())

	// 输出所有部门路径
	return &allDepartments
}

// getAllDepartmentsWithPath 递归获取所有部门并构建完整路径
func (ss *sealsuite) getAllDepartmentsWithPath(departmentData string) []Department {
	var allDepartments []Department

	result := gjson.Parse(departmentData)
	if !result.Get("data").Exists() {
		logx.Errorf("无效的部门数据格式")
		return allDepartments
	}

	// 获取data数组
	dataArray := result.Get("data").Array()
	if len(dataArray) == 0 {
		logx.Errorf("没有找到部门数据")
		return allDepartments
	}

	// 处理根部门（通常只有一个）
	for _, rootDeptData := range dataArray {
		rootDept := ss.parseDepartmentFromGJSON(rootDeptData, "", 0)
		allDepartments = append(allDepartments, rootDept)

		// 递归获取所有子部门
		childDepartments := ss.recursiveGetDepartments(rootDept.SubDepartments, rootDept.Path, 1)
		allDepartments = append(allDepartments, childDepartments...)
	}

	return allDepartments
}

// parseDepartmentFromGJSON 从gjson.Result解析部门数据
func (ss *sealsuite) parseDepartmentFromGJSON(deptData gjson.Result, parentPath string, level int) Department {
	dept := Department{
		ID:                     deptData.Get("id").String(),
		Name:                   deptData.Get("name").String(),
		Type:                   int(deptData.Get("type").Int()),
		ParentID:               deptData.Get("parent_id").String(),
		Seq:                    int(deptData.Get("seq").Int()),
		ThirdPartyType:         int(deptData.Get("third_party_type").Int()),
		ThirdPartyDepartmentID: deptData.Get("third_party_department_id").String(),
		Level:                  level,
	}

	// 构建路径
	if parentPath == "" {
		dept.Path = "/" + dept.Name
	} else {
		dept.Path = parentPath + "/" + dept.Name
	}

	// 解析子部门
	subDepartments := deptData.Get("sub_departments")
	if subDepartments.Exists() && subDepartments.IsArray() {
		for _, subDeptData := range subDepartments.Array() {
			subDept := ss.parseDepartmentFromGJSON(subDeptData, "", 0) // 路径在递归时会重新设置
			dept.SubDepartments = append(dept.SubDepartments, subDept)
		}
	}

	return dept
}

// recursiveGetDepartments 递归获取子部门
func (ss *sealsuite) recursiveGetDepartments(departments []Department, parentPath string, level int) []Department {
	var result []Department

	for _, dept := range departments {
		// 构建当前部门的完整路径
		dept.Path = parentPath + "/" + dept.Name
		dept.Level = level

		result = append(result, dept)

		// 递归处理子部门
		if len(dept.SubDepartments) > 0 {
			childDepartments := ss.recursiveGetDepartments(dept.SubDepartments, dept.Path, level+1)
			result = append(result, childDepartments...)
		}
	}

	return result
}
