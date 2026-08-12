package main

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestWailsAssetServerUsesExternalOverrideBeforePackagedAsset(t *testing.T) {
	packaged := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("frontend")},
		"images/datinginfo/dating-partner/female/11.webp": &fstest.MapFile{Data: []byte("packaged-portrait")},
	}
	imagesDir := t.TempDir()
	overridePath := filepath.Join(imagesDir, "datinginfo", "dating-partner", "female", "11.webp")
	if err := os.MkdirAll(filepath.Dir(overridePath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(overridePath, []byte("user-override"), 0644); err != nil {
		t.Fatal(err)
	}

	options := newAssetServerOptions(fs.FS(packaged), imagesDir, "")
	// Wails 会把 Middleware 包在默认的内置资源处理器外层。这里使用同样的
	// 公开契约验证：即使内置 FS 有同名文件，外部图片也必须先返回。
	handler := options.Middleware(http.FileServer(http.FS(packaged)))

	overrideResponse := httptest.NewRecorder()
	handler.ServeHTTP(overrideResponse, httptest.NewRequest("GET", "/images/datinginfo/dating-partner/female/11.webp", nil))
	if overrideResponse.Code != 200 || strings.TrimSpace(overrideResponse.Body.String()) != "user-override" {
		t.Fatalf("override response = %d %q", overrideResponse.Code, overrideResponse.Body.String())
	}
	if overrideResponse.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("override Cache-Control = %q, want no-store", overrideResponse.Header().Get("Cache-Control"))
	}

	if err := os.Remove(overridePath); err != nil {
		t.Fatal(err)
	}
	packagedResponse := httptest.NewRecorder()
	handler.ServeHTTP(packagedResponse, httptest.NewRequest("GET", "/images/datinginfo/dating-partner/female/11.webp", nil))
	if packagedResponse.Code != 200 || strings.TrimSpace(packagedResponse.Body.String()) != "packaged-portrait" {
		t.Fatalf("packaged response = %d %q", packagedResponse.Code, packagedResponse.Body.String())
	}

}

func TestWailsAssetServerUsesExternalAudioBeforePackagedAsset(t *testing.T) {
	packaged := fstest.MapFS{
		"index.html":               &fstest.MapFile{Data: []byte("frontend")},
		"audio/lifegame-theme.wav": &fstest.MapFile{Data: []byte("packaged-music")},
	}
	audioDir := t.TempDir()
	overridePath := filepath.Join(audioDir, "lifegame-theme.wav")
	if err := os.WriteFile(overridePath, []byte("user-music"), 0o644); err != nil {
		t.Fatal(err)
	}

	options := newAssetServerOptions(fs.FS(packaged), "", audioDir)
	handler := options.Middleware(http.FileServer(http.FS(packaged)))
	overrideResponse := httptest.NewRecorder()
	handler.ServeHTTP(overrideResponse, httptest.NewRequest("GET", "/audio/lifegame-theme.wav", nil))
	if overrideResponse.Code != 200 || strings.TrimSpace(overrideResponse.Body.String()) != "user-music" {
		t.Fatalf("override response = %d %q", overrideResponse.Code, overrideResponse.Body.String())
	}
	if overrideResponse.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("override Cache-Control = %q, want no-store", overrideResponse.Header().Get("Cache-Control"))
	}

	if err := os.Remove(overridePath); err != nil {
		t.Fatal(err)
	}
	packagedResponse := httptest.NewRecorder()
	handler.ServeHTTP(packagedResponse, httptest.NewRequest("GET", "/audio/lifegame-theme.wav", nil))
	if packagedResponse.Code != 200 || strings.TrimSpace(packagedResponse.Body.String()) != "packaged-music" {
		t.Fatalf("packaged response = %d %q", packagedResponse.Code, packagedResponse.Body.String())
	}
}

func TestAssetHandlerUsesAntiqueFallback(t *testing.T) {
	packaged := fstest.MapFS{
		"images/antiqueinfo/default.webp": &fstest.MapFile{Data: []byte("packaged-fallback")},
	}
	handler := newAssetHandler(fs.FS(packaged))
	fallbackResponse := httptest.NewRecorder()
	handler.ServeHTTP(fallbackResponse, httptest.NewRequest("GET", "/images/antiqueinfo/missing.png", nil))
	if fallbackResponse.Code != 200 || strings.TrimSpace(fallbackResponse.Body.String()) != "packaged-fallback" {
		t.Fatalf("fallback response = %d %q", fallbackResponse.Code, fallbackResponse.Body.String())
	}
}

func TestEnsureExternalImagesExtractsMissingAndPreservesOverrides(t *testing.T) {
	packaged := fstest.MapFS{
		"images/portraits/01.webp": &fstest.MapFile{Data: []byte("default-portrait")},
		"images/scenes/park.webp":  &fstest.MapFile{Data: []byte("default-scene")},
		"index.html":               &fstest.MapFile{Data: []byte("frontend")},
	}
	imagesDir := filepath.Join(t.TempDir(), "images")

	extracted, err := ensureExternalImages(packaged, imagesDir)
	if err != nil {
		t.Fatalf("ensureExternalImages() error = %v", err)
	}
	if extracted != 2 {
		t.Fatalf("extracted files = %d, want 2", extracted)
	}

	portraitPath := filepath.Join(imagesDir, "portraits", "01.webp")
	if err := os.WriteFile(portraitPath, []byte("user-portrait"), 0o644); err != nil {
		t.Fatal(err)
	}
	scenePath := filepath.Join(imagesDir, "scenes", "park.webp")
	if err := os.Remove(scenePath); err != nil {
		t.Fatal(err)
	}

	extracted, err = ensureExternalImages(packaged, imagesDir)
	if err != nil {
		t.Fatalf("second ensureExternalImages() error = %v", err)
	}
	if extracted != 1 {
		t.Fatalf("second extracted files = %d, want 1", extracted)
	}
	portrait, err := os.ReadFile(portraitPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(portrait) != "user-portrait" {
		t.Fatalf("portrait = %q, want preserved user file", portrait)
	}
	scene, err := os.ReadFile(scenePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(scene) != "default-scene" {
		t.Fatalf("scene = %q, want restored default", scene)
	}
}

func TestEnsureExternalAudioExtractsMissingAndPreservesOverride(t *testing.T) {
	packaged := fstest.MapFS{
		"audio/lifegame-theme.wav": &fstest.MapFile{Data: []byte("default-music")},
	}
	audioDir := filepath.Join(t.TempDir(), "audio")

	extracted, err := ensureExternalAudio(packaged, audioDir)
	if err != nil {
		t.Fatalf("ensureExternalAudio() error = %v", err)
	}
	if extracted != 1 {
		t.Fatalf("extracted files = %d, want 1", extracted)
	}
	musicPath := filepath.Join(audioDir, "lifegame-theme.wav")
	if err := os.WriteFile(musicPath, []byte("user-music"), 0o644); err != nil {
		t.Fatal(err)
	}

	extracted, err = ensureExternalAudio(packaged, audioDir)
	if err != nil {
		t.Fatalf("second ensureExternalAudio() error = %v", err)
	}
	if extracted != 0 {
		t.Fatalf("second extracted files = %d, want 0", extracted)
	}
	music, err := os.ReadFile(musicPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(music) != "user-music" {
		t.Fatalf("music = %q, want preserved user file", music)
	}
}
