package auth

import "context"

type contextKey string

const emailKey contextKey = "auth_email"

func SetEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, emailKey, email)
}

func GetEmail(ctx context.Context) string {
	email, _ := ctx.Value(emailKey).(string)
	return email
}
