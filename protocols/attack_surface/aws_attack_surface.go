package attack_surface

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"regexp"

// 	"github.com/tidwall/gjson"
// 	"github.com/xid-protocol/xidp/db"
// )

// const ExposureIprange = "0.0.0.0/0"

// type Rule struct {
// 	FromPort int      `json:"fromPort"`
// 	ToPort   int      `json:"toPort"`
// 	Protocol string   `json:"protocol"`
// 	Iprange  []string `json:"iprange"`
// }

// // path /protocols/external-attack-surface/aws-instance
// type AWSAttackSurface struct {
// 	InstanceID   string            `json:"instanceId"`
// 	InstanceName string            `json:"instanceName"`
// 	Tags         map[string]string `json:"tags"`
// 	PublicIPs    []string          `json:"publicIps"`
// 	PrivateIPs   []string          `json:"privateIps"`
// 	Rules        []Rule            `json:"rules"`
// }

// func NewAWSAttackSurface(id string) *AWSAttackSurface {
// 	var awsEas AWSAttackSurface
// 	ctx := context.Background()
// 	repository := db.NewXidInfoRepository()
// 	secgroupXID, err := repository.FindOneByXidAndPath(ctx, id, "/info/aws/secgroup")
// 	if err != nil {
// 		return nil
// 	}

// 	// logx.Infof("secgroupXID: %v", secgroupXID)

// 	instanceXid, err := repository.FindOneByXidAndPath(ctx, id, "/info/aws/instance")
// 	if err != nil {
// 		return nil
// 	}
// 	jsonBytes, _ := json.Marshal(instanceXid.Payload)
// 	jsonStr := string(jsonBytes)

// 	// InstanceID
// 	awsEas.InstanceID = gjson.Get(jsonStr, "#(Key==instanceid).Value").String()

// 	// Tags
// 	awsEas.Tags = make(map[string]string)
// 	tagArr := gjson.Get(jsonStr, `#(Key=="tags").Value`).Array()
// 	for _, tagKVList := range tagArr { // 每个 tagKVList 仍是一个数组
// 		var k, v string
// 		for _, kv := range tagKVList.Array() { // kv 形如 {"Key":"key","Value":"org"}
// 			if kv.Get("Key").String() == "key" {
// 				k = kv.Get("Value").String()
// 			} else if kv.Get("Key").String() == "value" {
// 				v = kv.Get("Value").String()
// 			}
// 		}
// 		if k != "" { // 填入结果
// 			awsEas.Tags[k] = v
// 		}
// 	}

// 	// InstanceName
// 	awsEas.InstanceName = gjson.Get(jsonStr, `#(Key=="Name").Value`).String()

// 	// ---------- 正则提取公网 IP ----------
// 	// 匹配 "Key":"publicip" ... "Value":"xxx"
// 	re := regexp.MustCompile(`"Key"\s*:\s*"publicip"[\s\S]*?"Value"\s*:\s*"([^"]+)"`)
// 	matches := re.FindAllStringSubmatch(jsonStr, -1)
// 	seen := make(map[string]struct{})
// 	for _, m := range matches {
// 		if len(m) >= 2 {
// 			ip := m[1]
// 			if _, ok := seen[ip]; !ok {
// 				awsEas.PublicIPs = append(awsEas.PublicIPs, ip)
// 				seen[ip] = struct{}{}
// 			}
// 		}
// 	}

// 	// ---------- 正则提取私网 IP ----------
// 	// 匹配 "Key":"privateip" ... "Value":"xxx"
// 	re = regexp.MustCompile(`"Key"\s*:\s*"privateipaddress"[\s\S]*?"Value"\s*:\s*"([^"]+)"`)
// 	matches = re.FindAllStringSubmatch(jsonStr, -1)
// 	seen = make(map[string]struct{})
// 	for _, m := range matches {
// 		if len(m) >= 2 {
// 			ip := m[1]
// 			if _, ok := seen[ip]; !ok {
// 				awsEas.PrivateIPs = append(awsEas.PrivateIPs, ip)
// 				seen[ip] = struct{}{}
// 			}
// 		}
// 	}

// 	//---------------------secgroup---------------------

// 	secgroupJsonBytes, err := json.Marshal(secgroupXID.Payload)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}
// 	secgroupJsonStr := string(secgroupJsonBytes)
// 	// fmt.Println(secgroupJsonStr)

// 	// ----------- 提取所有安全组规则 -------------
// 	securitygroups := gjson.Get(secgroupJsonStr, `0.#(Key=="securitygroups").Value`).Array()
// 	fmt.Println(securitygroups)
// 	var allRules []Rule

// 	for _, sg := range securitygroups {
// 		// 每个 sg 仍是一个 Key-Value 数组，需要再次过滤
// 		ipPerms := sg.Get(`#(Key=="ippermissions").Value`).Array()
// 		for _, perm := range ipPerms {
// 			from := int(perm.Get(`#(Key=="fromport").Value`).Int())
// 			to := int(perm.Get(`#(Key=="toport").Value`).Int())
// 			proto := perm.Get(`#(Key=="ipprotocol").Value`).String()

// 			var cidrs []string
// 			// ipranges 可能为 null
// 			ipranges := perm.Get(`#(Key=="ipranges").Value`).Array()
// 			for _, iprange := range ipranges {
// 				iprangeStr := iprange.Get(`#(Key=="cidrip").Value`).String()
// 				cidrs = append(cidrs, iprangeStr)
// 			}
// 			allRules = append(allRules, Rule{
// 				FromPort: from,
// 				ToPort:   to,
// 				Protocol: proto,
// 				Iprange:  cidrs,
// 			})

// 		}
// 	}

// 	awsEas.Rules = allRules

// 	return &awsEas
// }

// func GetSgByID(securityGroupID *string, ec2Cli *ec2.EC2) *ec2.DescribeSecurityGroupsOutput {

// repository := repositories.NewXidInfoRepository()
// //实例化安全组
// var secGroup SecGroup

// securityGroupInput := &ec2.DescribeSecurityGroupsInput{
// 	GroupIds: []*string{securityGroupID},
// }
// securityGroupOutput, err := ec2Cli.DescribeSecurityGroups(securityGroupInput)
// // log.Println(securityGroupOutput)
// if err != nil {
// 	log.Println("Failed to describe security groups:", err)
// }

// return securityGroupOutput
// return nil
// }
