package goo

const (
	grpcServerPort  = "server-port"
	grpcServiceName = "service-name"
	grpcConsul      = "consul"
)

func GRPCServiceName(serviceName string) Option {
	return NewOption(grpcServiceName, serviceName)
}

func GRPCConsul(address, username, password string) Option {
	return NewOption(grpcConsul, NewConsul(address, username, password))
}
