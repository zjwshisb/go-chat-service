// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Students is the golang structure of table students for DAO operations like Where/Data.
type Students struct {
	g.Meta             `orm:"table:students, do:true"`
	Id                 interface{} //
	CustomerId         interface{} // 客户id
	UserId             interface{} // 用户Id
	SchoolId           interface{} // 学校id
	GradeId            interface{} // 年纪id
	ClassId            interface{} // 班级Id
	EatGradeId         interface{} // 就餐年级
	EatClassId         interface{} // 就餐班级
	Name               interface{} // 名称
	NamePinyin         interface{} //
	Sex                interface{} // 性别
	CreatedAt          interface{} //
	UpdatedAt          interface{} //
	IsSpecificRecharge interface{} // 开启单独充值,不受学校充值设置影响
	IsSpecificRefund   interface{} // 开启单独退款，不受学校退款设置影响
	TypeId             interface{} //
	IsAuto             interface{} // 是否开启自动订餐
	IsHide             interface{} // 是否在前台隐藏：0显示，1隐藏
	IsFreeze           interface{} // 非选餐情况下是否不加单
	IsCancel           interface{} // 是否注销 0否, 1是
}
