package handler

import (
	"UserCrud/dto/request"
	"UserCrud/dto/response"
	"UserCrud/middleware"
	"UserCrud/model"
	"UserCrud/repository"
	"UserCrud/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"net/http"
	"reflect"
)

type userHandler struct {
	userService    service.UserService
	authMiddleware middleware.AuthMiddleware
}

func AddUserHandler(userService service.UserService, authMiddleware middleware.AuthMiddleware, r *gin.Engine) {
	u := &userHandler{
		userService:    userService,
		authMiddleware: authMiddleware,
	}
	userRoute := r.Group("/user")
	userRoute.POST("/register", u.register())
	userRoute.POST("/login", u.login())
	userRoute.PATCH("", u.authMiddleware.ValidateAndExtractJwt(), u.updateUser())
	userRoute.DELETE("", u.authMiddleware.ValidateAndExtractJwt(), u.deleteUser())
	userRoute.GET("/profile", u.authMiddleware.ValidateAndExtractJwt(), u.getUser())
}

func respondWithError(err error, c *gin.Context) {
	var uniqueErr *repository.UniqueConstraintError
	var validationErr validator.ValidationErrors
	switch {
	case errors.Is(err, service.ErrWrongPassword):
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
	case errors.As(err, &uniqueErr):
		c.JSON(http.StatusBadRequest, response.Response{Error: err.Error()})
	case errors.Is(err, repository.ErrUserNotFound):
		c.JSON(http.StatusNotFound, response.Response{Error: err.Error()})
	case errors.As(err, &validationErr):
		var errMessage string
		switch validationErr[0].Tag() {
		case "required":
			errMessage = "Field " + validationErr[0].Field() + " is required"
		case "email":
			errMessage = "Field " + validationErr[0].Field() + " must be a valid email"
		case "number":
			errMessage = "Field " + validationErr[0].Field() + " must be a valid number"
		}
		c.JSON(http.StatusBadRequest, response.Response{Error: errMessage})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
}

func (u *userHandler) register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.RegisterRequest
		if err := c.ShouldBind(&req); err != nil {
			respondWithError(err, c)
			return
		}
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			respondWithError(err, c)
			return
		}
		err := u.userService.Register(c, model.User{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       req.Email,
			Password:    req.Password,
			PhoneNumber: req.PhoneNumber,
			Gender:      req.Gender,
		})
		if err != nil {
			respondWithError(err, c)
			return
		}
		c.JSON(http.StatusCreated, response.Response{Message: "User register successfully"})
	}
}

func (u *userHandler) login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.LoginRequest
		if err := c.ShouldBind(&req); err != nil {
			respondWithError(err, c)
			return
		}
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			respondWithError(err, c)
			return
		}
		token, err := u.userService.Login(c, req.Email, req.Password)
		if err != nil {
			respondWithError(err, c)
			return
		}
		c.JSON(http.StatusOK, response.LoginResponse{AccessToken: token})
	}
}

func (u *userHandler) updateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.UpdateRequest
		if err := c.ShouldBind(&req); err != nil {
			respondWithError(err, c)
			return
		}
		if !reflect.ValueOf(req.PhoneNumber).IsZero() {
			validate := validator.New()
			if err := validate.StructPartial(req, "PhoneNumber"); err != nil {
				respondWithError(err, c)
				return
			}
		}
		claims, _ := c.Get(middleware.JWTClaimsContextKey)
		userClaims := claims.(jwt.MapClaims)
		id := uint(userClaims["userId"].(float64))
		err := u.userService.Update(c, model.User{
			ID:          id,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Password:    req.Password,
			PhoneNumber: req.PhoneNumber,
			Gender:      req.Gender,
		})
		if err != nil {
			respondWithError(err, c)
			return
		}
		c.JSON(http.StatusCreated, response.Response{Message: "User update successfully"})
	}
}

func (u *userHandler) deleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get(middleware.JWTClaimsContextKey)
		userClaims := claims.(jwt.MapClaims)
		id := uint(userClaims["userId"].(float64))
		err := u.userService.Delete(c, id)
		if err != nil {
			respondWithError(err, c)
			return
		}
		c.JSON(http.StatusOK, response.Response{Message: "User delete successfully"})
	}
}

func (u *userHandler) getUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get(middleware.JWTClaimsContextKey)
		userClaims := claims.(jwt.MapClaims)
		id := uint(userClaims["userId"].(float64))
		user, err := u.userService.GetById(c, id)
		if err != nil {
			respondWithError(err, c)
			return
		}
		c.JSON(http.StatusOK, response.GetUserResponse{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Gender:      user.Gender,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}
}
