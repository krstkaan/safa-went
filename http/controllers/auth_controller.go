package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"safa-went/database/models"
	"safa-went/http/middlewares"
	"safa-went/http/requests"
	"safa-went/http/resources"
	"safa-went/internal/responses"
)

type Auth struct {
    DB *gorm.DB
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body requests.RegisterPayload true "Register payload"
// @Success 201 {object} resources.AuthResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 409 {object} responses.ErrorBody
// @Failure 422 {object} responses.ErrorBody
// @Router /auth/register [post]
func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
    var payload requests.RegisterPayload
    if err := render.DecodeJSON(r.Body, &payload); err != nil {
        responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
        return
    }
    if payload.Name == "" || payload.Email == "" || payload.Password == "" {
        responses.JSONError(w, r, http.StatusUnprocessableEntity, "name, email and password are required")
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
    if err != nil {
        responses.JSONError(w, r, http.StatusInternalServerError, "could not hash password")
        return
    }

    user := models.User{Name: payload.Name, Email: payload.Email, Password: string(hashed)}
    if err := a.DB.Create(&user).Error; err != nil {
        responses.JSONError(w, r, http.StatusConflict, "email already in use")
        return
    }

    token, err := generateToken(user.ID)
    if err != nil {
        responses.JSONError(w, r, http.StatusInternalServerError, "could not generate token")
        return
    }

    render.Status(r, http.StatusCreated)
    render.JSON(w, r, resources.AuthResource{Token: token, User: resources.NewAuthUserResource(user)})
}

// Login godoc
// @Summary Login
// @Description Authenticate with email and password, returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body requests.LoginPayload true "Login payload"
// @Success 200 {object} resources.AuthResource
// @Failure 400 {object} responses.ErrorBody
// @Failure 401 {object} responses.ErrorBody
// @Router /auth/login [post]
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
    var payload requests.LoginPayload
    if err := render.DecodeJSON(r.Body, &payload); err != nil {
        responses.JSONError(w, r, http.StatusBadRequest, "invalid request body")
        return
    }

    var user models.User
    if err := a.DB.Where("email = ?", payload.Email).First(&user).Error; err != nil {
        responses.JSONError(w, r, http.StatusUnauthorized, "invalid credentials")
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
        responses.JSONError(w, r, http.StatusUnauthorized, "invalid credentials")
        return
    }

    token, err := generateToken(user.ID)
    if err != nil {
        responses.JSONError(w, r, http.StatusInternalServerError, "could not generate token")
        return
    }

    render.JSON(w, r, resources.AuthResource{Token: token, User: resources.NewAuthUserResource(user)})
}

// Me godoc
// @Summary Get current user
// @Description Returns the authenticated user's information
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} resources.AuthUserResource
// @Failure 401 {object} responses.ErrorBody
// @Router /auth/user [get]
func (a *Auth) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.UserContextKey).(*models.User)
	if !ok {
		responses.JSONError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	render.JSON(w, r, resources.NewAuthUserResource(*user))
}

// Logout godoc
// @Summary Logout
// @Description Invalidate the current session (client should discard the token)
// @Tags Auth
// @Security BearerAuth
// @Success 204
// @Router /auth/logout [post]
func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func generateToken(userID uint) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
        "iat": time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}