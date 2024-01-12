package models

// Credentials provided by the user during sign-up and sign-in
type RawCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct for data saved in the db
type StoredCredentials struct {
	Username string `gorm:"primaryKey"`
	Salt     int
	Hash     []byte
}
