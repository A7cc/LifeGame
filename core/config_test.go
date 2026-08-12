package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitConfigCreatesDefaultAndLoadsExistingOverride(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	previousConfig := AppConfig
	t.Cleanup(func() { AppConfig = previousConfig })

	if err := InitConfig(); err != nil {
		t.Fatalf("first InitConfig() error = %v", err)
	}
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("default config was not created: %v", err)
	}

	custom := *DefaultConfig
	custom.Game = DefaultConfig.Game
	custom.Game.Name = "自定义人生"
	if err := SaveConfigToYAML(configPath, &custom); err != nil {
		t.Fatal(err)
	}
	if err := InitConfig(); err != nil {
		t.Fatalf("second InitConfig() error = %v", err)
	}
	if got := GetGameName(); got != "自定义人生" {
		t.Fatalf("game name = %q, want existing override", got)
	}
}

func TestExternalAssetDirectoriesUseCurrentDataDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	imagesDir, err := GetImagesDir()
	if err != nil {
		t.Fatal(err)
	}
	audioDir, err := GetAudioDir()
	if err != nil {
		t.Fatal(err)
	}
	if imagesDir != filepath.Join(home, ".lifegame", "images") {
		t.Fatalf("images directory = %q", imagesDir)
	}
	if audioDir != filepath.Join(home, ".lifegame", "audio") {
		t.Fatalf("audio directory = %q", audioDir)
	}
	for _, directory := range []string{imagesDir, audioDir} {
		if info, err := os.Stat(directory); err != nil || !info.IsDir() {
			t.Fatalf("external asset directory %q was not created: %v", directory, err)
		}
	}
}
