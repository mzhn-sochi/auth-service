package usecase

import (
	"context"
	"github.com/mzhn-sochi/auth-service/internal/entity"
	"log/slog"
)

func (u *UseCase) FindById(ctx context.Context, id string) (*entity.User, error) {

	log := ctx.Value("logger").(*slog.Logger).With(slog.String("method", "FindById"))

	log.Debug("find user by id", slog.String("id", id))
	user, err := u.userStorage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	log.Debug("user found", slog.String("user", user.Phone))

	return user, nil
}
