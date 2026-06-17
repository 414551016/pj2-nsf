package model

type NnsfDiscGetRequest struct {
	Service string `json:"service"`
}

type NnsfDiscGetResponse struct {
	Service    string `json:"service"`
	NfName     string `json:"nf_name"`
	SbiAddress string `json:"sbi_address"`
	SbiPort    int    `json:"sbi_port"`
}

type NnsfDiscGroupGetRequest struct {
	Services []string `json:"services"`
}

type NnsfDiscGroupGetResponse struct {
	Services []NnsfDiscGetResponse `json:"services"`
}
