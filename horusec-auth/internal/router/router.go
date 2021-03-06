// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"github.com/ZupIT/horusec/development-kit/pkg/databases/relational"
	cacheRepository "github.com/ZupIT/horusec/development-kit/pkg/databases/relational/repository/cache"
	brokerLib "github.com/ZupIT/horusec/development-kit/pkg/services/broker"
	serverConfig "github.com/ZupIT/horusec/development-kit/pkg/utils/http/server"
	"github.com/ZupIT/horusec/horusec-auth/config/app"
	"github.com/ZupIT/horusec/horusec-auth/internal/handler/account"
	"github.com/ZupIT/horusec/horusec-auth/internal/handler/auth"
	"github.com/ZupIT/horusec/horusec-auth/internal/handler/health"
	"github.com/ZupIT/horusec/horusec-auth/internal/router/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	config *serverConfig.Server
	router *chi.Mux
}

func NewRouter(config *serverConfig.Server) *Router {
	return &Router{
		config: config,
		router: chi.NewRouter(),
	}
}

func (r *Router) setMiddleware() {
	r.EnableRealIP()
	r.EnableLogger()
	r.EnableRecover()
	r.EnableTimeout()
	r.EnableCompress()
	r.EnableRequestID()
	r.EnableCORS()
	r.RouterMetrics()
}

func (r *Router) GetRouter(postgresRead relational.InterfaceRead, postgresWrite relational.InterfaceWrite,
	broker brokerLib.IBroker, cache cacheRepository.Interface, appConfig *app.Config) *chi.Mux {
	r.setMiddleware()
	r.RouterAuth(postgresRead, postgresWrite, appConfig)
	r.RouterHealth(postgresRead, postgresWrite)
	r.RouterAccount(postgresRead, postgresWrite, broker, cache, appConfig)
	return r.router
}

func (r *Router) EnableRealIP() *Router {
	r.router.Use(middleware.RealIP)
	return r
}

func (r *Router) EnableLogger() *Router {
	r.router.Use(middleware.Logger)
	return r
}

func (r *Router) EnableRecover() *Router {
	r.router.Use(middleware.Recoverer)
	return r
}

func (r *Router) EnableTimeout() *Router {
	r.router.Use(middleware.Timeout(r.config.GetTimeout()))
	return r
}

func (r *Router) EnableCompress() *Router {
	r.router.Use(middleware.Compress(r.config.GetCompression()))
	return r
}

func (r *Router) EnableRequestID() *Router {
	r.router.Use(middleware.RequestID)
	return r
}

func (r *Router) EnableCORS() *Router {
	r.router.Use(r.config.Cors)
	return r
}

func (r *Router) RouterMetrics() *Router {
	r.router.Handle("/metrics", promhttp.Handler())
	return r
}

func (r *Router) RouterHealth(postgresRead relational.InterfaceRead, postgresWrite relational.InterfaceWrite) *Router {
	handler := health.NewHandler(postgresRead, postgresWrite)
	r.router.Route(routes.HealthHandler, func(router chi.Router) {
		router.Get("/", handler.Get)
		router.Options("/", handler.Options)
	})

	return r
}

func (r *Router) RouterAuth(
	postgresRead relational.InterfaceRead, postgresWrite relational.InterfaceWrite, appConfig *app.Config) *Router {
	handler := auth.NewAuthHandler(postgresRead, postgresWrite, appConfig)
	r.router.Route(routes.AuthHandler, func(router chi.Router) {
		router.Get("/config", handler.Config)
		router.Post("/authenticate", handler.AuthByType)
		router.Options("/", handler.Options)
	})

	return r
}

// nolint
func (r *Router) RouterAccount(postgresRead relational.InterfaceRead, postgresWrite relational.InterfaceWrite,
	broker brokerLib.IBroker, cache cacheRepository.Interface, appConfig *app.Config) *Router {
	handler := account.NewHandler(broker, postgresRead, postgresWrite, cache, appConfig)
	r.router.Route(routes.AccountHandler, func(router chi.Router) {
		router.Post("/create-account-from-keycloak", handler.CreateAccountFromKeycloak)
		router.Post("/create-account", handler.CreateAccount)
		router.Get("/validate/{accountID}", handler.ValidateEmail)
		router.Post("/send-code", handler.SendResetPasswordCode)
		router.Post("/validate-code", handler.ValidateResetPasswordCode)
		router.Post("/change-password", handler.ChangePassword)
		router.Post("/renew-token", handler.RenewToken)
		router.Post("/logout", handler.Logout)
		router.Delete("/delete", handler.DeleteAccount)
		router.Post("/verify-already-used", handler.VerifyAlreadyInUse)
		router.Patch("/update", handler.Update)
		router.Options("/", handler.Options)
	})

	return r
}
