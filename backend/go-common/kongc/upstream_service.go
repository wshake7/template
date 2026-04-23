package kongc

import (
	"fmt"
	"go-common/utils/httpc"
	"net/http"
)

type upstreamService struct {
	*Conf
}

func (u *upstreamService) Create(upstreamReq *Upstream) error {
	if upstreamReq == nil {
		return fmt.Errorf("upstreamReq is nil")
	}
	upstreamRes, err := httpc.PostJsonMarshal[*Upstream](httpc.Request{
		Url:  fmt.Sprintf("%s/%s/upstreams", u.Address, u.WorkSpace),
		Body: upstreamReq,
	}, func(request *http.Request) {})
	if err != nil {
		return err
	}
	if upstreamRes.ApiErr != nil {
		return upstreamRes.ApiErr
	}
	return nil
}

func (u *upstreamService) Save(upstreamReq *Upstream) (*Upstream, error) {
	if upstreamReq == nil {
		return nil, fmt.Errorf("upstreamReq is nil")
	}
	upstreamIdOrName := ""
	if upstreamReq.ID != nil {
		upstreamIdOrName = *upstreamReq.ID
	}
	if upstreamReq.Name != nil {
		upstreamIdOrName = *upstreamReq.Name
	}
	if upstreamIdOrName == "" {
		return nil, fmt.Errorf("upstreamIdOrName is nil")
	}
	upstreamRes, err := httpc.PutJsonMarshal[*Upstream](httpc.Request{
		Url:  fmt.Sprintf("%s/%s/upstreams/%s", u.Address, u.WorkSpace, upstreamIdOrName),
		Body: upstreamReq,
	}, func(request *http.Request) {})
	if err != nil {
		return nil, err
	}
	if upstreamRes.ApiErr != nil {
		return nil, upstreamRes.ApiErr
	}
	return upstreamRes, nil
}

func (u *upstreamService) Get(upstreamIdOrName string) (*Upstream, error) {
	return httpc.GetMarshal[*Upstream](httpc.Request{
		Url: fmt.Sprintf("%s/%s/upstreams/%s", u.Address, u.WorkSpace, upstreamIdOrName),
	}, func(request *http.Request) {})
}
