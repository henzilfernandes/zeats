package gateway

import (
	"net/http"
	"github.com/pressly/chi"
	"io/ioutil"
	"github.com/golang/glog"
	"encoding/json"
	"zeats/types"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"errors"
)

const (
	ZeatsRequestAcceptance = iota + 100
	ZeatsRequestRejection
	ZeatsRequestFailure

	JWTSignKey = "secret"
)

// ctx dependency injection
type ctx struct {
	store Service
	h func(Service, http.ResponseWriter, *http.Request)
	a func(http.ResponseWriter, *http.Request)
}

func (g *ctx) handle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		bearerToken := strings.TrimSpace(request.Header.Get("authorization"))
		token, err := ValidateToken(bearerToken)

		if err != nil {
			glog.V(0).Infof("Token validation failed :: %s", err)
			WrapResponse(writer, http.StatusUnauthorized, ErrInvalidRequest(ZeatsRequestRejection,
				err.Error()))
			return
		}

		context.Set(request, "decoded", token.Claims)
		glog.V(0).Infof("Token validation success")
		g.h(g.store, writer, request)
	}
}

func (g *ctx) handleInit() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		g.a(writer, request)
	}
}

func Handler(store Service) http.Handler {
	r := chi.NewRouter()

	createToken := &ctx{a: CreateToken}
	insertEats := &ctx{store: store, h: InsertEats}
	updateEats := &ctx{store: store, h: UpdateEats}
	fetchEats := &ctx{store: store, h: FetchEats}
	deleteEats := &ctx{store: store, h: DeleteEats}

	r.Route("/" , func(r chi.Router) {
		r.Post("/createToken", createToken.handleInit())
		r.Post("/insert", insertEats.handle())
		r.Post("/update", updateEats.handle())
		r.Get("/fetch/:productId", fetchEats.handle())
		r.Delete("/delete/:productId", deleteEats.handle())
	})

	return r
}

func ValidateToken(bearerToken string) (*jwt.Token, error) {
	if bearerToken != "" {
		token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return false, errors.New("invalid authorization bearerToken")
			}
			return []byte(JWTSignKey), nil
		})

		if err != nil {
			return nil, errors.New("invalid authorization bearerToken")
		}

		if token.Valid {
			return token, nil
		} else {
			return nil, errors.New("invalid authorization bearerToken")
		}
	}

	return nil, errors.New("provide authorization header")
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})
	tokenString, err := token.SignedString([]byte(JWTSignKey))
	if err != nil {
		glog.V(0).Infof("Token creation failed :: %s", err)
		WrapResponse(w, http.StatusInternalServerError, ErrInvalidRequest(ZeatsRequestFailure,
			"Token creation failed"))
		return
	}

	WrapResponse(w, http.StatusOK, NewZeatPayloadResponse(ZeatsRequestAcceptance,
		"Authentication Token created successfully", &ZeatsResponse{Token: tokenString}))
}

func InsertEats(store Service, w http.ResponseWriter, r *http.Request) {
	// Extracting request
	var rawProductObj []byte
	rawProductObj, bErr := ioutil.ReadAll(r.Body)
	if bErr != nil {
		glog.V(0).Infof("Product object body extracion failed :: %s", bErr)
		WrapResponse(w, http.StatusBadRequest, ErrInvalidRequest(ZeatsRequestRejection,
			"Invalid request body"))
		return
	}

	prod := &types.Product{}
	if mErr := json.Unmarshal(rawProductObj, &prod); mErr != nil {
		glog.V(0).Infof("Product object unmarshalling error :: [%s]", mErr)
		WrapResponse(w, http.StatusInternalServerError, ErrInvalidRequest(ZeatsRequestFailure,
			"Unable to process payload"))
		return
	}

	err := store.InsertEats(prod)
	if err != nil {
		glog.V(0).Infof("Insert failed:: %s", err)
		WrapResponse(w, http.StatusInternalServerError, ErrInvalidRequest(ZeatsRequestFailure, err.Error()))
		return
	}

	WrapResponse(w, http.StatusOK, NewSuccessResponse(ZeatsRequestAcceptance,
		"Product data inserted successfully"))
}

func UpdateEats(store Service, w http.ResponseWriter, r *http.Request) {
	// Extracting request
	var rawProductObj []byte
	rawProductObj, bErr := ioutil.ReadAll(r.Body)
	if bErr != nil {
		glog.V(0).Infof("Product object body extracion failed :: %s", bErr)
		WrapResponse(w, http.StatusBadRequest, ErrInvalidRequest(ZeatsRequestRejection,
			"Invalid request body"))
		return
	}

	prod := &types.Product{}
	if mErr := json.Unmarshal(rawProductObj, &prod); mErr != nil {
		glog.V(0).Infof("Product object unmarshalling error :: [%s]", mErr)
		WrapResponse(w, http.StatusInternalServerError, ErrInvalidRequest(ZeatsRequestFailure,
			"Unable to process payload"))
		return
	}

	if prod.Id == "" {
		glog.V(0).Infof("Product Id not supplied")
		WrapResponse(w, http.StatusBadRequest, ErrInvalidRequest(ZeatsRequestRejection,
			"Product Id not supplied"))
		return
	}

	err := store.UpdateEats(prod)
	if err != nil {
		glog.V(0).Infof("Update failed:: %s", err)
		WrapResponse(w, http.StatusInternalServerError, ErrInvalidRequest(ZeatsRequestFailure, err.Error()))
		return
	}

	WrapResponse(w, http.StatusOK, NewSuccessResponse(ZeatsRequestAcceptance,
		"Product data updated successfully"))
}

func FetchEats(store Service, w http.ResponseWriter, r *http.Request) {
	productId := strings.TrimSpace(chi.URLParam(r, "productId"))
	if productId == "" {
		glog.V(0).Info("Invalid Product Id supplied")
		WrapResponse(w, http.StatusBadRequest, ErrInvalidRequest(ZeatsRequestRejection,
			"Invalid Product Id supplied"))
		return
	}

	prod, err := store.FetchEats(productId)
	if err != nil {
		glog.V(0).Info("Fetch error :: %s", err)
		WrapResponse(w, http.StatusInternalServerError, ErrInvalidRequest(ZeatsRequestFailure,
			"Couldn't getch product details"))
		return
	}


	WrapResponse(w, http.StatusOK, NewZeatPayloadResponse(ZeatsRequestAcceptance,
		"Product data fetched successfully", &ZeatsResponse{Payload: prod}))
}

func DeleteEats(store Service, w http.ResponseWriter, r *http.Request) {
	productId := strings.TrimSpace(chi.URLParam(r, "productId"))
	if productId == "" {
		glog.V(0).Info("Invalid Product Id supplied")
		WrapResponse(w, http.StatusBadRequest, ErrInvalidRequest(ZeatsRequestRejection,
			"Invalid Product Id supplied"))
		return
	}

	err := store.DeleteEats(productId)
	if err != nil {
		glog.V(0).Info("Delete error :: %s", err)
		WrapResponse(w, http.StatusInternalServerError, ErrInvalidRequest(ZeatsRequestFailure,
			"Couldn't delete product details"))
	}

	WrapResponse(w, http.StatusOK, NewSuccessResponse(ZeatsRequestAcceptance,
		"Product data deleted successfully"))
}

// Wrap response
func WrapResponse(w http.ResponseWriter, httpStatusCode int, commResp *ZeatsResponse) {

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	json.NewEncoder(w).Encode(commResp)
}

// Generic invalid request case
func ErrInvalidRequest(statusCode int, errMessage string) *ZeatsResponse {
	return &ZeatsResponse{StatusCode: statusCode, StatusMessage: errMessage}
}

// For write call success case
func NewSuccessResponse(statusCode int, message string) *ZeatsResponse {
	return &ZeatsResponse{StatusCode: statusCode, StatusMessage: message}
}

// For read payload call success case
func NewZeatPayloadResponse(statusCode int, message string, commResp *ZeatsResponse) *ZeatsResponse {
	commResp.StatusCode = statusCode
	commResp.StatusMessage = message
	return commResp
}

// Response is a wrapper response structure
type ZeatsResponse struct {
	StatusCode    int               `json:"statusCode,omitempty"`
	StatusMessage string            `json:"statusMessage,omitempty"`
	Token string                    `json:"token,omitempty"`
	RequestId     string            `json:"requestId,omitempty"`
	Payload       *types.Product    `json:"payload,omitempty"`
	Status        map[string]string `json:"status,omitempty"`
	CreatedOn     int64             `json:"createdOn,omitempty"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}
