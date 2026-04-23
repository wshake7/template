package kongc

import (
	"fmt"
	"go-common/utils/httpc"
	"net/http"
)

type targetService struct {
	*Conf
}

func (u *targetService) Create(targetReq *Target) error {
	if targetReq == nil {
		return fmt.Errorf("targetReq is nil")
	}
	upstreamIdOrName := ""
	if targetReq.Upstream != nil {
		upstream := targetReq.Upstream
		if upstream == nil {
			return fmt.Errorf("targetReq.Upstream is nil")
		}
		upstreamIdOrName = upstream.FriendlyName()
	}
	if upstreamIdOrName == "" {
		return fmt.Errorf("upstreamIdOrName is nil")
	}
	targetRes, err := httpc.PostJsonMarshal[*Target](httpc.Request{
		Url:  fmt.Sprintf("%s/%s/upstreams/%s/targets", u.Address, u.WorkSpace, upstreamIdOrName),
		Body: targetReq,
	}, func(request *http.Request) {})
	if err != nil {
		return err
	}
	if targetRes.ApiErr != nil {
		return targetRes.ApiErr
	}
	return nil
}

func (u *targetService) Save(targetReq *Target) (*Target, error) {
	if targetReq == nil {
		return nil, fmt.Errorf("targetReq is nil")
	}
	upstreamIdForTarget := ""
	if targetReq.Upstream != nil {
		upstream := targetReq.Upstream
		if upstream == nil {
			return nil, fmt.Errorf("targetReq.Upstream is nil")
		}
		if upstream.ID == nil {
			return nil, fmt.Errorf("upstream.ID is nil")
		}
		upstreamIdForTarget = *upstream.ID
	}
	if upstreamIdForTarget == "" {
		return nil, fmt.Errorf("upstreamIdForTarget is nil")
	}
	targetIdOrTarget := ""
	if targetReq.ID != nil {
		targetIdOrTarget = *targetReq.ID
	}
	if targetReq.Target != nil {
		targetIdOrTarget = *targetReq.Target
	}
	if targetIdOrTarget == "" {
		return nil, fmt.Errorf("targetIdOrTarget is nil")
	}
	targetRes, err := httpc.PutJsonMarshal[*Target](httpc.Request{
		Url:  fmt.Sprintf("%s/%s/upstreams/%s/targets/%s", u.Address, u.WorkSpace, upstreamIdForTarget, targetIdOrTarget),
		Body: targetReq,
	}, func(request *http.Request) {})
	if err != nil {
		return nil, err
	}
	if targetRes.ApiErr != nil {
		return nil, targetRes.ApiErr
	}

	return targetRes, nil
}
