package service
import (
	"github.com/google/uuid"
)

type UserService struct {
	
}
func NewUserService() *UserService {
	return &UserService{}
}

// RegisterUser creates a new user and returns a unique userID
func (s *UserService) RegisterUser() string {
	// Generate a new unique user ID
	userID := uuid.New().String()
	return userID
}
