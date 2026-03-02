package routes

import "github.com/gin-gonic/gin"

func registerCRUD(r *gin.RouterGroup, path string, list gin.HandlerFunc, get gin.HandlerFunc, create gin.HandlerFunc, update gin.HandlerFunc, del gin.HandlerFunc) {
	r.GET(path, list)
	r.GET(path+"/:id", get)
	r.POST(path, create)

	if update != nil {
		r.PATCH(path+"/:id", update)
	}

	if del != nil {
		r.DELETE(path+"/:id", del)
	}
}
