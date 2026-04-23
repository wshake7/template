package kongc

import (
	"errors"
	"fmt"
	"go-common/utils/pt"
	"go.uber.org/zap"
)

type Protocol string

const (
	ProtocolGrpc           = Protocol("grpc")
	ProtocolGrpcs          = Protocol("grpcs")
	ProtocolHttp           = Protocol("http")
	ProtocolHttps          = Protocol("https")
	ProtocolTcp            = Protocol("tcp")
	ProtocolTls            = Protocol("tls")
	ProtocolTlsPassthrough = Protocol("tls_passthrough")
	ProtocolUdp            = Protocol("udp")
	ProtocolWs             = Protocol("ws")
	ProtocolWss            = Protocol("wss")
)

type Conf struct {
	Address   string
	WorkSpace string
}

type Client struct {
	*Conf
	*upstreamService
	*targetService
	*serviceService
	*routeService
}

func Init(address, workSpace string) *Client {
	cfg := &Conf{Address: address, WorkSpace: workSpace}
	c := &Client{
		Conf:            cfg,
		upstreamService: &upstreamService{Conf: cfg},
		targetService:   &targetService{Conf: cfg},
		serviceService:  &serviceService{Conf: cfg},
		routeService:    &routeService{Conf: cfg},
	}
	if c.WorkSpace == "" {
		c.WorkSpace = "default"
	}
	return c
}

type RegisterHttpOption struct {
	Host         string
	Port         int
	ServiceName  string
	Weight       int
	Protocol     Protocol
	RouteOptions []RouteOption
}

func (c *Client) RegisterHttpWithOption(option *RegisterHttpOption) error {
	if option == nil {
		return errors.New("option is nil")
	}
	// 创建 upstream
	upstream, err := c.upstreamService.Save(&Upstream{
		Name: pt.String(option.ServiceName),
	})
	if err != nil {
		return fmt.Errorf("create upstream failed: %w", err)
	}
	if upstream == nil || upstream.ID == nil {
		return errors.New("upstream create failed, upstream id can not be nil")
	}

	// 创建 target
	target := &Target{
		Upstream: &Upstream{
			ID: upstream.ID,
		},
		Target: pt.String(fmt.Sprintf("%s:%d", option.Host, option.Port)),
	}
	err = c.targetService.Create(target)
	if err != nil {
		var apiErr *ApiErr
		if errors.As(err, &apiErr) && apiErr.Code != nil && *apiErr.Code == UniqueErr {
			zap.S().Warn("upstream already includes target")
		} else {
			return fmt.Errorf("create target failed: %w", err)
		}
	}
	if option.Protocol == "" {
		option.Protocol = ProtocolHttp
	}
	// 创建 service
	service := &Service{
		Name:     pt.String(option.ServiceName),
		Host:     pt.String(option.ServiceName),
		Port:     pt.Int(0),
		Protocol: pt.String(string(option.Protocol)),
	}
	_, err = c.serviceService.Save(service)
	if err != nil {
		var apiErr *ApiErr
		if errors.As(err, &apiErr) && apiErr.Code != nil && *apiErr.Code == UniqueErr {
			zap.S().Warn("service already exists")
		} else {
			return fmt.Errorf("create service failed: %w", err)
		}
	}

	// 创建 routes
	for _, routeOption := range option.RouteOptions {
		route := &Route{
			Paths:     pt.Slice[string](routeOption.Paths...),
			StripPath: pt.Bool(routeOption.StripPath),
			Service:   service,
		}
		if routeOption.Name != "" {
			route.Name = pt.String(routeOption.Name)
		}
		for _, method := range routeOption.Methods {
			route.Methods = append(route.Methods, pt.String(string(method)))
		}
		if len(routeOption.Protocols) == 0 {
			route.Protocols = []*string{pt.String(string(ProtocolHttp)), pt.String(string(ProtocolHttps))}
		}
		for _, protocol := range routeOption.Protocols {
			route.Protocols = append(route.Protocols, pt.String(string(protocol)))
		}
		_, err = c.routeService.Save(route)
		if err != nil {
			var apiErr *ApiErr
			if errors.As(err, &apiErr) && apiErr.Code != nil && *apiErr.Code == UniqueErr {
				zap.S().Warn("route already exists")
			} else {
				return fmt.Errorf("create route failed: %w", err)
			}
		}
	}

	return nil
}

func (c *Client) RegisterHttp(host string, port int, serviceName string) error {
	return c.RegisterHttpWithOption(&RegisterHttpOption{
		Host:         host,
		Port:         port,
		ServiceName:  serviceName,
		RouteOptions: nil,
	})
}
