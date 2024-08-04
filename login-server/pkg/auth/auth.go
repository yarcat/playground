package auth

type (
	Email  string
	UserID string

	Auth struct {
		Email
		UserID
	}
)

func IsValid(a *Auth) bool {
	return a != nil && a.Email != "" && a.UserID != ""
}

func New(email Email, userID UserID) *Auth {
	return &Auth{
		Email:  email,
		UserID: userID,
	}
}

func (a *Auth) GetEmail() Email {
	if a == nil {
		return ""
	}
	return a.Email
}

func (a *Auth) GetUserID() UserID {
	if a == nil {
		return ""
	}
	return a.UserID
}
