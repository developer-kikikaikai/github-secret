package main

import (
	"fmt"
	"github-secret/api"
	"github-secret/arguments"
	"os"
)

const (
	//cmd name
	cmdUpdate  = "update"
	cmdDelete  = "delete"
	cmdGet     = "get"
	cmdList    = "list"
	cmdDefault = cmdUpdate
	//argument definition
	argRepo    = "repo"
	argOwner   = "owner"
	argToken   = "token"
	argSecName = "secname"
	argSecret  = "secret"
)

type argDefine struct {
	Usage    string
	Required bool
	Default  string
}

type argumentValues struct {
	vals    map[string]*string
	defines map[string]argDefine
}

func newArg() argumentValues {
	return argumentValues{
		vals: map[string]*string{},
		defines: map[string]argDefine{
			argRepo:    argDefine{Usage: "repository name to access", Required: true},
			argOwner:   argDefine{Usage: "repository owner", Required: true},
			argToken:   argDefine{Usage: "access token", Default: os.Getenv("GITHUB_TOKEN"), Required: false},
			argSecName: argDefine{Usage: "secret name", Required: true},
			argSecret:  argDefine{Usage: "secret value", Required: true},
		},
	}
}

func (arg argumentValues) setArgs(keys ...string) argumentValues {
	for _, key := range keys {
		arg.setArg(key)
	}
	return arg
}
func (arg argumentValues) setArg(key string) {
	define, ok := arg.defines[key]
	if !ok {
		return
	}
	if define.Required {
		arg.vals[key] = arguments.RequiredString(key, define.Usage)
	} else {
		arg.vals[key] = arguments.String(key, define.Default, define.Usage)
	}
}

func updateMain() {
	arg := newArg().setArgs(argRepo, argOwner, argToken, argSecName, argSecret)
	if !arguments.Parse() {
		arguments.Usage()
		return
	}

	client := api.NewSecretClient(*arg.vals[argOwner], *arg.vals[argRepo], *arg.vals[argToken])
	err := client.Update(*arg.vals[argSecName], *arg.vals[argSecret])
	if err != nil {
		fmt.Printf("Failed to call update, repos=%s\n", *arg.vals[argRepo])
	}
}

func deleteMain() {
	arg := newArg().setArgs(argRepo, argOwner, argToken, argSecName)
	if !arguments.Parse() {
		arguments.Usage()
		return
	}
	client := api.NewSecretClient(*arg.vals[argOwner], *arg.vals[argRepo], *arg.vals[argToken])
	err := client.Delete(*arg.vals[argSecName])
	if err != nil {
		fmt.Printf("Failed to call delete, repos=%s, secret=%s\n", *arg.vals[argRepo], *arg.vals[argSecName])
	}
}

func getMain() {
	arg := newArg().setArgs(argRepo, argOwner, argToken, argSecName)
	if !arguments.Parse() {
		arguments.Usage()
		return
	}
	client := api.NewSecretClient(*arg.vals[argOwner], *arg.vals[argRepo], *arg.vals[argToken])
	secret, err := client.Get(*arg.vals[argSecName])
	if err != nil {
		fmt.Printf("Failed to call get, repos=%s, secret=%s\n", *arg.vals[argRepo], *arg.vals[argSecName])
		return
	}
	fmt.Printf("Name:%v, CreatedAt:%v, UpdatedAt:%v\n", secret.Name, secret.CreatedAt, secret.UpdatedAt)
}

func listMain() {
	arg := newArg().setArgs(argRepo, argOwner, argToken)
	if !arguments.Parse() {
		arguments.Usage()
		return
	}
	client := api.NewSecretClient(*arg.vals[argOwner], *arg.vals[argRepo], *arg.vals[argToken])
	secrets, err := client.List()
	if err != nil {
		fmt.Printf("Failed to call get, repos=%s\n", *arg.vals[argRepo])
		return
	}
	for _, secret := range secrets.Secrets {
		fmt.Printf("    Name:%v, CreatedAt:%v, UpdatedAt:%v\n", secret.Name, secret.CreatedAt, secret.UpdatedAt)
	}
}

var cmds = map[string]func(){
	cmdUpdate: updateMain,
	cmdDelete: deleteMain,
	cmdGet:    getMain,
	cmdList:   listMain,
}

func getCmdNames() string {
	names := ""
	for name, _ := range cmds {
		names += name + " "
	}
	return names
}

func cmderr() {
	errmsg := fmt.Sprintf("command type is one of following commands:%s", getCmdNames())
	fmt.Printf("Usage of %s:\n", os.Args[0])
	fmt.Printf("\t<command type>, %s\n", errmsg)
}

func main() {
	if len(os.Args) <= 1 {
		cmderr()
		return
	}
	cmdType := os.Args[1]
	fnc, ok := cmds[cmdType]
	if !ok {
		cmderr()
		return
	}
	pname := os.Args[0] + " " + cmdType
	os.Args = os.Args[1:]
	os.Args[0] = pname
	fmt.Printf("Call %s\n", cmdType)
	fnc()
}
