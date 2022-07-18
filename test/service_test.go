package test

import (
	"fmt"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/service"
	"github.com/stretchr/testify/require"
	"testing"
)

/***** SERVICE TX START *****/

// 创建一个新的服务定义。
func TestServiceDefine(t *testing.T) {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	definition := service.DefineServiceRequest{
		ServiceName:       "test_serviceName",
		Description:       "this is a test service",
		Tags:              nil,
		AuthorDescription: "service provider",
		Schemas:           schemas,
	}

	result, err := txClient.Service.DefineService(definition, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Hash)
}

// 绑定一个服务。
func TestServiceBind(t *testing.T) {
	deposit, e := sdk.ParseDecCoins("5000stake")
	require.NoError(t, e)
	pricing := `{"price":"1stake"}`
	options := `{}`
	binding := service.BindServiceRequest{
		ServiceName: "test_serviceName",
		Deposit:     deposit,
		Pricing:     pricing,
		QoS:         10,
		Options:     options,
	}
	result, err := txClient.Service.BindService(binding, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Hash)
}

// 更新已存在的服务绑定。
func TestServiceUpdateBinding(t *testing.T) {
	deposit, e := sdk.ParseDecCoins("5000stake")
	require.NoError(t, e)
	pricing := `{"price":"2stake"}`
	binding := service.UpdateServiceBindingRequest{
		ServiceName: "test_serviceName",
		Deposit:     deposit,
		Pricing:     pricing,
		QoS:         10,
	}
	result, err := txClient.Service.UpdateServiceBinding(binding, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Hash)
}

// 禁用一个可用的服务绑定。
func TestServiceDisable(t *testing.T) {
	result, err := txClient.Service.DisableServiceBinding("test_serviceName", address, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Hash)
}

// 启用一个不可用的服务绑定。
func TestServiceEnable(t *testing.T) {
	deposit, e := sdk.ParseDecCoins("500stake")
	require.NoError(t, e)
	result, err := txClient.Service.EnableServiceBinding("test_serviceName", address, deposit, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Hash)
}

// 发起服务调用。
func TestServiceCall(t *testing.T) {
	serviceFeeCap, e := sdk.ParseDecCoins("200stake")
	require.NoError(t, e)
	input := `{"header":{},"body":{"param1":"hello","param2":"world"}}`
	invocation := service.InvokeServiceRequest{
		ServiceName:   "test_serviceName",
		Providers:     []string{address},
		Input:         input,
		ServiceFeeCap: serviceFeeCap,
		Timeout:       10,
		Repeated:      false,
		RepeatedTotal: -1,
	}

	requestContextID, result, err := txClient.Service.InvokeService(invocation, baseTx)
	require.NoError(t, err)
	fmt.Println("InvokeService success",
		"hash", result.Hash,
		"requestContextID", requestContextID,
	)
}

// 响应指定的服务请求。
func TestServiceRespond(t *testing.T) {
	param1 := "hello"
	param2 := "world"
	data := param1 + " " + param2
	output := `{"header":{},"body":{"data":"` + data + `"}}`
	testResult := `{"code":200,"message":"success"}`
	request := service.InvokeServiceResponseRequest{
		RequestId: "60A745D6C42FF68D8EB22015FE40EC3E87161974751798A54D887448E93BC33C0000000000000000000000000000000100000000000065FE0000",
		Output:    output,
		Result:    testResult,
	}

	result, err := txClient.Service.InvokeServiceResponse(request, baseTx)
	require.NoError(t, err)
	fmt.Println(result)
}

// 更新指定的请求上下文。
func TestServiceUpdate(t *testing.T) {
	serviceFeeCap, e := sdk.ParseDecCoins("210stake")
	require.NoError(t, e)
	request := service.UpdateRequestContextRequest{
		RequestContextID:  "530388B593AE1DE177B6DD57C52E08374996835FE47E25316DD8B73C38B6B5A10000000000000000",
		Providers:         []string{address},
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           20,
		RepeatedFrequency: 100,
		RepeatedTotal:     100,
	}
	result, err := txClient.Service.UpdateRequestContext(request, baseTx)
	require.NoError(t, err)
	fmt.Println(result)
}

// 暂停一个正在进行的请求上下文。
func TestServicePause(t *testing.T) {
	res, err := txClient.Service.PauseRequestContext("09CA08860F72F2E95131D64C3EDB4EED64ED6E46EC51DEE58FFD9A61BE7233DF0000000000000000", baseTx)
	require.NoError(t, err)
	fmt.Println(res)
}

// 启动一个暂停的请求上下文。
func TestServiceStart(t *testing.T) {
	res, err := txClient.Service.StartRequestContext("09CA08860F72F2E95131D64C3EDB4EED64ED6E46EC51DEE58FFD9A61BE7233DF0000000000000000", baseTx)
	require.NoError(t, err)
	fmt.Println(res)
}

// 永久终止一个请求上下文。。
func TestServiceKill(t *testing.T) {
	res, err := txClient.Service.KillRequestContext("09CA08860F72F2E95131D64C3EDB4EED64ED6E46EC51DEE58FFD9A61BE7233DF0000000000000000", baseTx)
	require.NoError(t, err)
	fmt.Println(res)
}

// 设置所有者的服务费提取地址。
func TestServiceSetWithdrawAddr(t *testing.T) {
	_, err := txClient.Service.SetWithdrawAddress(address, baseTx)
	require.NoError(t, err)
}

// 提取服务提供者赚取的服务费。如未指定服务提供者，则提取该所有者全部服务提供者的服务费。
func TestServiceWithdrawFees(t *testing.T) {
	_, err := txClient.Service.WithdrawEarnedFees(address, baseTx)
	require.NoError(t, err)
}

/***** SERVICE TX END *****/

/***** SERVICE QUERY START *****/

// 查询服务定义。
func TestServiceDefinition(t *testing.T) {
	defi, err := txClient.Service.QueryServiceDefinition("test_serviceName")
	require.NoError(t, err)
	fmt.Println(defi)
}

// 查询指定的服务绑定。
func TestServiceBinding(t *testing.T) {
	binding, err := txClient.Service.QueryServiceBinding("test_serviceName", address)
	require.NoError(t, err)
	fmt.Println(binding)
}

// 查询指定服务的绑定列表
func TestServiceBindings(t *testing.T) {
	bindings, err := txClient.Service.QueryServiceBindings("test_serviceName", pagination)
	require.NoError(t, err)
	fmt.Println(bindings)
}

// 通过服务绑定查询当前活跃的服务请求列表。
func TestServiceQueryServiceRequests(t *testing.T) {
	requestid, err := txClient.Service.QueryServiceRequests("test_serviceName", address, pagination)
	require.NoError(t, err)
	fmt.Println(requestid)
}

// 通过请求上下文 ID 查询当前活跃的服务请求列表。
func TestServiceQueryRequestsByReqCtx(t *testing.T) {
	requestContextID := "E17F2FD53E9B14B84CCA1169098B3D01B6384D830FFC6BEADCC64B8F7A5C04460000000000000000"
	requestid, err := txClient.Service.QueryRequestsByReqCtx(requestContextID, 1, pagination)
	require.NoError(t, err)
	fmt.Println(requestid)
}

// 通过请求 ID 查询服务请求。
func TestServiceRequest(t *testing.T) {
	res, err := txClient.Service.QueryServiceRequest("E17F2FD53E9B14B84CCA1169098B3D01B6384D830FFC6BEADCC64B8F7A5C0446000000000000000000000000000000010000000000006B150000")
	require.NoError(t, err)
	fmt.Println(res)
}

// 查询指定服务请求的服务响应。
func TestServiceResponse(t *testing.T) {
	res, err := txClient.Service.QueryServiceResponse("E17F2FD53E9B14B84CCA1169098B3D01B6384D830FFC6BEADCC64B8F7A5C0446000000000000000000000000000000010000000000006B150000")
	require.NoError(t, err)
	fmt.Println(res)
}

// 通过请求上下文 ID 以及批次计数器查询活跃的服务响应。
func TestServiceResponses(t *testing.T) {
	res, err := txClient.Service.QueryServiceResponses("E17F2FD53E9B14B84CCA1169098B3D01B6384D830FFC6BEADCC64B8F7A5C04460000000000000000", 1, pagination)
	require.NoError(t, err)
	fmt.Println(res)
}

// 查询指定的请求上下文。
func TestServiceRequestContext(t *testing.T) {
	res, err := txClient.Service.QueryRequestContext("E17F2FD53E9B14B84CCA1169098B3D01B6384D830FFC6BEADCC64B8F7A5C04460000000000000000")
	require.NoError(t, err)
	fmt.Println(res)
}

// 查询指定服务提供者赚取的服务费。
func TestServiceFees(t *testing.T) {
	res, err := txClient.Service.QueryFees(address)
	require.NoError(t, err)
	fmt.Println(res)
}

/***** SERVICE QUERY END *****/
