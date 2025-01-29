package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/DroppedHard/SWIFT-service/service/api"
)

func CustomPathParameterValidationMiddleware(validateFn func(r *http.Request) error) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := validateFn(r); err != nil {
				api.WriteError(w, http.StatusBadRequest, err)
				return
			}
			next(w, r)
		}
	}
}

func BodyValidationMiddleware[T any](validateFn func(context.Context, *T) error) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var payload T

			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				api.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON payload: %v", err))
				log.Println(fmt.Errorf("invalid JSON payload: %v", err))
				return
			}

			if validateFn != nil {
				if err := validateFn(r.Context(), &payload); err != nil {
					api.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation error: %v", err))
					log.Println(fmt.Errorf("validation error: %v", err))
					return
				}
			}

			ctx := context.WithValue(r.Context(), reflect.TypeOf(payload), payload)
			next(w, r.WithContext(ctx))
		}
	}
}
