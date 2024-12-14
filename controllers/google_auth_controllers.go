// controllers/google_auth_controllers.go
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server/config"
	"server/models"
	"server/utils"
	"server/validation"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type (
	GoogleService interface {
		GoogleAuthCommon(c *gin.Context, code string) *GoogleUserInfo
		GoogleSignIn(c *gin.Context)
		GoogleSignUp(c *gin.Context)
		GoogleSignInCallback(c *gin.Context)
		GoogleSignUpCallback(c *gin.Context)
	}

	GoogleUserInfo struct {
		ID            string        `json:"id"`
		Email         string        `json:"email"`
		VerifiedEmail bool          `json:"verified_email"`
		Name          string        `json:"name"`
		GivenName     string        `json:"given_name"`
		FamilyName    string        `json:"family_name"`
		Picture       string        `json:"picture"`
		Locale        string        `json:"locale"`
		Token         *oauth2.Token `json:"token"`
	}

	requestGoogleSignInData struct {
		Data []models.RequestSignInData `json:"data"`
	}

	requesGoogleSignUpData struct {
		Data []models.RequestSignUpData `json:"data"`
	}

	GoogleManager struct {
		GoogleService config.GoogleService
		UtilsFetcher  utils.UtilsFetcher
	}
)

func NewGoogleService(
	GoogleService config.GoogleService,
	utilsFetcher utils.UtilsFetcher,
) GoogleService {
	return &GoogleManager{
		GoogleService: GoogleService,
		UtilsFetcher:  utilsFetcher,
	}
}

func RedirectURL(env string, path string) string {
	var RedirectURI string
	if env != "local" {
		RedirectURI = fmt.Sprintf("%s/%s", os.Getenv("DOMAIN"), path)
	} else {
		RedirectURI = fmt.Sprintf("http://localhost:8080/%s", path)
	}
	return RedirectURI
}

func (gm *GoogleManager) GoogleSignIn(c *gin.Context) {
	url := gm.GoogleService.GoogleAuthURL(RedirectURL(os.Getenv("ENV"), "auth/google/signin/callback"))
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *GoogleManager) GoogleSignUp(c *gin.Context) {
	// リダイレクトURLを環境によって切り替える
	url := gm.GoogleService.GoogleAuthURL(RedirectURL(os.Getenv("ENV"), "auth/google/signup/callback"))
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *GoogleManager) GoogleAuthCommon(c *gin.Context, code string) *GoogleUserInfo {
	validator := validation.RequestGoogleCallbackData{
		Code: code,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return nil
	}

	// 認証コードからトークンを取得
	token, err := gm.GoogleService.GoogleOauth().Exchange(c, code)
	if err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return nil
	}

	// トークンを使ってユーザー情報を取得
	client := gm.GoogleService.GoogleOauth().Client(c, token)
	resp, err := client.Get(config.OauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return nil
	}
	defer resp.Body.Close()

	// ユーザー情報の取得
	var userInfo GoogleUserInfo
	userInfo.Token = token
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return nil
	}

	return &userInfo
}

func (gm *GoogleManager) GoogleSignInCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var err error
	var secure bool = false
	var domain string = "localhost"
	var httpOnly bool = false
	code := c.Query("code")

	validator := validation.RequestGoogleCallbackData{
		Code: code,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 認証コードからトークンを取得
	token, err := gm.GoogleService.GoogleOauth().Exchange(c, code)
	if err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// トークンを使ってユーザー情報を取得
	client := gm.GoogleService.GoogleOauth().Client(c, token)
	resp, err := client.Get(config.OauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer resp.Body.Close()

	// ユーザー情報の取得
	var userInfo GoogleUserInfo
	userInfo.Token = token
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	dbFetcher, _, _ := models.NewSignDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	registerData := requesGoogleSignUpData{
		Data: []models.RequestSignUpData{
			{
				UserName:     userInfo.Email,
				UserPassword: "google",
				NickName:     userInfo.Name,
			},
		},
	}
	if err := dbFetcher.PostSignUp(registerData.Data[0]); err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: "既に登録されたメールアドレスです。",
		}
		c.JSON(http.StatusConflict, response)
		return
	}
	getData := requestGoogleSignInData{
		Data: []models.RequestSignInData{
			{
				UserName:     userInfo.Email,
				UserPassword: "google",
			},
		},
	}
	dbFetcherSingIn, _, _ := models.NewSignDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := dbFetcherSingIn.GetSignIn(getData.Data[0])
	if err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	// UtilsFetcher を使用してトークンを生成
	newToken, err := gm.UtilsFetcher.NewToken(result[0].UserId, utils.AuthTokenHour)
	if err != nil {
		response := utils.ResponseWithSlice[requestSignInData]{
			ErrorMsg: "新規トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	refreshToken, err := gm.UtilsFetcher.RefreshToken(result[0].UserId, utils.RefreshAuthTokenHour)
	if err != nil {
		response := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "リフレッシュトークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// ローカルの場合
	if os.Getenv("ENV") != "local" {
		domain = os.Getenv("DOMAIN")
		secure = true
		httpOnly = true
	}

	c.SetCookie(utils.UserId, fmt.Sprintf("%d", result[0].UserId), 0, "/", domain, secure, httpOnly)
	c.SetCookie(utils.AuthToken, newToken, utils.AuthTokenHour*utils.SecondsInHour, "/", domain, secure, httpOnly)
	c.SetCookie(utils.RefreshAuthToken, refreshToken, utils.RefreshAuthTokenHour*utils.SecondsInHour, "/", domain, secure, httpOnly)
	c.SetCookie(utils.GoogleToken, userInfo.Token.AccessToken, -1, "/", domain, secure, httpOnly)

	fmt.Println("mail check: ", userInfo.Email)
	fmt.Println("mail check: ", userInfo.Token.AccessToken)

	c.JSON(http.StatusOK, userInfo)
}

func (gm *GoogleManager) GoogleSignUpCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var err error
	var secure bool = false
	var domain string = "localhost"
	var httpOnly bool = false
	code := c.Query("code")

	validator := validation.RequestGoogleCallbackData{
		Code: code,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 認証コードからトークンを取得
	token, err := gm.GoogleService.GoogleOauth().Exchange(c, code)
	if err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// トークンを使ってユーザー情報を取得
	client := gm.GoogleService.GoogleOauth().Client(c, token)
	resp, err := client.Get(config.OauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer resp.Body.Close()

	// ユーザー情報の取得
	var userInfo GoogleUserInfo
	userInfo.Token = token
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	dbFetcher, _, _ := models.NewSignDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	registerData := requesGoogleSignUpData{
		Data: []models.RequestSignUpData{
			{
				UserName:     userInfo.Email,
				UserPassword: "google",
				NickName:     userInfo.Name,
			},
		},
	}
	if err := dbFetcher.PostSignUp(registerData.Data[0]); err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: "既に登録されたメールアドレスです。",
		}
		c.JSON(http.StatusConflict, response)
		return
	}
	getData := requestGoogleSignInData{
		Data: []models.RequestSignInData{
			{
				UserName:     userInfo.Email,
				UserPassword: "google",
			},
		},
	}
	dbFetcherSingIn, _, _ := models.NewSignDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := dbFetcherSingIn.GetSignIn(getData.Data[0])
	if err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	// UtilsFetcher を使用してトークンを生成
	newToken, err := gm.UtilsFetcher.NewToken(result[0].UserId, utils.AuthTokenHour)
	if err != nil {
		response := utils.ResponseWithSlice[requestSignInData]{
			ErrorMsg: "新規トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	refreshToken, err := gm.UtilsFetcher.RefreshToken(result[0].UserId, utils.RefreshAuthTokenHour)
	if err != nil {
		response := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "リフレッシュトークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// ローカルの場合
	if os.Getenv("ENV") != "local" {
		domain = os.Getenv("DOMAIN")
		secure = true
		httpOnly = true
	}

	c.SetCookie(utils.UserId, fmt.Sprintf("%d", result[0].UserId), 0, "/", domain, secure, httpOnly)
	c.SetCookie(utils.AuthToken, newToken, utils.AuthTokenHour*utils.SecondsInHour, "/", domain, secure, httpOnly)
	c.SetCookie(utils.RefreshAuthToken, refreshToken, utils.RefreshAuthTokenHour*utils.SecondsInHour, "/", domain, secure, httpOnly)
	c.SetCookie(utils.GoogleToken, userInfo.Token.AccessToken, -1, "/", domain, secure, httpOnly)

	fmt.Println("mail check: ", userInfo.Email)
	fmt.Println("mail check: ", userInfo.Token.AccessToken)

	c.JSON(http.StatusOK, userInfo)
}
