package k8s

import "encoding/json"

type SetImage []struct {
	Name  string `json:"name"`
	Image string `json:"image" binding:"required" msg:"镜像地址不能为空"`
}

func (d SetImage) NameNotEmpty() bool {
	for _, item := range d {
		if item.Name == "" {
			return false
		}
	}
	return true
}

func (d SetImage) String() string {
	data, err := json.Marshal(d)
	if err != nil {
		return ""
	}
	return string(data)
}
