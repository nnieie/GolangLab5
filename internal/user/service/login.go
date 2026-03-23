package service

import (
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/nnieie/golanglab5/internal/user/dal/db"
	"github.com/nnieie/golanglab5/internal/user/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/tracer"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func (s *UserService) Login(username, password string, code *string) (*base.User, error) {
	logger.Debugf("Login request(rpc): username=%s, password=%s, code=%v", username, password, code)

	user, err := db.QueryUserByNameWithPassword(s.ctx, username)
	if errors.Is(err, errno.UserIsNotExistErr) || user == nil || !utils.VerifyPassword(password, user.Password) {
		tracer.UserLoginCounter.Add(s.ctx, 1, metric.WithAttributes(
			attribute.String("status", "fail"),
		))
		return nil, errno.UserIsNotExistOrPasswordErr
	}
	if err != nil {
		tracer.UserLoginCounter.Add(s.ctx, 1, metric.WithAttributes(
			attribute.String("status", "fail"),
		))
		return nil, err
	}

	if user.TOTP != "" {
		if code == nil || !utils.CheckTotp(*code, user.TOTP) {
			logger.Infof("user %d login err: invalid mfa code, %v %s", user.ID, code, user.TOTP)
			tracer.UserLoginCounter.Add(s.ctx, 1, metric.WithAttributes(
				attribute.String("status", "fail"),
			))
			return nil, errno.MFAInvalidCodeErr
		}
	}

	tracer.UserLoginCounter.Add(s.ctx, 1, metric.WithAttributes(
		attribute.String("status", "success"),
	))
	return pack.DBUserTobaseUser(user), nil
}
