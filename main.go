package main

import "net/http"

func main(){
	// ヘルスチェック用のエンドポイント
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	// 8088万ポートでサーバーの起動
	http.ListenAndServe(":8080", nil)

}

// import (
//         "fmt"
//         "os"

//         "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
//         "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
//         "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
//         "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/regions"
//         cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
// 		trtc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/trtc/v20190722"
// )

// func main() {
//         // Essential steps:
//         // Instantiate an authentication object. The Tencent Cloud account key pair `secretId` and `secretKey` need to be passed in as the input parameters.
//         // The example here uses the way to read from the environment variable, so you need to set these two values in the environment variable first.
//         // You can also write the key pair directly into the code, but be careful not to copy, upload, or share the code to others;
//         // otherwise, the key pair may be leaked, causing damage to your properties.
//         credential := common.NewCredential(
//                 os.Getenv("TENCENTCLOUD_SECRET_ID"),
//                 os.Getenv("TENCENTCLOUD_SECRET_KEY"),
//         )

//         // Nonessential steps
//         // Instantiate a client configuration object. You can specify the timeout period and other configuration items
//         cpf := profile.NewClientProfile()
//         // The SDK uses the POST method by default.
//         // If you have to use the GET method, you can set it here, but the GET method cannot handle some large requests.
//         // Do not modify the default settings unless absolutely necessary.
//         //cpf.HttpProfile.ReqMethod = "GET"
//         // The SDK has a default timeout period. Do not adjust it unless absolutely necessary.
//         // If needed, check in the code to get the latest default value.
//         //cpf.HttpProfile.ReqTimeout = 10
//         // The SDK automatically specifies the domain name. Generally, you don't need to specify a domain name, but if you are accessing a service in a finance availability zone,
//         // you have to manually specify the domain name, such as cvm.ap-shanghai-fsi.tencentcloudapi.com for the Shanghai Finance availability zone in the CVM
//         //cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
//         // The SDK uses HmacSHA256 for signing by default, which is more secure but slightly reduces the performance.
//         // Do not modify the default settings unless absolutely necessary.
//         //cpf.SignMethod = "HmacSHA1"
//         // The SDK uses `zh-CN` calls to return Chinese content by default. You can also set it to `en-US` to return English content.
//         // However, most products or APIs do not fully support returns in English.
//         // Do not modify the default settings unless absolutely necessary.
//         //cpf.Language = "en-US"

//         // Instantiate the client object of the requested product (with CVM as an example)
//         // The second parameter is the region information. You can enter the string `ap-guangzhou` directly or import the preset constant
//         client, _ := cvm.NewClient(credential, regions.Guangzhou, cpf)
//         // Instantiate a request object. You can further set the request parameters according to the API called and actual conditions
//         // You can check the SDK source code directly to determine which attributes of `DescribeInstancesRequest` can be set.
//         // An attribute may be of a basic type or import another data structure.
//         // You are recommended to use the IDE for development where you can redirect to and view the documentation of each API and data structure easily
//         request := cvm.NewDescribeInstancesRequest()

//         // Settings of a basic parameter.
//         // This API allows setting the number of instances returned, which is specified as only one here.
//         // The SDK uses the pointer style to specify parameters, so even for basic parameters, you need to use pointers to assign values to them.
//         // The SDK provides encapsulation functions for importing the pointers of basic parameters
//         request.Limit = common.Int64Ptr(1)

//         // Settings of an array.
//         // This API allows filtering by specified instance ID; however, as it conflicts with the `Filter` parameter to be demonstrated next, it is commented out here.
//         // request.InstanceIds = common.StringPtrs([]string{"ins-r8hr2upy"})

//         // Settings of a complex object.
//         // In this API, `Filters` is an array whose elements are complex objects `Filter`, and the member `Values` of `Filter` are string arrays.
//         request.Filters = []*cvm.Filter{
//             &cvm.Filter{
//                 Name: common.StringPtr("zone"),
//                 Values: common.StringPtrs([]string{"ap-guangzhou-1"}),
//             },
//         }

//         // Use a JSON string to set a request. Note that this is actually an update request, that is, `Limit=1` will be retained,
//         // and the `zone` of the filter will be changed to `ap-guangzhou-2`.
//         // If you need a new request, you should create it with `cvm.NewDescribeInstancesRequest()`.
//         err := request.FromJsonString(`{"Filters":[{"Name":"zone","Values":["ap-guangzhou-2"]}]}`)
//         if err != nil {
//                 panic(err)
//         }
//         // Call the API you want to access through the client object. You need to pass in the request object
//         response, err := client.DescribeInstances(request)
//         // Handle the exception
//         if _, ok := err.(*errors.TencentCloudSDKError); ok {
//                 fmt.Printf("An API error has returned: %s", err)
//                 return
//         }
//         // This is a direct failure instead of SDK exception. You can add other troubleshooting measures in the real code.
//         if err != nil {
//                 panic(err)
//         }
//         // Print the returned JSON string
//         fmt.Printf("%s", response.ToJsonString())
// }