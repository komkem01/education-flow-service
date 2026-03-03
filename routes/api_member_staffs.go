package routes

import (
	"education-flow/app/modules"
	"education-flow/app/modules/entities/ent"

	"github.com/gin-gonic/gin"
)

func apiMemberStaffs(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/staffs/register", mod.Staff.Ctl.Register)

	protected := r.Group("")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff))

	registerCRUD(protected, "/staffs", mod.Staff.Ctl.List, mod.Staff.Ctl.Get, mod.Staff.Ctl.Create, mod.Staff.Ctl.Update, mod.Staff.Ctl.Delete)
	registerCRUD(protected, "/inventory-items", mod.InventoryItem.Ctl.List, mod.InventoryItem.Ctl.Get, mod.InventoryItem.Ctl.Create, mod.InventoryItem.Ctl.Update, mod.InventoryItem.Ctl.Delete)
	registerCRUD(protected, "/inventory-requests", mod.InventoryRequest.Ctl.List, mod.InventoryRequest.Ctl.Get, mod.InventoryRequest.Ctl.Create, mod.InventoryRequest.Ctl.Update, mod.InventoryRequest.Ctl.Delete)
	registerCRUD(protected, "/document-tracking", mod.DocumentTracking.Ctl.List, mod.DocumentTracking.Ctl.Get, mod.DocumentTracking.Ctl.Create, mod.DocumentTracking.Ctl.Update, mod.DocumentTracking.Ctl.Delete)
	registerCRUD(protected, "/school-announcements", mod.SchoolAnnouncement.Ctl.List, mod.SchoolAnnouncement.Ctl.Get, mod.SchoolAnnouncement.Ctl.Create, mod.SchoolAnnouncement.Ctl.Update, mod.SchoolAnnouncement.Ctl.Delete)
	registerCRUD(protected, "/storages", mod.Storage.Ctl.List, mod.Storage.Ctl.Get, mod.Storage.Ctl.Create, mod.Storage.Ctl.Update, mod.Storage.Ctl.Delete)
	registerCRUD(protected, "/storage-links", mod.StorageLink.Ctl.List, mod.StorageLink.Ctl.Get, mod.StorageLink.Ctl.Create, mod.StorageLink.Ctl.Update, mod.StorageLink.Ctl.Delete)
}
