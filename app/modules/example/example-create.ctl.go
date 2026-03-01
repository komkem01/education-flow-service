package example

import (
	"education-flow/app/utils"
	"education-flow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CreateRequest struct {
	Filename   string `json:"filename" binding:"required,filename"`
	Visibility string `json:"visibility" default:"private"`
}

type CreateResponse struct {
	ID string `json:"id"`
}

func (c *Controller) Create(ctx *gin.Context) {
	var req CreateRequest
	var res CreateResponse
	span, log := utils.LogSpanFromGin(ctx)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}
	span.AddEvent("example.create.request", trace.WithAttributes(
		attribute.String("filename", req.Filename),
		attribute.String("visibility", req.Visibility),
	))

	userID, err := uuid.Parse(ctx.GetString("userID"))
	if err != nil {
		userID = uuid.Nil
	}
	example, err := c.svc.Create(ctx, userID)
	if err != nil {
		log.Errf("example.create.error: %s", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	copier.CopyWithOption(
		&res,
		example,
		copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
		},
	)
	base.Success(ctx, res, "example-created")
}
