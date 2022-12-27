package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gary-stroup-developer/bkend-dms/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	DB       *mongo.Database
	UserInfo models.User
}

var Repo *Repository

func CreateRepo(d *mongo.Database) *Repository {
	return &Repository{DB: d}
}

func NewRepo(r *Repository) {
	Repo = r
}

func (m *Repository) Login(res http.ResponseWriter, req *http.Request) {
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Login attempt failed", http.StatusBadRequest)
	}

	var user models.User
	var userPayload models.User

	err = json.Unmarshal(payload, &userPayload)
	if err != nil {
		http.Error(res, "Sorry. There seems to be an issue connecting to the database", http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filter := bson.D{{Key: "uid", Value: userPayload.UID}}
	if err = m.DB.Collection("User").FindOne(ctx, filter).Decode(&user); err != nil {
		http.Error(res, "Sorry. User information unavailable at this time", http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPayload.Password))
	if err != nil {
		http.Error(res, "username and/or password do not match", http.StatusBadRequest)
	}

	//store the user info in the repository to be accessed at other routes
	m.UserInfo = user

	response, err := json.Marshal(&user)
	if err != nil {
		http.Error(res, "Unable to send user info", http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (m *Repository) Dashboard(res http.ResponseWriter, req *http.Request) {
	// create the stages to get users who are active, then return everything except the password
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: true}, {Key: "role", Value: "user"}}}}
	unsetStage := bson.D{{Key: "$unset", Value: bson.A{"password"}}}

	var users []bson.M //will hold the user data retrieved from the database

	//use aggregate func to query results using pipeline
	cursor, err := m.DB.Collection("User").Aggregate(context.TODO(), mongo.Pipeline{matchStage, unsetStage})
	if err != nil {
		http.Error(res, "trouble connecting to server", http.StatusInternalServerError)
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &users) //store all active users in database

	//create stages for job query
	stageOne := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: "wip"}}}}
	stageTwo := bson.D{{Key: "$unset", Value: bson.A{"cat_num", "cat_desc", "cat_lot", "raw_pn", "raw_desc", "qty", "start_date", "end_date", "notes"}}}
	stageThree := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$uid"},
			{Key: "jobs", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}},
	}

	var jobs []bson.M
	//use aggregate func to query results using pipeline
	c, err := m.DB.Collection("Jobs").Aggregate(context.TODO(), mongo.Pipeline{stageOne, stageTwo, stageThree})
	if err != nil {
		http.Error(res, "trouble connecting to server", http.StatusInternalServerError)
	}
	defer c.Close(context.TODO())

	c.All(context.TODO(), &jobs) //store all active users in database
	// userinfo, _ := json.Marshal(&users)
	// data, _ := json.Marshal(&jobs)

	marshalData := models.DashboardResponse{Users: users, Jobs: jobs}
	payload, _ := json.Marshal(&marshalData)
	res.Write(payload)

}

func (m *Repository) TestDash(res http.ResponseWriter, req *http.Request) {

}

func (m *Repository) UserProfile(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "This is the user profile page. Specific to the user with a UID")
}

func (m *Repository) CreateJob(res http.ResponseWriter, req *http.Request) {
	// payload, err := io.ReadAll(req.Body)
	// if err != nil {
	// 	http.Error(res, "request has been denied", 404)
	// }
	// insert code goes here

	// 	Cat_Num        string
	// Cat_Desc       string
	// Cat_Lot        string
	// Raw_PN         string
	// Raw_Desc       string
	// Qty            string
	// Start_date     string
	// End_date       string
	// Notes          string
	// Weight         int
	// Status         string
	// UID            string

	// const shortForm = "2006-01-02"
	// time1, _ := time.Parse(shortForm, "2022-10-29")
	// time2, _ := time.Parse(shortForm, "2022-10-31")
	// time3, _ := time.Parse(shortForm, "2022-11-02")
	// time4, _ := time.Parse(shortForm, "2022-12-02")
	// time5, _ := time.Parse(shortForm, "2022-12-16")
	// time6, _ := time.Parse(shortForm, "2022-12-29")
	// start1 := primitive.NewDateTimeFromTime(time1)
	// start2 := primitive.NewDateTimeFromTime(time2)
	// start3 := primitive.NewDateTimeFromTime(time3)
	// start4 := primitive.NewDateTimeFromTime(time4)
	// start5 := primitive.NewDateTimeFromTime(time5)
	// start6 := primitive.NewDateTimeFromTime(time6)

	docs := []interface{}{
		bson.D{{Key: "cat_num", Value: "M349800"}, {Key: "cat_lot", Value: "210999"}, {Key: "raw_pn", Value: "15042"}, {Key: "raw_desc", Value: "prod1"}, {Key: "qty", Value: "100"}, {Key: "weight", Value: 0.5}, {Key: "status", Value: "wip"}, {Key: "uid", Value: "X202202"}},
		bson.D{{Key: "cat_num", Value: "M349801"}, {Key: "cat_lot", Value: "210945"}, {Key: "raw_pn", Value: "15743"}, {Key: "raw_desc", Value: "prod2"}, {Key: "qty", Value: "500"}, {Key: "weight", Value: 0.75}, {Key: "status", Value: "wip"}, {Key: "uid", Value: "X202203"}},
		bson.D{{Key: "cat_num", Value: "M350311"}, {Key: "cat_lot", Value: "211000"}, {Key: "raw_pn", Value: "16210"}, {Key: "raw_desc", Value: "prod3"}, {Key: "qty", Value: "25"}, {Key: "weight", Value: 1}, {Key: "status", Value: "wip"}, {Key: "uid", Value: "X202204"}},
	}
	_, err := m.DB.Collection("Jobs").InsertMany(context.TODO(), docs)

	if err != nil {
		panic(err)
	}

	// res.Header().Set("Content-Type", "application/json")
	// res.WriteHeader(http.StatusOK)
	// res.Write(payload)
}

func (m *Repository) ReadJob(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Reading job with ID...")
}

func (m *Repository) UpdateJob(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Updating job with ID...")
}

func (m *Repository) DeleteJob(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Deleting job with ID...")
}

func (m *Repository) UpdateJobStatus(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Job was drag and dropped into new bucket. Updating the status")
}

//*****************************Admin operations**************************************************

func (m *Repository) CreateUser(res http.ResponseWriter, req *http.Request) {
	request, _ := io.ReadAll(req.Body) //sent from axios post request. Json data and not form data

	var user models.User

	err := json.Unmarshal(request, &user) //parses JSON-encoded data[]bytes into user fields
	if err != nil {
		http.Error(res, "data not available. Please submit again", 400)
	}

	//context created and provided to MongoDB query funcs
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// filter the Employee DB and return user data where UID matches what was sent from frontend
	filter := bson.D{{Key: "uid", Value: user.UID}}
	if err = m.DB.Collection("Employee").FindOne(ctx, filter).Decode(&user); err != nil {
		http.Error(res, "employee not inserted", http.StatusBadRequest)
		return
	}

	//when i create the frontend page for this route, will send role along with UID.
	user.Status = true
	user.Role = "user"

	//since user is an Employee, the credentials get stored in the User Collection
	_, err = m.DB.Collection("User").InsertOne(ctx, &user)
	if err != nil {
		http.Error(res, "employee not inserted into user DB", http.StatusBadRequest)
		return
	}
	//need to send the user's role to front in order to direct the routing
	response, _ := json.Marshal(&user.Role)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (m *Repository) SetToInactive(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Admin is changing user status to inactive")
}

func (m *Repository) CreateProductInfo(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Admin is adding new product info to system")
}

//needed only to generate an employee in the database. won't exist in production
// bs, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
// if err != nil {
// 	http.Error(res, "having trouble receiving the data", http.StatusBadRequest)
// 	return
// }

// pass := string(bs)

// _, err = m.DB.Collection("Employee").InsertOne(ctx, models.Employee{
// 	UID:       user.UID,
// 	Firstname: "Employee",
// 	Lastname:  "One",
// 	Password:  pass,
// 	Email:     "employeeone@dmsapp.io",
// })

// if err != nil {
// 	http.Error(res, "employee not inserted", 500)
// 	return
// }

// CountDocuments returns the number of documents in the collection. For a fast count of the documents in the
// collection, see the EstimatedDocumentCount method.

// The filter parameter must be a document and can be used to select which documents contribute to the count.
// m.DB.Collection().CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions)
