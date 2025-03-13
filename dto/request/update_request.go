package request

type UpdateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number" validate:"number"`
	Gender      string `json:"gender"`
}
