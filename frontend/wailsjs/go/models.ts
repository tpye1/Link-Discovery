export namespace main {
	
	export class ConnectDTO {
	    id: number;
	    name: string;
	    network_card: string;
	    mac_addr: string;
	    ip_addr: string;
	    status: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ConnectDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.network_card = source["network_card"];
	        this.mac_addr = source["mac_addr"];
	        this.ip_addr = source["ip_addr"];
	        this.status = source["status"];
	    }
	}
	export class LinkDataDTO {
	    switch_name: string;
	    port_id: string;
	    vlan_id: string;
	    switch_ip: string;
	    switch_model: string;
	    duplex_option: string;
	    vpt_mgmt_domain: string;
	    protocol: string;
	
	    static createFrom(source: any = {}) {
	        return new LinkDataDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.switch_name = source["switch_name"];
	        this.port_id = source["port_id"];
	        this.vlan_id = source["vlan_id"];
	        this.switch_ip = source["switch_ip"];
	        this.switch_model = source["switch_model"];
	        this.duplex_option = source["duplex_option"];
	        this.vpt_mgmt_domain = source["vpt_mgmt_domain"];
	        this.protocol = source["protocol"];
	    }
	}

}

