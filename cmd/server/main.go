package main

import (
	"net/http"

	"github.com/Mehrbod2002/lcp/internal/adapter/graphql"
	"github.com/Mehrbod2002/lcp/internal/adapter/repository/lcp"
	"github.com/Mehrbod2002/lcp/internal/config"
	lcpencrypt "github.com/Mehrbod2002/lcp/internal/lcp/encrypt"
	lcplicense "github.com/Mehrbod2002/lcp/internal/lcp/license"
	"github.com/Mehrbod2002/lcp/internal/usecase/lcp/license"
	"github.com/Mehrbod2002/lcp/internal/usecase/lcp/publication"
)

// @title LCP License Server API
// @version 1.0
// @description API for managing LCP licenses and publications
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	lcpEnc := lcpencrypt.NewFileCopyEncrypter(cfg.LCP.Storage.FS.Directory)
	lcpSrv := lcplicense.NewService()
	pubRepo := lcp.NewPublicationRepository()
	licRepo := lcp.NewLicenseRepository()
	pubUsecase := publication.NewPublicationUsecase(pubRepo, lcpEnc)
	licUsecase := license.NewLicenseUsecase(licRepo, lcpSrv)

	mux := http.NewServeMux()

	gqlHandler := graphql.NewHandler(&graphql.Resolver{
		PublicationUsecase: pubUsecase,
		LicenseUsecase:     licUsecase,
	})
	mux.Handle("/graphql", gqlHandler)

	port := cfg.Server.Port
	if port == "" {
		port = ":8080"
	}

	if err := http.ListenAndServe(port, mux); err != nil {
		panic(err)
	}
}
