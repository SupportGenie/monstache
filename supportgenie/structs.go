package supportgenie

type Company struct {
	CompanyEmail string `bson:"email"`
}

type User struct {
	User_name  string `bson:"name"`
	User_phone string `bson:"phone"`
	User_email string `bson:"email"`
}

type Ticket struct {
	Company Company
	User    User
	Ticket  map[string]interface{}
}
