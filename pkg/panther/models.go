package panther

type User struct {
	ID         string `json:"id"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Email      string `json:"email"`
	Status     string `json:"status"`
	Role       Role   `json:"role"`
}

type Role struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
