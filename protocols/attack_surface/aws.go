package attack_surface

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
)

// 获取所有可用区域
func GetAllRegions() ([]string, error) {
	// 使用默认区域创建一个会话，用于获取区域列表
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"), // 使用一个默认区域
		Credentials: credentials.NewStaticCredentials(viper.GetString("AWS.AwsApiKey"), viper.GetString("AWS.AwsSecretKey"), ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	// 创建 EC2 服务客户端
	svc := ec2.New(sess)

	// 调用 DescribeRegions API
	result, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(true), // 包括已启用和未启用的区域
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe regions: %v", err)
	}

	// 提取区域名称
	var regions []string
	for _, region := range result.Regions {
		if region.RegionName != nil {
			regions = append(regions, *region.RegionName)
		}
	}

	return regions, nil
}

// 获取所有区域的EC2客户端
// func GetAllRegionEc2Clients() map[string]*ec2.EC2 {
// 	// 获取所有可用区域
// 	regions, err := GetAllRegions()
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	// 创建所有区域的客户端
// 	clients := make(map[string]*ec2.EC2)
// 	for _, region := range regions {
// 		sess, err := session.NewSession(&aws.Config{
// 			Region:      aws.String(region),
// 			Credentials: credentials.NewStaticCredentials(viper.GetString("AWS.AwsApiKey"), viper.GetString("AWS.AwsSecretKey"), ""),
// 		})
// 		if err != nil {
// 			log.Printf("Failed to create session for region %s: %v", region, err)
// 			continue
// 		}
// 		clients[region] = ec2.New(sess)
// 	}

// 	return clients
// }
