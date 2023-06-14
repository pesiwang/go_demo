package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

type NsqNodeInfo struct {
	RemoteAddr    string   `json:"remote_address"`
	Hostname      string   `json:"hostname"`
	BroadcastAddr string   `json:"broadcast_address"`
	TcpPort       uint32   `json:"tcp_port"`
	HttpPort      uint32   `json:"http_port"`
	Version       string   `json:"version"`
	Tombstones    []bool   `json:"tombstones"`
	Topics        []string `json:"topics"`
}

type Nodes struct {
	Producers []NsqNodeInfo `json:"producers"`
}

type NsqdHelper struct {
}

func NewNsqdHelper() *NsqdHelper {
	return &NsqdHelper{}
}

func (s *NsqdHelper) getNodesIps(ip string) ([]string, error) {

	nodeUrl := fmt.Sprintf("http://%s/nodes", ip)

	resp, body, errs := gorequest.New().Timeout(2000 * time.Millisecond).Get(nodeUrl).End()
	if errs != nil && len(errs) > 0 && errs[0] != nil {
		var err error
		for _, err = range errs {
			fmt.Printf("[GetNodesIps] http get nodes ip failed, %v\n", err)
		}
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http get nodes ip invalid , http code is:%v", resp.StatusCode)
	}

	result := Nodes{}

	if err := json.Unmarshal([]byte(body), &result); err != nil {
		fmt.Printf("error %v body %s\n", err, body)
		return nil, err
	}

	var nsqdsUrl []string

	for _, v := range result.Producers {
		ip := strings.Split(v.RemoteAddr, ":")
		if len(ip) == 2 {
			url := fmt.Sprintf("%v:%v", ip[0], v.TcpPort)
			nsqdsUrl = append(nsqdsUrl, url)
		}
	}

	return nsqdsUrl, nil

}

func (s *NsqdHelper) GetNodesIps(lookupAddrs []string) ([]string, error) {
	for _, v := range lookupAddrs {
		nodes, err := s.getNodesIps(v)
		if err == nil {
			return nodes, nil
		}
		fmt.Printf("error %v\n", err)
	}
	return nil, fmt.Errorf("get nsqd nodes from lookups[%+v] failed", lookupAddrs)
}
