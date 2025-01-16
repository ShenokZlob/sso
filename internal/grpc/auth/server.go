package auth

import (
	"context"

	ssov1 "github.com/ShenokZlob/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	// Имя ssov1 мы использовали в контракте protos/proto/sso/sso.prot
	// option go_package = "tuzov.sso.v1;ssov1";
	//
	// Мы этой строчкой оснащаем нашу струтуру методами, которые описали в прото файле
	// Это временная заглушка (?)
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

// Ниже методы нашей структуры
// Двойные заглушки(?)

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}
