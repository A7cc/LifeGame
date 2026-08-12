package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	packagedImagesPath = "images"
	packagedAudioPath  = "audio"
)

// ensureExternalImages 将发布包中的默认图片释放到用户数据目录。
// 已存在的文件始终视为用户自定义内容，不会被默认资源覆盖。
func ensureExternalImages(packagedAssets fs.FS, imagesDir string) (int, error) {
	return ensureExternalAssetDirectory(packagedAssets, packagedImagesPath, imagesDir, "图片")
}

// ensureExternalAudio 将默认音频释放到用户数据目录，并保留用户同名替换。
func ensureExternalAudio(packagedAssets fs.FS, audioDir string) (int, error) {
	return ensureExternalAssetDirectory(packagedAssets, packagedAudioPath, audioDir, "音频")
}

func ensureExternalAssetDirectory(packagedAssets fs.FS, packagedRoot, targetDir, label string) (int, error) {
	if targetDir == "" {
		return 0, fmt.Errorf("%s目录不能为空", label)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return 0, fmt.Errorf("创建%s目录: %w", label, err)
	}

	extracted := 0
	err := fs.WalkDir(packagedAssets, packagedRoot, func(assetPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if assetPath == packagedRoot {
			return nil
		}

		relativePath := strings.TrimPrefix(assetPath, packagedRoot+"/")
		if relativePath == assetPath || !fs.ValidPath(relativePath) {
			return fmt.Errorf("无效的内置%s路径 %q", label, assetPath)
		}
		targetPath := filepath.Join(targetDir, filepath.FromSlash(relativePath))
		relativeTarget, err := filepath.Rel(targetDir, targetPath)
		if err != nil || relativeTarget == ".." || strings.HasPrefix(relativeTarget, ".."+string(filepath.Separator)) {
			return fmt.Errorf("%s路径越界 %q", label, assetPath)
		}

		if entry.IsDir() {
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return fmt.Errorf("创建%s子目录 %q: %w", label, relativePath, err)
			}
			return nil
		}

		created, err := copyEmbeddedFileIfMissing(packagedAssets, assetPath, targetPath)
		if err != nil {
			return fmt.Errorf("释放%s %q: %w", label, relativePath, err)
		}
		if created {
			extracted++
		}
		return nil
	})
	if err != nil {
		return extracted, err
	}
	return extracted, nil
}

func copyEmbeddedFileIfMissing(packagedAssets fs.FS, assetPath, targetPath string) (bool, error) {
	if info, err := os.Lstat(targetPath); err == nil {
		if info.IsDir() {
			return false, fmt.Errorf("目标位置是目录")
		}
		return false, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return false, err
	}
	source, err := packagedAssets.Open(assetPath)
	if err != nil {
		return false, err
	}
	defer source.Close()

	temporary, err := os.CreateTemp(filepath.Dir(targetPath), ".lifegame-asset-*")
	if err != nil {
		return false, err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if _, err := io.Copy(temporary, source); err != nil {
		temporary.Close()
		return false, err
	}
	if err := temporary.Chmod(0o644); err != nil {
		temporary.Close()
		return false, err
	}
	if err := temporary.Close(); err != nil {
		return false, err
	}

	// 再检查一次，避免释放期间出现的用户文件被覆盖。
	if _, err := os.Lstat(targetPath); err == nil {
		return false, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}
	if err := os.Rename(temporaryPath, targetPath); err != nil {
		return false, err
	}
	return true, nil
}
