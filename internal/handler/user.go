package handler

import (
	"strconv"

	"server/internal/service"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname" example:"张三"`
	Email    string `json:"email" binding:"omitempty,email" example:"zhangsan@example.com"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3" example:"newuser"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
	Email    string `json:"email" binding:"required,email" example:"newuser@example.com"`
}

// UpdateStatusRequest 更新状态请求
type UpdateStatusRequest struct {
	Status int `json:"status" binding:"oneof=0 1" example:"1"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据ID获取用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.User}
// @Failure 400 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的用户ID")
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		response.Error(c, 404, "用户不存在")
		return
	}

	response.Success(c, user)
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Security BearerAuth
// @Success 200 {object} response.Response{data=ListUsersResponse}
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	users, total, err := h.userService.GetAll(page, size)
	if err != nil {
		response.Error(c, 500, "获取用户列表失败")
		return
	}

	response.Success(c, ListUsersResponse{
		List:  users,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "用户信息"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.User}
// @Failure 400 {object} response.Response
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	user, err := h.userService.Create(req.Username, req.Password, req.Email)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateUser 更新用户
// @Summary 更新用户信息
// @Description 更新用户的昵称和邮箱（本人或管理员）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body UpdateUserRequest true "更新信息"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.User}
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的用户ID")
		return
	}

	// 只允许本人或管理员修改
	currentUserID, _ := c.Get("userID")
	role, _ := c.Get("role")
	roleStr, _ := role.(string)
	if currentUserID.(uint) != uint(id) && roleStr != "admin" {
		response.Error(c, 403, "无权限修改此用户")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	user, err := h.userService.Update(uint(id), req.Nickname, req.Email)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, user)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 根据ID删除用户（仅管理员）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// 仅管理员可删除用户
	role, _ := c.Get("role")
	roleStr, _ := role.(string)
	if roleStr != "admin" {
		response.Error(c, 403, "无权限执行此操作")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的用户ID")
		return
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 启用或禁用用户（仅管理员）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body UpdateStatusRequest true "状态信息"
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/users/{id}/status [patch]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	// 仅管理员可更新状态
	role, _ := c.Get("role")
	roleStr, _ := role.(string)
	if roleStr != "admin" {
		response.Error(c, 403, "无权限执行此操作")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "无效的用户ID")
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	if err := h.userService.UpdateStatus(uint(id), req.Status); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetCurrentUser 获取当前用户
// @Summary 获取当前登录用户信息
// @Description 获取当前Token对应的用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.User}
// @Failure 401 {object} response.Response
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	user, err := h.userService.GetByID(userID.(uint))
	if err != nil {
		response.Error(c, 404, "用户不存在")
		return
	}

	response.Success(c, user)
}
