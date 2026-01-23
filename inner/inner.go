package inner

import "context"

const grpcSignServiceName = "sign_service_name"

func GetRequestServiceName(ctx context.Context) string {
	value := ctx.Value(grpcSignServiceName)
	if value == nil {
		return ""
	}
	if serviceName, ok := value.([]string); ok && len(serviceName) > 0 {
		return serviceName[0]
	}
	return ""
}
