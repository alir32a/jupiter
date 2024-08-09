package ocserv

type ConnectionEntity struct {
	ID                 int           `json:"ID"`
	Username           string        `json:"UserID"`
	Groupname          string        `json:"Groupname"`
	State              string        `json:"State"`
	Vhost              string        `json:"vhost"`
	Device             string        `json:"Device"`
	MTU                string        `json:"MTU"`
	RemoteIP           string        `json:"Remote IP"`
	Location           string        `json:"Location"`
	LocalDeviceIP      string        `json:"Local Device IP"`
	IPv4               string        `json:"IPv4"`
	PTPIPv4            string        `json:"P-t-P IPv4"`
	IPv6               string        `json:"IPv6"`
	PTPIPv6            string        `json:"P-t-P IPv6"`
	UserAgent          string        `json:"User-Agent"`
	RX                 string        `json:"RX"`
	TX                 string        `json:"TX"`
	AverageRX          string        `json:"Average RX"`
	AverageTX          string        `json:"Average TX"`
	DPD                string        `json:"DPD"`
	KeepAlive          string        `json:"KeepAlive"`
	Hostname           string        `json:"Hostname"`
	ConnectedAt        string        `json:"Connected at"`
	RawConnectedAt     int           `json:"raw_connected_at"`
	FullSession        string        `json:"Full session"`
	Session            string        `json:"Session"`
	TLSCiphersuite     string        `json:"TLS ciphersuite"`
	DNS                []string      `json:"DNS"`
	NBNS               []interface{} `json:"NBNS"`
	SplitDNSDomains    []interface{} `json:"Split-DNS-Domains"`
	Routes             string        `json:"Routes"`
	NoRoutes           []interface{} `json:"No-routes"`
	IRoutes            []interface{} `json:"iRoutes"`
	RestrictedToRoutes string        `json:"Restricted to routes"`
	RestrictedToPorts  []interface{} `json:"Restricted to ports"`
}
