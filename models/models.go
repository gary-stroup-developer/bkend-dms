package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UID       string  `bson:"uid" json:"uid"`
	Firstname string  `bson:"firstname" json:"firstname,omitempty"`
	Lastname  string  `bson:"lastname" json:"lastname,omitempty"`
	Password  string  `bson:"password" json:"password,omitempty"`
	Email     string  `bson:"email" json:"email,omitempty"`
	Status    bool    `bson:"status" json:"status,omitempty"`
	Role      string  `bson:"role" json:"role,omitempty"`
	Capacity  float64 `bson:"capacity" json:"capacity,omitempty"`
}

type Employee struct {
	UID       string `bson:"uid" json:"uid"`
	Firstname string `bson:"firstname" json:"fname"`
	Lastname  string `bson:"lastname" json:"lname"`
	Password  string `bson:"password" json:"password"`
	Email     string `bson:"email" json:"email"`
}

type Product struct {
	Cat_Num     string `bson:"cat_num" json:"cat_num"`
	Description string `bson:"description" json:"description"`
	Raw_PN      string
	Raw_Desc    string
	Weight      float64 `bson:"weight" json:"weight"`
}

type Location struct {
	Tank        string
	Cane        string
	Box         string
	Position    string
	PN          string
	LN          string
	Description string
	Density     string
	Vial_Init   string
	Media       string
}

type CBR struct {
	Source         string     `json:"cbr_source,omitempty"`
	Orig_Viability string     `json:"orig_viability,omitempty"`
	Orig_Count     string     `json:"orig_count,omitempty"`
	Num_Vials      string     `json:"num_vials,omitempty"`
	PN             string     `json:"cbr_pn,omitempty"`
	LN             string     `json:"cbr_ln,omitempty"`
	Description    string     `json:"cbr_desc,omitempty"`
	Density        string     `json:"density,omitempty"`
	Vial_Init      string     `json:"vial_init,omitempty"`
	Location       []Location `json:"location,omitempty"`
}

type Job struct {
	OID        string    `bson:"_id,omitempty" json:"_id"`
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
}

type DashboardResponse struct {
	Users []primitive.M
	Jobs  []primitive.M
}
