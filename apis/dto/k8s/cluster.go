package k8s

type TlsClientConfig struct {
	Insecure bool   `json:"insecure"`
	CertData string `json:"certData" binding:"pem=cert" pem_err:"客户端证书格式错误"`
	KeyData  string `json:"keyData" binding:"required_with=CertData,pem=key" msg:"客户端私钥必填" pem_err:"私钥格式错误"`
	CAData   string `json:"CAData" binding:"required_if=Insecure false,len=0|pem=cert" msg:"未启用Insecure, CA证书必填" pem_err:"CA证书格式错误"`
}

type ClusterInfo struct {
	ClusterName     string          `json:"clusterName"`
	Host            string          `json:"host" binding:"http_url" msg:"Host必须是http(s) url"`
	BearerToken     string          `json:"bearerToken" binding:"jwt|len=0,required_without=TLSClientConfig.CertData" msg:"token必须为jwt格式" required_without_err:"BearerToken和tls客户端认证二选一"`
	TLSClientConfig TlsClientConfig `json:"TLSClientConfig"`
	Version         string          `json:"version,omitempty"`
}
