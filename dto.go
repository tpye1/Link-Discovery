package main

type ConnectDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	NetworkCard string `json:"network_card"`
	MacAddr     string `json:"mac_addr"`
	IpAddr      string `json:"ip_addr"`
	Status      bool   `json:"status"`
}

type LinkDataDTO struct {
	SwitchName    string `json:"switch_name"`
	PortID        string `json:"port_id"`
	VlanID        string `json:"vlan_id"`
	SwitchIP      string `json:"switch_ip"`
	SwitchModel   string `json:"switch_model"`
	DuplexOption  string `json:"duplex_option"`
	VptMgmtDomain string `json:"vpt_mgmt_domain"`
	Protocol      string `json:"protocol"`
}

func toConnectDTO(c connect_struct) ConnectDTO {
	return ConnectDTO{
		ID:          c.id,
		Name:        c.name,
		NetworkCard: c.network_card,
		MacAddr:     c.mac_addr,
		IpAddr:      c.ip_addr,
		Status:      c.status,
	}
}

func toLinkDataDTO(l link_data) LinkDataDTO {
	return LinkDataDTO{
		SwitchName:    l.switch_name,
		PortID:        l.port_id,
		VlanID:        l.vlan_id,
		SwitchIP:      l.switch_ip,
		SwitchModel:   l.switch_model,
		DuplexOption:  l.duplex_option,
		VptMgmtDomain: l.vpt_mgmt_domain,
		Protocol:      l.protocol,
	}
}
