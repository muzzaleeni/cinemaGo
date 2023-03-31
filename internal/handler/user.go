package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/forChin/my-project/internal/handler/request"
	"github.com/forChin/my-project/internal/handler/response"
	"github.com/forChin/my-project/internal/model"
	"github.com/forChin/my-project/internal/service"
	"github.com/forChin/my-project/pkg/liberror"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var defaultPageSize uint64 = 500

type UserHandler struct {
	userService *service.UserService
	logger      *zap.SugaredLogger
}

func NewUserHandler(
	userService *service.UserService,
	logger *zap.SugaredLogger,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req request.UserCreate
	if err := c.BodyParser(&req); err != nil {
		return newInvalidJSONErr(err)
	}

	req.Normalise()

	if err := req.Validate(); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	user := model.User{
		Name: req.Name,
	}
	created, err := h.userService.Create(c.Context(), user)
	if err != nil {
		h.logger.Errorw("create user",
			zap.Error(err), zap.Any("user", user))
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusCreated).JSON(created)
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	var filter model.UserFilter
	if err := c.QueryParser(&filter); err != nil {
		return newQueryParamErr(err)
	}

	page := c.Query("page", "1")
	pageNum, err := strconv.ParseUint(page, 10, 64)
	if err != nil {
		return newQueryParamErr(errors.New(`"page" expected to be a number`))
	}

	offset := (pageNum - 1) * defaultPageSize
	users, totalRows, err := h.userService.GetAll(c.Context(), defaultPageSize, offset, filter)
	if err != nil {
		h.logger.Errorf("get all users: %v", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	pageResp := response.NewPage(
		pageNum, defaultPageSize, totalRows, len(users), users,
	)
	return c.Status(http.StatusOK).JSON(pageResp)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return newRouteParamErr(`"id" must be an integer`)
	}

	err = h.userService.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, liberror.ErrNotFound) {
			return newNotFoundErr("user", id)
		}

		h.logger.Errorw("delete user",
			zap.Error(err), zap.Int("id", id))
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusNoContent).Send(nil)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return newRouteParamErr(`"id" must be an integer`)
	}

	var req request.UserUpdate
	if err := c.BodyParser(&req); err != nil {
		return newInvalidJSONErr(err)
	}

	req.Normalise()

	if err := req.Validate(); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	user := model.User{
		Name: req.Name,
	}
	err = h.userService.Update(c.Context(), user)
	if err != nil {
		if errors.Is(err, liberror.ErrNotFound) {
			return newNotFoundErr("user", id)
		}

		h.logger.Errorw("update user",
			zap.Error(err), zap.Any("user", user))
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusNoContent).Send(nil)
}
