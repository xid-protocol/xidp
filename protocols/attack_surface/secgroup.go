package attack_surface

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/colin-404/logx"
	"github.com/spf13/viper"
)

// type SecGroup struct {
// 	GroupId          string `json:"groupId"`
// 	GroupName        string `json:"groupName"`
// 	GroupDescription string `json:"groupDescription"`
// 	GroupVpcId       string `json:"groupVpcId"`
// 	GroupOwnerId     string `json:"groupOwnerId"`
// 	GroupStatus      string `json:"groupStatus"`
// 	GroupTags        []Tag  `json:"groupTags"`
// }

func RunSecGroup() {
	//安全组
	// ec2SecGroup := aws_assets.SecGroupMonitor()
	// aws_assets.SetAWSSecGroupInfo(ec2SecGroup)
	// log.Println(ec2SecGroup)
}

// 通过security group ID 获取安全组规则
func GetSgByID(securityGroupID *string, ec2Cli *ec2.EC2) *ec2.DescribeSecurityGroupsOutput {

	// //实例化安全组
	// var secGroup SecGroup

	securityGroupInput := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{securityGroupID},
	}
	securityGroupOutput, err := ec2Cli.DescribeSecurityGroups(securityGroupInput)
	// log.Println(securityGroupOutput)
	if err != nil {
		logx.Errorf("Failed to describe security groups: %v", err)
	}

	return securityGroupOutput
}

func SecGroupMonitor() *map[string][]*ec2.DescribeSecurityGroupsOutput {
	// 获取所有区域的客户端
	clients := GetAllRegionEc2Clients()
	// 遍历每个区域获取实例信息
	ec2SecGroup := make(map[string][]*ec2.DescribeSecurityGroupsOutput)
	for region, client := range clients {
		result, err := client.DescribeInstances(nil)
		if err != nil {
			logx.Errorf("no instances in this region %s: %v", region, err)
			continue
		}

		// 处理该区域的实例
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				securityGroups := instance.SecurityGroups
				for _, securityGroup := range securityGroups {
					secGroup := GetSgByID(securityGroup.GroupId, client)
					ec2SecGroup[*instance.InstanceId] = append(ec2SecGroup[*instance.InstanceId], secGroup)
					// log.Println(secGroup)
				}
				// instances[*instance.InstanceId] = instance
			}
		}
	}
	return &ec2SecGroup
}

func GetAllRegionEc2Clients() map[string]*ec2.EC2 {
	// 获取所有可用区域
	regions, err := GetAllRegions()
	if err != nil {
		logx.Errorf("GetAllRegionEc2Clients: %v", err)
	}

	// 创建所有区域的客户端
	clients := make(map[string]*ec2.EC2)
	for _, region := range regions {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(viper.GetString("AWS.AwsApiKey"), viper.GetString("AWS.AwsSecretKey"), ""),
		})
		if err != nil {
			logx.Errorf("Failed to create session for region %s: %v", region, err)
			continue
		}
		clients[region] = ec2.New(sess)
	}

	return clients
}
