package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/context"

	"github.com/tespo/buddha/util"
)

//
// ExplicitRouterAuthenticationWrapper wraps all developer
// handlers validating tokens and roles/permissions on endpoints
//
func ExplicitRouterAuthenticationWrapper(route, method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", os.Getenv("VIJNANA_URL")+"/validate-token-permissions", nil)
		req.Header = r.Header
		req.Header.Set("Route", route)
		req.Header.Set("Method", method)
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
		context.Set(r, "token", tokenString)
		next(w, r)
	}
}

//
// LambdaRouterAuthenticationWrapper wraps all developer
// handlers validating tokens and roles/permissions on endpoints
//
func LambdaRouterAuthenticationWrapper(route, method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", os.Getenv("VIJNANA_URL")+"/validate-token-permissions", nil)
		req.Header = r.Header
		req.Header.Set("Route", route)
		req.Header.Set("Method", method)
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
		next(w, r)
	}
}

//
// ImplicitRouterAuthenticationWrapper validates tokens and
// adds context to the handlers for the user data
//
func ImplicitRouterAuthenticationWrapper(requiredScope string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", os.Getenv("VIJNANA_URL")+"/validate-token", nil)
		req.Header = r.Header
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
		scopePermissions, ok := claims["scope_permissions"].([]interface{})
		if !ok {
			util.ErrorResponder(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		scopedFields, ok := claims["scoped_fields"].([]interface{})
		if !ok {
			util.ErrorResponder(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		for count, scope := range scopePermissions {
			if checkPattern(scope.(string), requiredScope) {
				break
			}
			if count == len(scopePermissions)-1 {
				util.ErrorResponder(w, http.StatusUnauthorized, errors.New("Unauthorized"))
				return
			}
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			util.ErrorResponder(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		accountID, ok := claims["account_id"].(string)
		if !ok {
			util.ErrorResponder(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		owner, ok := claims["owner"].(bool)
		if !ok {
			util.ErrorResponder(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		scopedFieldSlice := make([]string, len(scopedFields))
		for i, v := range scopedFields {
			scopedFieldSlice[i] = v.(string)
		}
		context.Set(r, "user_id", userID)
		context.Set(r, "account_id", accountID)
		context.Set(r, "owner", owner)
		context.Set(r, "token", tokenString)
		context.Set(r, "scoped_fields", scopedFieldSlice)
		next(w, r)
	}
}

func checkPattern(pattern, match string) bool {
	var validator *regexp.Regexp
	if pattern == "*" {
		pattern = "." + pattern
	}
	validator = regexp.MustCompile("^" + pattern)
	return validator.MatchString(match)
}
