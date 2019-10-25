package main

import (
	"github.com/codeuniversity/smag-mvp/awsService"
	"github.com/codeuniversity/smag-mvp/utils"
)

func main() {

	kafkaAddress := utils.GetStringFromEnvWithDefault("KAFKA_ADDRESS", "127.0.0.1:9092")
	//
	s := awsService.New(kafkaAddress, 9900)
	s.Listen()

	//renewingElastic := pb.RenewingElasticIp{PodIp: "192.168.84.217", InstanceId: "i-0bc755555f0aaa272"}
	//result, err := s.RenewElasticIp(nil, &renewingElastic)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(result)
	//

	//s.ErrorTest()
}
