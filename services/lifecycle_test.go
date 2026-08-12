package services

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStartupFailureIsReportedAndBlocksGameInitialization(t *testing.T) {
	blockedHome := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(blockedHome, []byte("blocked"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", blockedHome)

	app := NewApp()
	t.Cleanup(func() { app.Shutdown(context.Background()) })
	app.Startup(context.Background())

	statusResponse := app.GetStartupStatus()
	status := statusResponse.Status
	if status.Ready || status.Error == "" || status.Stage != "加载配置" {
		t.Fatalf("unexpected startup status: %#v", status)
	}
	if !strings.Contains(status.Error, "加载配置失败") {
		t.Fatalf("startup error = %q, want stage context", status.Error)
	}
	if code := responseCode(t, app.InitGame("测试", true, 1)); code == 200 {
		t.Fatal("InitGame() succeeded after startup failure")
	}

	if err := os.Remove(blockedHome); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(blockedHome, 0o755); err != nil {
		t.Fatal(err)
	}
	retry := app.RetryStartup()
	if responseCode(t, retry) != 200 {
		t.Fatalf("RetryStartup() response = %#v", retry)
	}
	retriedStatus := retry.Status
	if !retriedStatus.Ready || retriedStatus.Error != "" {
		t.Fatalf("retry status = %#v, want ready", retriedStatus)
	}
}
