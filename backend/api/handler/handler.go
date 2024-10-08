package handler

import (
	"context"
	"log"
	"proxy-checker-server/api/service"
	pb "proxy-checker-server/generated/grpc/proxy-checker-api.v1"
	"sync"

	"google.golang.org/grpc"
)

func mapProxyType(proxyType string) pb.ProxyType {
	switch proxyType {
	case "sock5":
		return pb.ProxyType_PROXY_TYPE_SOCK5
	case "https":
		return pb.ProxyType_PROXY_TYPE_HTTPS
	case "http":
		return pb.ProxyType_PROXY_TYPE_HTTP
	default:
		log.Printf("Unsupported proxy type %v", proxyType)
		panic("Unsupported proxy type")
	}
}

func checkProxy(proxy string) *pb.ProxyCheckResult {
	proxyInfo := service.GetProxyInfo(proxy)

	if proxyInfo != nil {
		return &pb.ProxyCheckResult{
			Value:  proxy,
			Active: true,
			Info: &pb.ProxyInfo{
				Type:       mapProxyType(proxyInfo.Type),
				ExternalIp: proxyInfo.ExternalIp,
				Country:    proxyInfo.Country,
				Region:     proxyInfo.Region,
			},
		}
	} else {
		return &pb.ProxyCheckResult{
			Value:  proxy,
			Active: false,
		}
	}
}

type ApiServer struct{}

func (s *ApiServer) Check(
	context context.Context,
	request *pb.ProxyRequest,
) (*pb.ProxyResponse, error) {

	wg := sync.WaitGroup{}

	var result []*pb.ProxyCheckResult

	var activeCounter int64 = 0

	for _, value := range request.Value {
		wg.Add(1)

		go func() {
			log.Printf("Handling proxy: %v", value)

			defer wg.Done()

			info := checkProxy(value)

			if info.Active {
				activeCounter++
			}
			result = append(result, info)
		}()
	}

	wg.Wait()
	return &pb.ProxyResponse{
		Total:       int64(len(result)),
		Active:      activeCounter,
		CheckResult: result,
	}, nil
}

func (s *ApiServer) CheckStream(
	request *pb.ProxyRequest,
	stream pb.ProxyChecker_CheckStreamServer,
) error {
	wg := sync.WaitGroup{}

	for _, value := range request.Value {
		wg.Add(1)

		go func() {
			log.Printf("Handling proxy: %v", value)

			defer wg.Done()

			stream.Send(checkProxy(value))
		}()
	}

	wg.Wait()

	return nil
}

func ServerRegister(server *grpc.Server) {
	pb.RegisterProxyCheckerServer(server, &ApiServer{})
}
