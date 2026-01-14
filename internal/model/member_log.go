package model

import "time"

type MemberLog struct {
	ID                  int64     `db:"id" json:"id"`
	SiteBid             string    `db:"site_bid" json:"site_bid"`
	Type                int16     `db:"type" json:"type"`
	ExtMemberID         int       `db:"ext_member_id" json:"ext_member_id"`
	ExtUsername         string    `db:"ext_username" json:"ext_username"`
	MemberID            int       `db:"member_id" json:"member_id"`
	Scene               int16     `db:"scene" json:"scene"`
	AdminID             int       `db:"admin_id" json:"admin_id"`
	OperatorType        int16     `db:"operator_type" json:"operator_type"`
	OperatorExtMemberID int64     `db:"operator_ext_member_id" json:"operator_ext_member_id"`
	OperatorUsername    string    `db:"operator_username" json:"operator_username"`
	OperatorMemberID    int64     `db:"operator_member_id" json:"operator_member_id"`
	Remark              string    `db:"remark" json:"remark"`
	CreateTime          time.Time `db:"create_time" json:"create_time"`
}
