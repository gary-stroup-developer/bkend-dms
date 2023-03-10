package models

import (
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
	Firstname string `bson:"firstname" json:"firstname"`
	Lastname  string `bson:"lastname" json:"lastname"`
	Password  string `bson:"password" json:"password"`
	Email     string `bson:"email" json:"email"`
}

type Product struct {
	Cat_Num  string  `bson:"cat_num" json:"cat_num"`
	Cat_Desc string  `bson:"cat_desc" json:"cat_desc"`
	Raw_PN   string  `bson:"raw_pn" json:"raw_pn"`
	Raw_Desc string  `bson:"raw_desc" json:"raw_desc"`
	Weight   float64 `bson:"weight" json:"weight"`
}

type FSR struct {
	Num_Vials   string `bson:"num_vials,omitempty" json:"num_vials,omitempty"`
	Location    string `bson:"location,omitempty" json:"location,omitempty"`
	Description string `bson:"fsr_desc,omitempty" json:"fsr_desc,omitempty"`
	Due_Date    string `bson:"due_date,omitempty" json:"due_date,omitempty"`
	JobID       string `bson:"jobid,omitempty" json:"jobid"`
	UID         string `bson:"uid" json:"uid"`
}

type CBR struct {
	Source         string `json:"cbr_source,omitempty"`
	Orig_Viability string `json:"orig_viability,omitempty"`
	Orig_Count     string `json:"orig_count,omitempty"`
	Media          string `json:"media,omitempty"`
	Num_Vials      string `json:"num_vials,omitempty"`
	PN             string `json:"cbr_pn,omitempty"`
	LN             string `json:"cbr_ln,omitempty"`
	Description    string `json:"cbr_desc,omitempty"`
	Density        string `json:"density,omitempty"`
	Vial_Init      string `json:"vial_init,omitempty"`
	Tank           string `bson:"tank" json:"tank"`
	Cane           string `bson:"cane" json:"cane"`
	Box            string `bson:"box" json:"box"`
	Position       string `bson:"pos" json:"pos"`
	JobID          string `bson:"jobid,omitempty" json:"jobid"`
	UID            string `bson:"uid" json:"uid"`
}

type Job struct {
	OID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ID         string             `bson:"id,omitempty" json:"id"`
	Cat_Num    string             `bson:"cat_num,omitempty" json:"cat_num,omitempty"`
	Cat_Desc   string             `bson:"cat_desc,omitempty" json:"cat_desc,omitempty"`
	Cat_Lot    string             `bson:"cat_lot,omitempty" json:"cat_lot,omitempty"`
	Raw_PN     string             `bson:"raw_pn,omitempty" json:"raw_pn,omitempty"`
	Raw_Desc   string             `bson:"raw_desc,omitempty" json:"raw_desc,omitempty"`
	Qty        string             `bson:"qty,omitempty" json:"qty,omitempty"`
	Start_date string             `bson:"start_date,omitempty" json:"start_date,omitempty"`
	End_date   string             `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Notes      string             `bson:"notes,omitempty" json:"notes,omitempty"`
	Weight     float64            `bson:"weight,omitempty" json:"weight,omitempty"`
	Status     string             `bson:"status,omitempty" json:"status,omitempty"`
	UID        string             `bson:"uid,omitempty" json:"uid,omitempty"`
	FSR        FSR                `bson:"fsr,omitempty" json:"fsr,omitempty"`
	CBR        CBR                `bson:"cbr,omitempty" json:"cbr,omitempty"`
}

type DashboardResponse struct {
	Users []primitive.M
	Jobs  []primitive.M
}
