package routes

import (
	"net/http"

	"github.com/gary-stroup-developer/bkend-dms/handlers"
)

func Routes(r *http.ServeMux) {
	r.HandleFunc("/dms/login", handlers.Repo.Login)
	r.HandleFunc("/dms/dashboard", handlers.Repo.Dashboard)
	r.HandleFunc("/dms/userprofile/", handlers.Repo.UserProfile)
	r.HandleFunc("/dms/create-job", handlers.Repo.CreateJob)
	r.HandleFunc("/dms/read-job/", handlers.Repo.ReadJob)
	r.HandleFunc("/dms/update-job/", handlers.Repo.UpdateJob)
	r.HandleFunc("/dms/delete-job/", handlers.Repo.DeleteJob)
	r.HandleFunc("/dms/update/job-status/", handlers.Repo.UpdateJobStatus)
	r.HandleFunc("/dms/user/create", handlers.Repo.CreateUser)
	r.HandleFunc("/dms/user/inactive/", handlers.Repo.SetToInactive)
	r.HandleFunc("/dms/product/create", handlers.Repo.CreateProductInfo)
}