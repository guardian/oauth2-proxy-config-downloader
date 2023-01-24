package main

import "fmt"

type AppId struct {
	App   string
	Stack string
	Stage string
}

func NewAppId(app string, stack string, stage string) *AppId {
	return &AppId{
		App:   app,
		Stack: stack,
		Stage: stage,
	}
}

func (a *AppId) SsmPath(key string) string {
	return fmt.Sprintf("/%s/%s/%s/%s", a.Stage, a.Stack, a.App, key)
}

func (a *AppId) CookieName() string {
	return fmt.Sprintf("_%s_%s_%s", a.App, a.Stack, a.Stage)
}
