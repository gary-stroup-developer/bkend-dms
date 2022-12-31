package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gary-stroup-developer/bkend-dms/models"
	uuid "github.com/satori/go.uuid"
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
		return
	}

	var user models.User
	var userPayload models.User

	err = json.Unmarshal(payload, &userPayload)
	if err != nil {
		http.Error(res, "Sorry. There seems to be an issue connecting to the database", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filter := bson.D{{Key: "uid", Value: userPayload.UID}}
	if err = m.DB.Collection("User").FindOne(ctx, filter).Decode(&user); err != nil {
		http.Error(res, "Sorry. User information unavailable at this time", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPayload.Password))
	if err != nil {
		http.Error(res, "username and/or password do not match", http.StatusBadRequest)
		return
	}

	//store the user info in the repository to be accessed at other routes
	m.UserInfo = user

	response, err := json.Marshal(struct {
		Firstname string
		Lastname  string
		Role      string
	}{Firstname: user.Firstname, Lastname: user.Lastname, Role: user.Role})
	if err != nil {
		http.Error(res, "Unable to send user info", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (m *Repository) Dashboard(res http.ResponseWriter, req *http.Request) {
	// create the stages to get users who are active, then return everything except the password & status
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: true}, {Key: "role", Value: "user"}}}}
	unsetStage := bson.D{{Key: "$unset", Value: bson.A{"password", "status"}}}

	var users []bson.M //will hold the user data retrieved from the database

	//use aggregate func to query results using pipeline
	cursor, err := m.DB.Collection("User").Aggregate(context.TODO(), mongo.Pipeline{matchStage, unsetStage})
	if err != nil {
		http.Error(res, "trouble connecting to server", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &users) //store all active users in database

	payload, _ := json.Marshal(&users)
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func (m *Repository) UserProfile(res http.ResponseWriter, req *http.Request) {
	var id models.User
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Trouble retrieving information", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(payload, &id); err != nil {
		http.Error(res, "user profile unmarshal", http.StatusBadRequest)
		return
	}
	log.Println(id)
	var user []bson.M //will hold the user data retrieved from the database
	filter := bson.D{{Key: "uid", Value: id.UID}}
	//search for user with the id passed to the payload
	err = m.DB.Collection("User").FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		http.Error(res, "trouble connecting to server", http.StatusInternalServerError)
		return
	}

	var jobs []bson.M //stores the results of query. Advantage is that this wont throw errors as may decoding into struct might

	//create stages for job query
	stageOne := bson.D{{Key: "$match", Value: bson.D{{Key: "uid", Value: id}}}}
	stageTwo := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "jobs", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}},
	}

	c, err := m.DB.Collection("Jobs").Aggregate(context.TODO(), mongo.Pipeline{stageOne, stageTwo})
	if err != nil {
		http.Error(res, "no jobs found for user", http.StatusInternalServerError)
		return
	}
	defer c.Close(context.TODO())

	if err = c.All(context.TODO(), &jobs); err != nil {
		http.Error(res, "database connection is faulty", http.StatusInternalServerError)
		return
	}
	var profile models.DashboardResponse
	profile.Users = user
	profile.Jobs = jobs
	response, _ := json.Marshal(&profile)

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (m *Repository) CreateJob(res http.ResponseWriter, req *http.Request) {
	//read in the data sent by axios request
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "request has been denied", 404)
	}

	var job models.Job
	jobID := uuid.NewV4().String()
	//store data into the job variable of type Job
	if err = json.Unmarshal(payload, &job); err != nil {
		http.Error(res, "could not parse data", http.StatusBadRequest)
		return
	}
	job.ID = jobID
	job.UID = m.UserInfo.UID

	_, err = m.DB.Collection("Jobs").InsertOne(context.TODO(), job) //insert the job into the Jobs Collection
	if err != nil {
		http.Error(res, "The job was not saved! Please resubmit", http.StatusInternalServerError)
		return
	}

	//extract weight of job and increment the capacity field of user by the weight value
	weight := job.Weight
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "uid", Value: job.UID}}}}
	updateStage := bson.D{{Key: "$inc", Value: bson.D{{Key: "capacity", Value: weight}}}}
	//run pipeline
	_, err = m.DB.Collection("User").Aggregate(context.TODO(), mongo.Pipeline{matchStage, updateStage})
	if err != nil {
		http.Error(res, "could not update the user capacity", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("The job was successfully created. User capacity has been updated"))
}

func (m *Repository) ReadJob(res http.ResponseWriter, req *http.Request) {
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "trouble parsing data", http.StatusBadRequest)
		return
	}

	var id string
	json.Unmarshal(payload, &id)
	//create stages for job query
	filter := bson.D{{Key: "uid", Value: m.UserInfo.UID}, {Key: "id", Value: id}}

	var job []bson.M
	//use aggregate func to query results using pipeline
	c := m.DB.Collection("Jobs").FindOne(context.TODO(), filter).Decode(&job)
	if c.Error() != "" {
		http.Error(res, "trouble connecting to server", http.StatusInternalServerError)
	}

	response, err := json.Marshal(&job)
	if err != nil {
		http.Error(res, "trouble parsing data", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (m *Repository) UpdateJob(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Updating job with ID...")
}

func (m *Repository) DeleteJob(res http.ResponseWriter, req *http.Request) {
	//read in the data sent by axios request
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "request has been denied", 404)
	}

	var job models.Job

	//store data into the job variable of type Job
	if err = json.Unmarshal(payload, &job); err != nil {
		http.Error(res, "could not parse data", http.StatusBadRequest)
		return
	}

	filter := bson.D{{Key: "id", Value: job.ID}}
	_, err = m.DB.Collection("Jobs").DeleteOne(context.TODO(), filter) //insert the job into the Jobs Collection
	if err != nil {
		http.Error(res, "The job was not deleted", http.StatusInternalServerError)
		return
	}

	//extract weight of job and increment the capacity field of user by the weight value
	weight := job.Weight
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "uid", Value: job.UID}}}}
	updateStage := bson.D{{Key: "$inc", Value: bson.D{{Key: "capacity", Value: -weight}}}}
	//run pipeline
	_, err = m.DB.Collection("User").Aggregate(context.TODO(), mongo.Pipeline{matchStage, updateStage})
	if err != nil {
		http.Error(res, "could not update the user capacity", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("The job was deleted and the user capacity has been updated!"))
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

// docs := []interface{}{
// 	bson.D{{Key: "cat_num", Value: "M349800"}, {Key: "cat_lot", Value: "210999"}, {Key: "raw_pn", Value: "15042"}, {Key: "raw_desc", Value: "prod1"}, {Key: "qty", Value: "100"}, {Key: "weight", Value: 0.5}, {Key: "status", Value: "wip"}, {Key: "uid", Value: "X202202"}},
// 	bson.D{{Key: "cat_num", Value: "M349801"}, {Key: "cat_lot", Value: "210945"}, {Key: "raw_pn", Value: "15743"}, {Key: "raw_desc", Value: "prod2"}, {Key: "qty", Value: "500"}, {Key: "weight", Value: 0.75}, {Key: "status", Value: "wip"}, {Key: "uid", Value: "X202203"}},
// 	bson.D{{Key: "cat_num", Value: "M350311"}, {Key: "cat_lot", Value: "211000"}, {Key: "raw_pn", Value: "16210"}, {Key: "raw_desc", Value: "prod3"}, {Key: "qty", Value: "25"}, {Key: "weight", Value: 1}, {Key: "status", Value: "wip"}, {Key: "uid", Value: "X202204"}},
// }
