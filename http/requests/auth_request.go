package requests

type RegisterPayload struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginPayload struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}