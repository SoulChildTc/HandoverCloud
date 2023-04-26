package k8s

import (
	"database/sql"
	"soul/model/common"
)

type Cluster struct {
	common.ID
	ClusterName string         `json:"clusterName" gorm:"size:32;not null;uniqueIndex;comment:集群名称"`
	Host        string         `json:"host" gorm:"size:64;not null;comment:ApiServer地址"`
	BearerToken sql.NullString `json:"bearerToken" gorm:"size:256;comment:访问ApiServer的Token"`
	Insecure    bool           `json:"insecure" gorm:"default:false;common:是否不验证服务端TLS证书"`
	CertData    sql.NullString `json:"certData" gorm:"comment:客户端证书"`
	KeyData     sql.NullString `json:"keyData"  gorm:"comment:客户端私钥"`
	CAData      sql.NullString `json:"CAData"  gorm:"comment:CA证书"`
	common.Timestamps
}

func (c Cluster) TableName() string {
	return "t_k8s_cluster"
}
