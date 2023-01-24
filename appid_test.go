package main

import "testing"

func TestSsmPath(t *testing.T) {
	appId := NewAppId("test-app", "test-stack", "PROD")
	p := appId.SsmPath("some/ssm/path")
	if p != "/PROD/test-stack/test-app/some/ssm/path" {
		t.Error("Unexpected SSM path, got ", p)
	}
}
