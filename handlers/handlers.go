package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	response, err := json.Marshal(&user)
	if err != nil {
		http.Error(res, "Unable to send user info", http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (m *Repository) Dashboard(res http.ResponseWriter, req *http.Request) {

	filter := bson.D{{Key: "email", Value: `usertest@gmail.com`}}

	var result models.User
	if err := m.DB.Collection("User").FindOne(context.TODO(), filter).Decode(&result); err != nil {
		panic(err)
	}

	//data, err := bson.Marshal(result)
	marsh, err := json.Marshal(result)

	if err != nil {
		log.Println("oppsy someone had a poopsie")
	}
	fmt.Fprintln(res, string(marsh))

}

func (m *Repository) UserProfile(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "This is the user profile page. Specific to the user with a UID")
}

func (m *Repository) CreateJob(res http.ResponseWriter, req *http.Request) {
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "request has been denied", 404)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(payload)
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

// Admin operations
func (m *Repository) CreateUser(res http.ResponseWriter, req *http.Request) {
	request, _ := io.ReadAll(req.Body)

	var user models.User

	err := json.Unmarshal(request, &user)
	if err != nil {
		http.Error(res, "data not available. Please submit again", 400)
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

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// filter the Employee DB and return user data where UID matches what was sent from frontend
	filter := bson.D{{Key: "uid", Value: user.UID}}
	if err = m.DB.Collection("Employee").FindOne(ctx, filter).Decode(&user); err != nil {
		http.Error(res, "employee not inserted", http.StatusBadRequest)
		return
	}

	user.Status = true
	user.Role = "user"

	_, err = m.DB.Collection("User").InsertOne(ctx, &user)
	if err != nil {
		http.Error(res, "employee not inserted into user DB", http.StatusBadRequest)
		return
	}
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

// //use this context to disconnect from mongo

// _, err := m.DB.Collection("User").InsertOne(ctx, models.User{
// 	UID:      "X345904",
// 	FName:    "User3",
// 	LName:    "Testing3",
// 	Password: "password3",
// 	Email:    "usertest@gmail.com",
// 	Status:   false,
// 	Role:     "user",
// })

// if err != nil {
// 	io.WriteString(res, "Document not inserted")
// }
// io.WriteString(res, "You are now on the login page")
