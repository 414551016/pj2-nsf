package model

type NnsfNfmPostRequest struct {
	NfName     string   `json:"nf_name"`
	SbiAddress string   `json:"sbi_address"`
	SbiPort    int      `json:"sbi_port"`
	Services   []string `json:"services"`
}

type NnsfNfmPostResponse struct {
	ServiceToken string `json:"service_token"`
	ServiceUuid  string `json:"service_uuid"`
}

type NnsfNfmDeleteRequest struct {
	NfName   string   `json:"nf_name"`
	Services []string `json:"services"`
}
