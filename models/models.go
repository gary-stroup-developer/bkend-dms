package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UID       string `bson:"uid" json:"uid"`
	Firstname string `bson:"firstname" json:"firstname"`
	Lastname  string `bson:"lastname" json:"lastname"`
	Password  string `bson:"password" json:"password"`
	Email     string `bson:"email" json:"email"`
	Status    bool   `bson:"status" json:"status"`
	Role      string `bson:"role" json:"role"`
}

type Employee struct {
	UID       string `bson:"uid" json:"uid"`
	Firstname string `bson:"firstname" json:"fname"`
	Lastname  string `bson:"lastname" json:"lname"`
	Password  string `bson:"password" json:"password"`
	Email     string `bson:"email" json:"email"`
}

type Product struct {
	CatNum      string
	Description string
	RawPN       string
	RawDesc     string
	Weight      float64
}

type FSR struct {
	Tank        string
	Cane        string
	Box         string
	Position    string
	PN          string
	LN          string
	Description string
	Qty         int64
}

type CBR struct {
	Source         string
	Orig_Viability string
	Orig_Count     string
	Num_Vials      string
	PN             string
	LN             string
	Description    string
	Density        string
	Vial_Init      string
	Location       []FSR
}

type Job struct {
	ID         string    `bson:"id,omitempty" json:"id"`
	Cat_Num    string    `bson:"cat_num,omitempty" json:"cat_num,omitempty"`
	Cat_Desc   string    `bson:"cat_desc,omitempty" json:"cat_desc,omitempty"`
	Cat_Lot    string    `bson:"cat_lot,omitempty" json:"cat_lot,omitempty"`
	Raw_PN     string    `bson:"raw_pn,omitempty" json:"raw_pn,omitempty"`
	Raw_Desc   string    `bson:"raw_desc,omitempty" json:"raw_desc,omitempty"`
	Qty        string    `bson:"qty,omitempty" json:"qty,omitempty"`
	Start_date time.Time `bson:"start_date,omitempty" json:"start_date,omitempty"`
	End_date   time.Time `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Notes      string    `bson:"notes,omitempty" json:"notes,omitempty"`
	Weight     float64   `bson:"weight,omitempty" json:"weight,omitempty"`
	Status     string    `bson:"status,omitempty" json:"status,omitempty"`
	UID        string    `bson:"uid,omitempty" json:"uid,omitempty"`
	FSR        string    `bson:"fsr,omitempty" json:"fsr,omitempty"`
	CBR        CBR       `bson:"cbr,omitempty" json:"cbr,omitempty"`
	Late       bool      `bson:"late,omitempty" json:"late,omitempty"`
	Count      int       `bson:"count,omitempty" json:"count,omitempty"`
}

type DashboardResponse struct {
	Users []primitive.M
	Jobs  []primitive.M
}
