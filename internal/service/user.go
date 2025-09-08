package service

import (
	"github.com/ray-d-song/go-echo-monolithic/internal/model"
	"github.com/ray-d-song/go-echo-monolithic/internal/repository"
	"github.com/ray-d-song/go-echo-monolithic/internal/types"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user business logic
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// GetByUsername retrieves a user by username
func (s *UserService) GetByUsername(username string) (*model.User, error) {
	return s.userRepo.GetByUsername(username)
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*model.User, error) {
	return s.userRepo.GetByEmail(email)
}

// Update updates user information
func (s *UserService) Update(userID uint, req *types.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Email != nil {
		// Check if email is already taken by another user
		existingUser, err := s.userRepo.GetByEmail(*req.Email)
		if err == nil && existingUser.ID != userID {
			return nil, types.ErrUserAlreadyExists
		}
		user.Email = *req.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete soft deletes a user
func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

// List retrieves users with pagination
func (s *UserService) List(offset, limit int) ([]*model.User, error) {
	return s.userRepo.List(offset, limit)
}

// Count returns total number of users
func (s *UserService) Count() (int64, error) {
	return s.userRepo.Count()
}

// ToResponse converts User model to UserResponse
func (s *UserService) ToResponse(user *model.User) *types.UserResponse {
	return &types.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// hashPassword hashes a password using bcrypt
func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// verifyPassword verifies a password against its hash
func (s *UserService) verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}