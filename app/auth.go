package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
)

const NoPassFlag = "nopass"

type User struct {
	name      string
	flags     map[string]bool
	passwords []string
}

var users map[string]User

func init() {
	users = make(map[string]User)
	users["default"] = User{name: "default", flags: map[string]bool{NoPassFlag: true}, passwords: make([]string, 0)}
}

func aclHandler(ctx context.Context, args []string) any {
	switch strings.ToUpper(args[0]) {
	case "WHOAMI":
		return BulkStr("default")
	case "GETUSER":
		return getUserHandler(args[1:])
	case "SETUSER":
		return setUserHandler(args[1:])
	}
	return errors.New("ERR not implemented")
}

func getUserHandler(args []string) any {
	username := args[0]
	user, found := users[username]
	if !found {
		return errors.New("ERR user not found")
	}
	flags := make([]BulkStr, 0)
	for flag, enabled := range user.flags {
		if enabled {
			flags = append(flags, BulkStr(flag))
		}
	}

	fmt.Println(users)
	passwords := make([]BulkStr, len(user.passwords))
	for i, passwd := range user.passwords {
		passwords[i] = BulkStr(passwd)
	}

	return []any{
		BulkStr("flags"), flags,
		BulkStr("passwords"), passwords,
	}
}

func setUserHandler(args []string) any {
	username := args[0]
	user, found := users[username]
	if !found {
		return errors.New("ERR user not found")
	}
	passwd, _ := strings.CutPrefix(args[1], ">")
	user.passwords = append(user.passwords, hash(passwd))
	user.flags[NoPassFlag] = false
	users[username] = user
	return RespStr("OK")
}

func hash(str string) string {
	hashed := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hashed[:])
}

func authHandler(ctx context.Context, args []string) any {
	username := args[0]
	password := hash(args[1])
	user, found := users[username]
	if !found {
		return errors.New("WRONGPASS invalid username-password pair or user is disabled.")
	}
	if !slices.Contains(user.passwords, password) {
		return errors.New("WRONGPASS invalid username-password pair or user is disabled.")
	}
	ctxUsername := ctx.Value(UsernameContextKey).(*string)
	*ctxUsername = user.name
	return RespStr("OK")
}

func authMiddleware(f handler) handler {
	return func(ctx context.Context, s []string) any {
		ctxUsername := ctx.Value(UsernameContextKey).(*string)
		if *ctxUsername != "" {
			return f(ctx, s)
		}
		username := "default"
		user, found := users[username]
		if !found {
			return f(ctx, s)
		}
		if user.flags[NoPassFlag] {
			*ctxUsername = user.name
			return f(ctx, s)
		}
		return errors.New("NOAUTH Authentication required.")
	}
}
