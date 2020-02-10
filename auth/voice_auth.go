package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gorilla/context"
	"github.com/tespo/buddha/util"
	"github.com/tespo/satya/v2/types"
)

//
// AuthenticateVoiceRequest handles parsing and validating
// a jwt token from different voice command services
//
func AuthenticateVoiceRequest(provider string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", os.Getenv("VIJNANA_URL")+"/validate-token", nil)
		req.Header = r.Header
		if provider == "alexa" {
			req.Header.Set("Authorization", "Bearer "+getAlexaToken(r))
		}
		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if err = json.NewEncoder(w).Encode(map[string]string{"Error": err.Error()}); err != nil {
				panic(err)
			}
			return
		}
		if resp.StatusCode != 200 {
			w.WriteHeader(http.StatusUnauthorized)
			if err = json.NewEncoder(w).Encode(map[string]string{"Error": "Unauthorized"}); err != nil {
				panic(err)
			}
			return
		}
		tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
		claims, err := util.ParseToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if err = json.NewEncoder(w).Encode(map[string]string{"Error": err.Error()}); err != nil {
				panic(err)
			}
			return
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			util.ErrorResponder(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		context.Set(r, "user_id", userID)
		next(w, r)
	}
}

func getAlexaToken(r *http.Request) string {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	var d map[string]interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return ""
	}
	alexaRequest := types.AlexaRequest{}
	if err := json.Unmarshal(data, &alexaRequest); err != nil {
		return ""
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	if !reflect.DeepEqual(alexaRequest.Directive.Endpoint, types.AlexaDiscoverResponseEventPayloadEndpoint{}) {
		return alexaRequest.Directive.Endpoint.Scope.Token
	}
	return alexaRequest.Directive.Payload.Scope.Token
}
