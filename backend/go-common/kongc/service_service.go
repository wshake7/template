package kongc

import (
	"fmt"
	"go-common/utils/httpc"
	"net/http"
)

type serviceService struct {
	*Conf
}

func (s *serviceService) Create(serviceReq *Service) error {
	if serviceReq == nil {
		return fmt.Errorf("serviceReq is nil")
	}
	serviceRes, err := httpc.PostJsonMarshal[*Service](httpc.Request{
		Url:  fmt.Sprintf("%s/%s/services", s.Address, s.WorkSpace),
		Body: serviceReq,
	}, func(request *http.Request) {})
	if err != nil {
		return err
	}
	if serviceRes.ApiErr != nil {
		return serviceRes.ApiErr
	}
	return nil
}

func (s *serviceService) Save(serviceReq *Service) (*Service, error) {
	if serviceReq == nil {
		return nil, fmt.Errorf("serviceReq is nil")
	}
	serviceIdOrName := ""
	if serviceReq.ID != nil {
		serviceIdOrName = *serviceReq.ID
	}
	if serviceReq.Name != nil {
		serviceIdOrName = *serviceReq.Name
	}
	if serviceIdOrName == "" {
		return nil, fmt.Errorf("serviceIdOrName is nil")
	}
	serviceRes, err := httpc.PutJsonMarshal[*Service](httpc.Request{
		Url:  fmt.Sprintf("%s/%s/services/%s", s.Address, s.WorkSpace, serviceIdOrName),
		Body: serviceReq,
	}, func(request *http.Request) {})
	if err != nil {
		return nil, err
	}
	if serviceRes.ApiErr != nil {
		return nil, serviceRes.ApiErr
	}
	return serviceRes, nil
}
