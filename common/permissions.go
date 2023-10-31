package common

var Permissions = []string{
	"GET/console.post.index",
	"GET/console.post.create",
	"POST/console.post.store",
	"GET/console.post.edit",
	"PUT/console.post.update",
	"DELETE/console.post.destroy",
	"GET/console.post.trash",
	"POST/console.post.unTrash",
	"POST/console.post.imgUpload",
	"GET/console.cate.index",
	"GET/console.cate.edit",
	"PUT/console.cate.update",
	"POST/console.cate.store",
	"DELETE/console.cate.destroy",
	"GET/console.tag.index",
	"POST/console.tag.store",
	"GET/console.tag.edit",
	"PUT/console.tag.update",
	"DELETE/console.tag.destroy",
	"GET/console.system.index",
	"PUT/console.system.update",
	"GET/console.link.index",
	"POST/console.link.store",
	"GET/console.link.edit",
	"PUT/console.link.update",
	"DELETE/console.link.destroy",
	"DELETE/console.auth.logout",
	"GET/console.home.index",
	"DELETE/console.auth.cache",
}

func CheckPermissions(permission string, method string) bool {
	for _, v := range Permissions {
		if v == method+"/"+permission {
			return true
		}
	}
	return false
}
