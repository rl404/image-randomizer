package internal

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

// RegisterBaseRoutes registers base routes.
func RegisterBaseRoutes(router *chi.Mux) {
	// Root route for testing.
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, "root", nil)
	})

	// Ping route also for testing.
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, "pong", nil)
	})

	// Handle page not found 404.
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusNotFound, "page not found", nil)
	}))
}

// GetUserRoutes to get all user routes.
func GetUserRoutes(cfg Config) (r http.Handler, err error) {
	router := chi.NewRouter()

	uh, err := newUserHandler(cfg)
	if err != nil {
		return router, err
	}

	router.Get("/{username}", uh.list)
	router.Get("/{username}/image.jpg", uh.random)
	router.Post("/register", uh.register)
	router.Post("/login", uh.login)
	router.Post("/update", uh.update)

	return router, nil
}

// UserHandler to handle all user activities.
type UserHandler struct {
	DB     *gorm.DB
	Config Config
}

// newUserHandler to create new instance.
func newUserHandler(cfg Config) (uh UserHandler, err error) {
	if cfg.Masterkey == "" {
		return uh, ErrRequiredKey
	}
	uh.Config = cfg

	// Init db connection.
	uh.DB, err = cfg.InitDB()
	return uh, err
}

// register to register new user.
func (uh *UserHandler) register(w http.ResponseWriter, r *http.Request) {
	var request User

	// Get request body.
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Required username.
	if request.Username == "" {
		RespondWithJSON(w, http.StatusBadRequest, ErrRequiredUser.Error(), nil)
		return
	}

	// Required password.
	if request.Password == "" {
		RespondWithJSON(w, http.StatusBadRequest, ErrRequiredPass.Error(), nil)
		return
	}

	// Is username exist.
	var tmp User
	uh.DB.Where("username = ?", Encrypt(request.Username, uh.Config.Masterkey)).First(&tmp)
	if tmp.ID != 0 {
		RespondWithJSON(w, http.StatusBadRequest, ErrUserExist.Error(), nil)
		return
	}

	// Insert new user.
	newUser := User{
		Username: Encrypt(request.Username, uh.Config.Masterkey),
		Password: Encrypt(request.Password, uh.Config.Masterkey),
	}

	err = uh.DB.Create(&newUser).Error
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
	} else {
		RespondWithJSON(w, http.StatusCreated, http.StatusText(http.StatusCreated), nil)
	}
}

// login to login and get user token.
func (uh *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	var request User

	// Get request body.
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Required username.
	if request.Username == "" {
		RespondWithJSON(w, http.StatusBadRequest, ErrRequiredUser.Error(), nil)
		return
	}

	// Required password.
	if request.Password == "" {
		RespondWithJSON(w, http.StatusBadRequest, ErrRequiredPass.Error(), nil)
		return
	}

	// Get user.
	var user User
	uh.DB.Where("username = ?", Encrypt(request.Username, uh.Config.Masterkey)).First(&user)
	if user.ID == 0 {
		RespondWithJSON(w, http.StatusNotFound, ErrNotFound.Error(), nil)
		return
	}

	if user.Password != Encrypt(request.Password, uh.Config.Masterkey) {
		RespondWithJSON(w, http.StatusBadRequest, ErrWrongUserPass.Error(), nil)
		return
	}

	// Generate token.
	user.Token = RandomStr(16)
	user.TokenExpired = time.Now().AddDate(0, 0, 1)

	err = uh.DB.Save(&user).Error
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
	} else {
		RespondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK), user.Token)
	}
}

// update to update user's image list.
func (uh *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	// Header check.
	token := r.Header.Get("token")

	if token == "" {
		RespondWithJSON(w, http.StatusUnauthorized, ErrRequiredToken.Error(), nil)
		return
	}

	var user User
	uh.DB.Where("token = ?", token).First(&user)

	if user.ID == 0 {
		RespondWithJSON(w, http.StatusForbidden, ErrInvalidToken.Error(), nil)
		return
	}

	if user.TokenExpired.Before(time.Now()) {
		RespondWithJSON(w, http.StatusForbidden, ErrInvalidToken.Error(), nil)
		return
	}

	// Get request body.
	var request []string
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Delete existing data.
	uh.DB.Where("user_id = ?", user.ID).Delete(&Image{})

	if len(request) > 0 {
		// Prepare new data.
		var data []Image
		for _, i := range request {
			i = strings.TrimSpace(i)
			if i != "" {
				data = append(data, Image{
					UserID: user.ID,
					Image:  Encrypt(i, uh.Config.Masterkey),
				})
			}
		}

		err := uh.DB.Create(&data).Error
		if err != nil {
			RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	RespondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK), len(request))
}

// list to get user's image list.
func (uh *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	// Header check.
	token := r.Header.Get("token")

	if token == "" {
		RespondWithJSON(w, http.StatusUnauthorized, ErrRequiredToken.Error(), nil)
		return
	}

	var user User
	uh.DB.Where("token = ?", token).First(&user)

	if user.ID == 0 {
		RespondWithJSON(w, http.StatusForbidden, ErrInvalidToken.Error(), nil)
		return
	}

	if user.TokenExpired.Before(time.Now()) {
		RespondWithJSON(w, http.StatusForbidden, ErrInvalidToken.Error(), nil)
		return
	}

	var images []Image
	uh.DB.Where("user_id = ?", user.ID).Order("id asc").Find(&images)

	list := []string{}
	for _, i := range images {
		list = append(list, Decrypt(i.Image, uh.Config.Masterkey))
	}

	RespondWithJSON(w, http.StatusOK, http.StatusText(http.StatusOK), list)
}

// random to get random image from user's image list.
func (uh *UserHandler) random(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	var user User
	uh.DB.Where("username = ?", Encrypt(username, uh.Config.Masterkey)).First(&user)

	if user.ID == 0 {
		return
	}

	var images []Image
	uh.DB.Where("user_id = ?", user.ID).Find(&images)
	if len(images) == 0 {
		return
	}

	randIndex := rand.Intn(len(images))
	image := Decrypt(images[randIndex].Image, uh.Config.Masterkey)

	ResponseWithImage(w, image)
}
