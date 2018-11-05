package structure

// User the blog user struct
type User struct {
	Username     string `json:"username" bson:"username"`
	PasswordHash []byte `bson:"password_hash"`
}
