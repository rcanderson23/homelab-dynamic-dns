package ip

type Lookup interface {
	GetIP() (string, error)
}

type HttpLookup struct {
	Url string `json:"url"`
}

type EdgeRouter struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
	Host      string `json:"host"`
	Port      uint16 `json:"port"`
}
