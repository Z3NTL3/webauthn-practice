package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"charm.land/log/v2"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/z3ntl3/webauthn-practice/globals"
	"github.com/z3ntl3/webauthn-practice/models"
	"github.com/z3ntl3/webauthn-practice/utils"
)

func LoginBegin(res http.ResponseWriter, req *http.Request) {
	webauth := globals.WEBAUTH
	STORAGE := globals.STORAGE

	var email string
	body, err := io.ReadAll(req.Body)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed reading response body")
		return
	}

	err = json.Unmarshal(body, &email)
	if err != nil {
		utils.Fail(res, http.StatusBadRequest, "bad request")
		return
	}

	if !(len(email) > 0) {
		utils.Fail(res, http.StatusBadRequest, "cannot provide empty email")
		return
	}

	user, ok := STORAGE.Load(email)
	if !ok {
		utils.Fail(res, http.StatusBadRequest, "user has not registered with us")
		return
	}

	user_, ok := user.(map[any]any)
	if !ok {
		utils.Fail(res, http.StatusInternalServerError, "couldnt find user")
		return
	}

	rawUserData, ok := user_["user"].([]byte)
	if !ok {
		utils.Fail(res, http.StatusInternalServerError, "tried parsing invalid raw data")
		return
	}

	userModel := new(models.User)
	if err = json.Unmarshal(rawUserData, userModel); err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed unmarshalling raw data to user")
		return
	}

	log.Debug("login-begin-handle", "email", email)
	opts, sessData, err := webauth.BeginLogin(userModel, webauthn.WithAllowedCredentials([]protocol.CredentialDescriptor{
		{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: userModel.Credentials[0].ID,
		},
	}))

	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, err.Error())
		return
	}

	rawUser, err := json.Marshal(userModel)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, err.Error())
		return
	}

	rawSessionData, err := json.Marshal(sessData)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, err.Error())
		return
	}

	STORAGE.Store(email, map[any]any{
		"user":        rawUser,
		"sessionData": rawSessionData,
	})

	log.Debug("session-login-begin", "data", *sessData)
	utils.OK(res, http.StatusOK, opts)
}
