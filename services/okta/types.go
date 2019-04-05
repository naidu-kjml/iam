package okta

// User contains formatted user data provided by Okta
type User struct {
	EmployeeNumber string   `json:"employeeNumber"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Position       string   `json:"position"`
	Department     string   `json:"department"`
	Email          string   `json:"email"`
	Location       string   `json:"location"`
	IsVendor       bool     `json:"isVendor"`
	TeamMembership []string `json:"teamMembership"`
	Manager        string   `json:"manager"`
}
