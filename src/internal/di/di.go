package di

import (
	"WebTic-tac-toe2/internal/app"
	"WebTic-tac-toe2/internal/config"
	"WebTic-tac-toe2/internal/infra"
	"WebTic-tac-toe2/internal/service/auth"
	"WebTic-tac-toe2/internal/service/jwtauth"
	web "WebTic-tac-toe2/internal/transport/http"
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/fx"
)

func invokeServer(lf fx.Lifecycle, handler http.Handler, conf *config.Config) {
	server := &http.Server{
		Addr:    conf.ServerAddr,
		Handler: handler,
	}

	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			fmt.Printf("Server started on %q\n", server.Addr)
			go server.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Printf("Server stopped!")
			return server.Shutdown(ctx)
		},
	})
}

var Module = fx.Options(
	fx.Provide(config.ParseConfig),
	fx.Provide(infra.NewPool),

	fx.Provide(fx.Annotate(
		infra.NewRepo,
		fx.As(new(app.GameRepository)),
		fx.As(new(auth.AuthRepo)),
		fx.As(new(jwtauth.JwtRepoService)),
	)),

	fx.Provide(fx.Annotate(
		auth.NewAuthService,
		fx.As(new(auth.UserService)),
	)),

	fx.Provide(fx.Annotate(
		jwtauth.NewJwtProvider,
		fx.As(new(jwtauth.JWTService)),
	)),

	fx.Provide(app.NewGameService),
	fx.Provide(web.NewMainHandler),
	fx.Provide(http.NewServeMux),
	fx.Provide((*web.MainHandler).RegisterRoutes),

	fx.Invoke(
		infra.MakeMigrations,
		invokeServer,
	),
)
