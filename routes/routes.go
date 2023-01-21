package routes

import (
	"net/http"

	"github.com/gary-stroup-developer/bkend-dms/handlers"
)

func Routes(r *http.ServeMux) {
	r.HandleFunc("/dms/login", handlers.Repo.Login)
	r.HandleFunc("/dms/dashboard", handlers.Repo.Dashboard)
	r.HandleFunc("/dms/userprofile/jobs", handlers.Repo.UserProfile)
	r.HandleFunc("/dms/job/create", handlers.Repo.CreateJob)
	r.HandleFunc("/dms/job/create/search", handlers.Repo.SearchJob)
	r.HandleFunc("/dms/job/update", handlers.Repo.UpdateJob)
	r.HandleFunc("/dms/job/delete", handlers.Repo.DeleteJob)
	r.HandleFunc("/dms/update/job-status", handlers.Repo.UpdateJobStatus)
	r.HandleFunc("/dms/job/update/fsr", handlers.Repo.FSRHandler)
	r.HandleFunc("/dms/job/update/cbr", handlers.Repo.CBRHandler)
	r.HandleFunc("/dms/user/create", handlers.Repo.CreateUser)
	r.HandleFunc("/dms/user/inactive/", handlers.Repo.SetToInactive)
	r.HandleFunc("/dms/product/create", handlers.Repo.CreateProductInfo)

	r.Handle("/favicon.ico", http.NotFoundHandler())

}
