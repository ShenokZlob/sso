// severAPI обрабатывает все запросы сервера

package authgrpc

import (
	"context"
	"errors"
	"sso/internal/services/auth"

	ssov1 "github.com/ShenokZlob/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Повторюсь и тут. Узнать почему мы пишем интерфейс туть, а не в другом файле
// По поводу самого интерфейса. Он описывает основную бизнес-логику
type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int32,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	// Имя ssov1 мы использовали в контракте protos/proto/sso/sso.prot
	// option go_package = "tuzov.sso.v1;ssov1";
	//
	// Мы этой строчкой оснащаем нашу струтуру методами, которые описали в прото файле
	// Без нее мы не сможем зарегать наш обработчик
	// Мы делаем так, чтобы наша структура удолетворяла интерфейсу регистр сервера
	// На начальном этапе разработки выступает в роли заглушки
	ssov1.UnimplementedAuthServer
	auth Auth
}

const (
	emptyValue = 0
)

// Приложжение не запуститься, если просто вписать структуру

// Регистрирует наш обработчик
func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

// Ниже методы нашей структуры
// Двойные заглушки :) Наверное так надо

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	// Валидация данных
	// Обычно для этой задачи используют готовые библиотеки :)
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	// TODO: implement login via auth service
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		// Важно помнить, что мы пишем в ошибке
		// Например здесь мы не пишем подробно причину ошибки, чтобы клиенты не знали, че у нас под капотом

		// TODO: ...
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	// Validate
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO:...
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "alreasy exist error")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	// Validate
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO: ...
		if errors.Is(err, auth.ErrInvalidAppID) {
			return nil, status.Error(codes.PermissionDenied, "invalid app ID error")
		}

		return nil, status.Error(codes.InvalidArgument, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is not admin")
	}
	return nil
}
