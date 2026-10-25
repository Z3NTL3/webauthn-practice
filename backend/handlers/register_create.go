package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/z3ntl3/webauthn-practice/globals"
	"github.com/z3ntl3/webauthn-practice/models"
	"github.com/z3ntl3/webauthn-practice/utils"
)

func RegisterCreateOptions(res http.ResponseWriter, req *http.Request) {
	webauth := globals.WEBAUTH
	STORAGE := globals.STORAGE

	body, err := io.ReadAll(req.Body)
	if err != nil {
		utils.Fail(res, http.StatusBadRequest, "bad request body")
		return
	}

	data := map[string]string{}
	if err := json.Unmarshal(body, &data); err != nil {
		utils.Fail(res, http.StatusBadRequest, "bad request body")
		return
	}

	email := data["email"]
	if !(len(email) > 0) {
		utils.Fail(res, http.StatusBadRequest, "invalid email provided")
		return
	}

	user := models.NewUser(email)
	creation, sessionData, err := webauth.BeginMediatedRegistration(user, protocol.MediationRequired,
		webauthn.WithCredentialParameters([]protocol.CredentialParameter{
			{
				Type:      protocol.PublicKeyCredentialType,
				Algorithm: webauthncose.AlgRS256,
			},
			{
				Type:      protocol.PublicKeyCredentialType,
				Algorithm: webauthncose.AlgES256,
			},
		}),
	)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed registiration handshake")
		return
	}

	sessionDataRaw, err := json.Marshal(sessionData)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed marshalling session data")
		return
	}

	userRaw, err := json.Marshal(user)
	if err != nil {
		utils.Fail(res, http.StatusInternalServerError, "failed marshalling user data")
		return
	}

	STORAGE.Store(email, map[any]any{
		"user":        userRaw,
		"sessionData": sessionDataRaw,
	})

	utils.OK(res, http.StatusOK, creation)
}
