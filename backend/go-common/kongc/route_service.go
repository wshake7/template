package kongc

import (
	"fmt"
	"go-common/utils/httpc"
	"net/http"
)

type routeService struct {
	*Conf
}

func (s *routeService) Create(routeReq *Route) error {
	if routeReq == nil {
		return fmt.Errorf("routeReq is nil")
	}
	routeRes, err := httpc.PostJsonMarshal[*Route](httpc.Request{
		Url:  fmt.Sprintf("%s/%s/routes", s.Address, s.WorkSpace),
		Body: routeReq,
	}, func(request *http.Request) {})
	if err != nil {
		return err
	}
	if routeRes.ApiErr != nil {
		return routeRes.ApiErr
	}
	return nil
}

func (s *routeService) Save(routeReq *Route) (*Route, error) {
	if routeReq == nil {
		return nil, fmt.Errorf("routeReq is nil")
	}
	routeIdOrName := ""
	if routeReq.ID != nil {
		routeIdOrName = *routeReq.ID
	}
	if routeReq.Name != nil {
		routeIdOrName = *routeReq.Name
	}
	if routeIdOrName == "" {
		return nil, fmt.Errorf("routeIdOrName is nil")
	}
	routeRes, err := httpc.PutJsonMarshal[*Route](httpc.Request{
		Url:  fmt.Sprintf("%s/%s/routes/%s", s.Address, s.WorkSpace, routeIdOrName),
		Body: routeReq,
	}, func(request *http.Request) {})
	if err != nil {
		return nil, err
	}
	if routeRes.ApiErr != nil {
		return nil, routeRes.ApiErr
	}
	return routeRes, nil
}

func (s *routeService) GetByServiceIdOrName(serviceIdOrName string) (*Route, error) {
	return httpc.GetMarshal[*Route](httpc.Request{
		Url: fmt.Sprintf("%s/%s/services/%s/routes", s.Address, s.WorkSpace, serviceIdOrName),
	}, func(request *http.Request) {})
}
