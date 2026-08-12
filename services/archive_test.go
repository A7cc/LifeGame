package services

import (
	"LifeGame/core"
	"LifeGame/internal/db"
	"context"
	"strings"
	"testing"
)

func newArchiveTestApp(t *testing.T) *App {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	app := NewApp()
	app.Startup(context.Background())
	t.Cleanup(func() { app.Shutdown(context.Background()) })
	if !app.startupStatus.Ready {
		t.Fatalf("startup failed: %#v", app.startupStatus)
	}
	if code := responseCode(t, app.InitGame("存档测试", true, core.DifficultyNormal)); code != 200 {
		t.Fatalf("InitGame() code = %d", code)
	}
	return app
}

func TestLoadGameRejectsMismatchedSaveWithoutMutatingActiveState(t *testing.T) {
	app := newArchiveTestApp(t)
	saved := app.SaveGame("非当前格式")
	if saved.Code != 200 {
		t.Fatalf("SaveGame() = %#v", saved)
	}
	if _, err := db.GetDB().Exec("UPDATE saves SET save_version = ? WHERE id = ?", currentSaveVersion-1, saved.SaveID); err != nil {
		t.Fatal(err)
	}

	app.Userinfo.UCash = 12345
	loaded := app.LoadGame(int(saved.SaveID))
	if loaded.Code == 200 || !strings.Contains(loaded.Msg, "不一致") {
		t.Fatalf("mismatched save response = %#v", loaded)
	}
	if app.Userinfo.UCash != 12345 {
		t.Fatalf("failed load mutated active cash to %d", app.Userinfo.UCash)
	}
}

func TestLoadGameRequiresCompleteCurrentSave(t *testing.T) {
	app := newArchiveTestApp(t)
	saved := app.SaveGame("完整格式")
	if saved.Code != 200 {
		t.Fatalf("SaveGame() = %#v", saved)
	}
	if _, err := db.GetDB().Exec("DELETE FROM announce_saves WHERE save_id = ?", saved.SaveID); err != nil {
		t.Fatal(err)
	}

	loaded := app.LoadGame(int(saved.SaveID))
	if loaded.Code == 200 || !strings.Contains(loaded.Msg, "未找到完整存档") {
		t.Fatalf("incomplete save response = %#v", loaded)
	}
}
