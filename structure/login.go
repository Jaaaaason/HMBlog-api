package structure

// Login used to bind POST request data for /admin/login
type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
