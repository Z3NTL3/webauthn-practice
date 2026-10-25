package handlers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/z3ntl3/webauthn-practice/globals"
	"github.com/z3ntl3/webauthn-practice/models"
	"github.com/z3ntl3/webauthn-practice/utils"
)

func LoginFinalize(res http.ResponseWriter, req *http.Request) {
	webauth := globals.WEBAUTH
	STORAGE := globals.STORAGE

	type ReqBodySchema struct {
		Email     string          `json:"email"`
		Assertion json.RawMessage `json:"assertion"`
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		utils.Fail(res, http.StatusBadRequest, "bad request body")
		return
	}

	var schema ReqBodySchema
	if err := json.Unmarshal(body, &schema); err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed unmarshalling body")
		return
	}

	email := schema.Email
	if !(len(email) > 0) {
		utils.Fail(res, http.StatusBadRequest, "invalid email")
		return
	}

	sess_, ok := STORAGE.Load(email)
	if !ok {
		utils.Fail(res, http.StatusBadRequest, "user does not exist")
		return
	}

	sess, ok := sess_.(map[any]any)
	if !ok {
		utils.Fail(res, http.StatusInternalServerError, "failed loading data from existing user")
		return
	}

	userRawData, ok := sess["user"].([]byte)
	sessRawData, ok := sess["sessionData"].([]byte)

	if !ok {
		utils.Fail(res, http.StatusInternalServerError, "failed loading data from user/session data")
		return
	}

	userModel := new(models.User)
	if err := json.Unmarshal(userRawData, userModel); err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed unmarshalling raw user data into model")
		return
	}

	sessData := new(webauthn.SessionData)
	if err := json.Unmarshal(sessRawData, sessData); err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed unmarshalling raw session data into model")
		return
	}

	parsed, err := protocol.ParseCredentialRequestResponseBytes(schema.Assertion)
	if err := json.Unmarshal(sessRawData, sessData); err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed parsing credential response from client")
		return
	}

	credential, err := webauth.ValidateLogin(userModel, *sessData, parsed)
	if err := json.Unmarshal(sessRawData, sessData); err != nil {
		utils.Fail(res, http.StatusBadRequest, err.Error())
		return
	}

	userModel.Credentials[0] = *credential

	rawUserData, err := json.Marshal(userModel)
	if err := json.Unmarshal(sessRawData, sessData); err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed marshalling user data")
		return
	}

	STORAGE.Store(email, map[any]any{"user": rawUserData})
	utils.OK(res, http.StatusOK, map[string]string{
		"id":     base64.StdEncoding.EncodeToString(credential.ID),
		"pubKey": base64.StdEncoding.EncodeToString(credential.PublicKey),
	})
}
