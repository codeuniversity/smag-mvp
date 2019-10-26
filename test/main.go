package main

import (
	"fmt"
	"strconv"
)

func main() {
	//grpcClient, err := grpc.Dial("35.156.92.28:9900", grpc.WithInsecure())
	//if err != nil {
	//	panic(err)
	//}
	//
	//renewIp := pb.RenewingElasticIp{
	//	InstanceId: "i-0b09844b5b5344654",
	//	Node:       "",
	//	Pod:        "",
	//	PodIp:      "172.31.44.135",
	//}
	//
	//awsClient := pb.NewRouteGuideClient(grpcClient)
	//result, err := awsClient.RenewElasticIp(context.Background(), &renewIp)
	//if err != nil {
	//	fmt.Println("sendRenewElasticIpRequestToAmazonService Err: ", err)
	//	return
	//}
	//
	//fmt.Println(result.IsRenewed)

	test := 23

	f := strconv.FormatInt(int64(test), 10)
	fmt.Println(f)
}
