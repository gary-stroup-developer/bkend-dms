package models

type User struct {
	UID      string `bson:"uid" json:"uid"`
	FName    string `bson:"fname" json:"fname"`
	LName    string `bson:"lname" json:"lname"`
	Password string `bson:"password" json:"password"`
	Email    string `jbson:"email" son:"email"`
	Status   bool   `bson:"status" json:"status"`
	Role     string `bson:"role" json:"role"`
}

type Employee struct {
	UID      string `json:"uid"`
	FName    string `json:"fname"`
	LName    string `json:"lname"`
	Password string `json:"password"`
	Email    string `json:"email"`
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
	ID             string
	Cat_Num        string
	Cat_Desc       string
	Cat_Lot        string
	Raw_PN         string
	Raw_Desc       string
	Qty            string
	Start_date     string
	End_date       string
	Notes          string
	Weight         int
	Status         string
	UID            string
	Late           int //maybe bool
	Frozen_Request []FSR
	Bank_Request   CBR
}
