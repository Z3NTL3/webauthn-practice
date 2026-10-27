package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/z3ntl3/webauthn-practice/globals"
	"github.com/z3ntl3/webauthn-practice/models"
	"github.com/z3ntl3/webauthn-practice/utils"
)

func RegisterFinalize(res http.ResponseWriter, req *http.Request) {
	webauth := globals.WEBAUTH
	STORAGE := globals.STORAGE

	type ReqBodySchema struct {
		Email      string          `json:"email"`
		Credential json.RawMessage `json:"credential"`
	}

	var schema ReqBodySchema
	if err := json.NewDecoder(req.Body).Decode(&schema); err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed unmarshalling body")
		return
	}

	email := schema.Email
	if !(len(email) > 0) {
		utils.Fail(res, http.StatusBadRequest, "email cannot be empty")
		return
	}

	sess, ok := STORAGE.Load(email)
	if !ok {
		utils.Fail(res, http.StatusBadRequest, "must call create options endpoint first")
		return
	}

	sess_, ok := sess.(map[any]any)
	if !ok {
		utils.Fail(res, http.StatusInternalServerError, "failed retrieving previous step session")
		return
	}

	user_, ok := sess_["user"].([]byte)
	if !ok {
		utils.Fail(res, http.StatusInternalServerError, "failed retrieving user")
		return
	}

	sessData_, ok := sess_["sessionData"].([]byte)
	if !ok {
		utils.Fail(res, http.StatusInternalServerError, "failed retrieving webauthn session data")
		return
	}

	userModel := new(models.User)
	err := json.Unmarshal(user_, userModel)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed unmarshalling user to model")
		return
	}

	sessDataModel := new(webauthn.SessionData)
	err = json.Unmarshal(sessData_, sessDataModel)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed unmarshalling session data to model")
		return
	}

	parsedCredential, err := protocol.ParseCredentialCreationResponseBytes(schema.Credential)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed parsing credential creation response")
		return
	}

	credential, err := webauth.CreateCredential(userModel, *sessDataModel, parsedCredential)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed webauthn registiration finalization")
		return
	}

	userModel.Credentials = []webauthn.Credential{*credential}

	raw, err := json.Marshal(userModel)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed marshalling user")
		return
	}

	STORAGE.Store(email, map[any]any{
		"user": raw,
	})

	utils.OK(res, http.StatusOK, "registiration complete")
}
