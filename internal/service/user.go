package service

import (
	"errors"

	"server/internal/model"
	"server/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Create 创建用户
func (s *UserService) Create(username, password, email string) (*model.User, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	user := &model.User{
		Username: username,
		Password: password,
		Email:    email,
		Status:   1,
	}

	// 统一使用 model 的密码加密方法
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

// GetAll 获取所有用户
func (s *UserService) GetAll(page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.userRepo.FindAll(page, pageSize)
}

// Update 更新用户
func (s *UserService) Update(id uint, nickname, email string) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	user.Nickname = nickname
	user.Email = email

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	if _, err := s.userRepo.FindByID(id); err != nil {
		return errors.New("用户不存在")
	}
	return s.userRepo.Delete(id)
}

// UpdateStatus 更新用户状态
func (s *UserService) UpdateStatus(id uint, status int) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.Status = status
	return s.userRepo.Update(user)
}
