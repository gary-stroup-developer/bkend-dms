package main

import "net/http"

func Routes(r *http.ServeMux) {
	r.HandleFunc("/dms/login", Repo.Login)
	r.HandleFunc("/dms/dashboard", Repo.Dashboard)
	r.HandleFunc("/dms/userprofile/", Repo.UserProfile)
	r.HandleFunc("/dms/create-job/", Repo.CreateJob)
	r.HandleFunc("/dms/read-job/", Repo.ReadJob)
	r.HandleFunc("/dms/update-job/", Repo.UpdateJob)
	r.HandleFunc("/dms/delete-job/", Repo.DeleteJob)
	r.HandleFunc("/dms/update/job-status/", Repo.UpdateJobStatus)
	r.HandleFunc("/dms/user/create", Repo.CreateUser)
	r.HandleFunc("/dms/user/inactive/", Repo.SetToInactive)
	r.HandleFunc("/dms/product/create", Repo.CreateProductInfo)
}
