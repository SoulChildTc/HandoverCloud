package model

var MigrateModels []any = []any{
	//your model. eg: system.User{},
	&SystemUser{},
	&K8sCluster{},
}
