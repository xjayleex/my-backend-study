package cli

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	clitool "gopkg.in/urfave/cli.v2"
)

type User struct {
	Username		string
	HashedPassword	string
}


var SignUp = clitool.Command{
	Name: "sign-up",
	Usage: "Sign up for a gRPC Server.",
	Action: signupAction,
	Flags: []clitool.Flag{
		&clitool.StringFlag{
			Name: "username",
			Usage: "username",
		},
		&clitool.StringFlag{
			Name: "password",
			Usage: "password",
		},
	},
}

func signupAction(c *clitool.Context) (err error) {
	var (
		userName		= c.String("username")
		password		= c.String("password")
	)

	if userName == "" {
		trap(errors.New("User Name Required."))
	}

	if len(password) <= 4 {
		trap(errors.New("Required over 4 characters for your password "))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)

	trap(err)
	// gRPC Client -> Info data -> call rpc SignUp
}