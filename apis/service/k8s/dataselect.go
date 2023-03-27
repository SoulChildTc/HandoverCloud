package k8s

import (
	"sort"
	"strings"
	"time"
)

// DataCell 目的是将k8s对象的类型转换为DataCell类型，做统一处理，可以处理分页排序
// corev1.Pod -> podCell -> DataCell
// appsv1.deployment -> deployCell -> DataCell
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

type DataSelectQuery struct {
	Filter   *FilterQuery
	Paginate *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

// DataSelect 用于对GenericDataList数据排序、过滤、分页
type DataSelect struct {
	GenericDataList []DataCell
	DataSelect      *DataSelectQuery
}

// Len 实现sort接口Len方法，返回GenericDataList的长度
func (d *DataSelect) Len() int {
	return len(d.GenericDataList)
}

// Swap 实现切片数据位置交换方法
func (d *DataSelect) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less 按照创建时间比大小
func (d *DataSelect) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	// a是否再b之后
	return a.After(b)
}

func (d *DataSelect) Sort() *DataSelect {
	sort.Sort(d)
	return d
}

// Filter 根据名称过滤数据
func (d *DataSelect) Filter() *DataSelect {
	if d.DataSelect.Filter.Name == "" {
		return d
	}
	var filtered []DataCell
	for _, value := range d.GenericDataList {
		objName := value.GetName()
		if !strings.Contains(objName, d.DataSelect.Filter.Name) {
			continue
		}
		filtered = append(filtered, value)
	}
	d.GenericDataList = filtered
	return d
}

// Paginate 数据分页
func (d *DataSelect) Paginate() *DataSelect {
	limit := d.DataSelect.Paginate.Limit // limit代表一页有多少条数据
	page := d.DataSelect.Paginate.Page   //  page代表要第几页的数据
	// 分页参数不合法, 返回空数据
	if limit <= 0 || page <= 0 {
		d.GenericDataList = []DataCell{}
		return d
	}

	startIndex := limit * (page - 1) // (page-1)=上一页  上一页*limit=上一页最后一条数据
	endIndex := limit * page

	if startIndex > len(d.GenericDataList)-1 { // startIndex不能超过切片的最大下标
		d.GenericDataList = []DataCell{}
		return d
	}

	if endIndex > len(d.GenericDataList) { // endIndex超过切片最大大小，直接获取剩余所有数据
		d.GenericDataList = d.GenericDataList[startIndex:]
	} else {
		d.GenericDataList = d.GenericDataList[startIndex:endIndex]

	}

	return d
}
