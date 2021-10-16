package request

type UserDetailsUpdate struct {
	Name     string `json:"name" binding:"omitempty,min=0,max=100"`
	About    string `json:"about" binding:"omitempty,min=0,max=1000"`
	Username string `json:"username" binding:"omitempty,min=0,max=100"`
}
