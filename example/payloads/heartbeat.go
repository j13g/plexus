package payloads

import "github.com/j13g/plexus/config"

func GetHeartbeat() *HeartbeatPayload {
	return &HeartbeatPayload{NodeInfo: GetNodeInfo()}
}

type HeartbeatPayload struct {
	NodeInfo NodeInfo `json:"node_info"`
}

func (h HeartbeatPayload) GetName() string {
	return "Heartbeat"
}

func (h HeartbeatPayload) GetVersion() string {
	return "1.0.0"
}

func GetNodeInfo() NodeInfo {
	cfg := config.Get()
	return NodeInfo{
		NodeID:     cfg.NodeID,
		NodeArea:   cfg.NodeArea,
		AppName:    cfg.AppName,
		AppVersion: "", // TODO
	}
}

type NodeInfo struct {
	NodeID     string `json:"node_id"`
	NodeArea   string `json:"node_area"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
}
