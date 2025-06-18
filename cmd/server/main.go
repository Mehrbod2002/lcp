package main

import (
	"context"

	"github.com/Mehrbod2002/lcp.git/internal/adapter/repository/lcp"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/readium/readium-lcp-server/lcpencrypt"
	"github.com/readium/readium-lcp-server/lcpserver"
	"github.com/swaggo/fiber-swagger"
	_ "github.com/swaggo/swag/example/celler/httputil"
	"github.com/yourusername/lcp-project/internal/adapter/graphql"
	"github.com/yourusername/lcp-project/internal/adapter/jwt"
	"github.com/yourusername/lcp-project/internal/config"
	"github.com/yourusername/lcp-project/internal/usecase/lcp/license"
	"github.com/yourusername/lcp-project/internal/usecase/lcp/publication"
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

	db, err := pgx.Connect(context.Background(), cfg.Database.DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close(context.Background())

	lcpEnc := lcpencrypt.NewEncrypter(cfg.LCP.Certificate, cfg.LCP.PrivateKey, cfg.LCP.Storage)
	lcpSrv := lcpserver.NewServer(cfg.LCP)
	pubRepo := lcp.NewPublicationRepository(db)
	licRepo := lcp.NewLicenseRepository(db)
	pubUsecase := publication.NewPublicationUsecase(pubRepo, lcpEnc)
	licUsecase := license.NewLicenseUsecase(licRepo, lcpSrv)

	app := fiber.New()
	app.Use(jwt.Middleware(cfg.JWT.Secret))

	// GraphQL endpoint
	gqlHandler := graphql.NewHandler(&graphql.Resolver{
		PublicationUsecase: pubUsecase,
		LicenseUsecase:     licUsecase,
	})
	app.Post("/graphql", adaptor.HTTPHandler(gqlHandler))

	// Swagger endpoint
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Listen(cfg.Server.Port)
}
