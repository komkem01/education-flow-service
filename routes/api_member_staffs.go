package routes

import (
	"education-flow/app/modules"

	"github.com/gin-gonic/gin"
)

func apiMemberStaffs(r *gin.RouterGroup, mod *modules.Modules) {
	registerCRUD(r, "/staffs", mod.Staff.Ctl.List, mod.Staff.Ctl.Get, mod.Staff.Ctl.Create, mod.Staff.Ctl.Update, mod.Staff.Ctl.Delete)
	registerCRUD(r, "/inventory-items", mod.InventoryItem.Ctl.List, mod.InventoryItem.Ctl.Get, mod.InventoryItem.Ctl.Create, mod.InventoryItem.Ctl.Update, mod.InventoryItem.Ctl.Delete)
	registerCRUD(r, "/inventory-requests", mod.InventoryRequest.Ctl.List, mod.InventoryRequest.Ctl.Get, mod.InventoryRequest.Ctl.Create, mod.InventoryRequest.Ctl.Update, mod.InventoryRequest.Ctl.Delete)
	registerCRUD(r, "/document-tracking", mod.DocumentTracking.Ctl.List, mod.DocumentTracking.Ctl.Get, mod.DocumentTracking.Ctl.Create, mod.DocumentTracking.Ctl.Update, mod.DocumentTracking.Ctl.Delete)
	registerCRUD(r, "/school-announcements", mod.SchoolAnnouncement.Ctl.List, mod.SchoolAnnouncement.Ctl.Get, mod.SchoolAnnouncement.Ctl.Create, mod.SchoolAnnouncement.Ctl.Update, mod.SchoolAnnouncement.Ctl.Delete)
	registerCRUD(r, "/storages", mod.Storage.Ctl.List, mod.Storage.Ctl.Get, mod.Storage.Ctl.Create, mod.Storage.Ctl.Update, mod.Storage.Ctl.Delete)
	registerCRUD(r, "/storage-links", mod.StorageLink.Ctl.List, mod.StorageLink.Ctl.Get, mod.StorageLink.Ctl.Create, mod.StorageLink.Ctl.Update, mod.StorageLink.Ctl.Delete)
}
