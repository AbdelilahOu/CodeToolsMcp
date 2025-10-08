package runners

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type FileRunner struct{}

type DirEntry struct {
	Path    string
	Name    string
	IsDir   bool
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time
}

func NewFileRunner() *FileRunner {
	return &FileRunner{}
}

type ListDirInput struct {
	Path       string
	Recursive  bool
	ShowHidden bool
	Limit      int
}

func (r *FileRunner) ListDir(ctx context.Context, input ListDirInput) ([]DirEntry, error) {
	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", absPath)
	}

	var entries []DirEntry
	limit := input.Limit

	collect := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if path == absPath {
			return nil
		}

		name := d.Name()
		if !input.ShowHidden && strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		dirEntry := DirEntry{
			Path:    path,
			Name:    name,
			IsDir:   d.IsDir(),
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
		}
		entries = append(entries, dirEntry)

		if limit > 0 && len(entries) >= limit {
			return context.Canceled
		}

		if !input.Recursive && d.IsDir() {
			return fs.SkipDir
		}

		return nil
	}

	var walkErr error
	if input.Recursive {
		walkErr = filepath.WalkDir(absPath, collect)
	} else {
		files, err := os.ReadDir(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory: %w", err)
		}
		for _, entry := range files {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			name := entry.Name()
			if !input.ShowHidden && strings.HasPrefix(name, ".") {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			dirEntry := DirEntry{
				Path:    filepath.Join(absPath, name),
				Name:    name,
				IsDir:   entry.IsDir(),
				Size:    info.Size(),
				Mode:    info.Mode(),
				ModTime: info.ModTime(),
			}
			entries = append(entries, dirEntry)
			if limit > 0 && len(entries) >= limit {
				break
			}
		}
	}

	if walkErr != nil && !errors.Is(walkErr, context.Canceled) {
		return nil, fmt.Errorf("failed to walk directory: %w", walkErr)
	}

	return entries, nil
}

type DeleteInput struct {
	Path string
}

func (r *FileRunner) Delete(ctx context.Context, input DeleteInput) error {
	info, err := os.Stat(input.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", input.Path)
		}
		return fmt.Errorf("failed to stat path: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("delete expects a file but directory was provided: %s", input.Path)
	}

	return os.Remove(input.Path)
}

type RemoveInput struct {
	Path      string
	Recursive bool
}

func (r *FileRunner) Remove(ctx context.Context, input RemoveInput) error {
	if input.Recursive {
		return os.RemoveAll(input.Path)
	}
	return os.Remove(input.Path)
}

type CopyInput struct {
	Source    string
	Target    string
	Overwrite bool
}

func (r *FileRunner) Copy(ctx context.Context, input CopyInput) error {
	srcAbs, err := filepath.Abs(input.Source)
	if err != nil {
		return fmt.Errorf("failed to resolve source: %w", err)
	}
	dstAbs, err := filepath.Abs(input.Target)
	if err != nil {
		return fmt.Errorf("failed to resolve target: %w", err)
	}

	srcInfo, err := os.Lstat(srcAbs)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	dstInfo, err := os.Stat(dstAbs)
	if err == nil {
		if !input.Overwrite {
			return fmt.Errorf("target already exists: %s", dstAbs)
		}
		if dstInfo.IsDir() {
			if err := os.RemoveAll(dstAbs); err != nil {
				return fmt.Errorf("failed to clear target: %w", err)
			}
		} else {
			if err := os.Remove(dstAbs); err != nil {
				return fmt.Errorf("failed to overwrite target: %w", err)
			}
		}
	}

	if srcInfo.IsDir() {
		return copyDirectory(srcAbs, dstAbs)
	}

	return copyFile(srcAbs, dstAbs, srcInfo.Mode())
}

type MoveInput struct {
	Source    string
	Target    string
	Overwrite bool
}

func (r *FileRunner) Move(ctx context.Context, input MoveInput) error {
	srcAbs, err := filepath.Abs(input.Source)
	if err != nil {
		return fmt.Errorf("failed to resolve source: %w", err)
	}
	dstAbs, err := filepath.Abs(input.Target)
	if err != nil {
		return fmt.Errorf("failed to resolve target: %w", err)
	}

	if _, err := os.Lstat(srcAbs); err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	if info, err := os.Stat(dstAbs); err == nil {
		if !input.Overwrite {
			return fmt.Errorf("target already exists: %s", dstAbs)
		}
		if info.IsDir() {
			if err := os.RemoveAll(dstAbs); err != nil {
				return fmt.Errorf("failed to clear target directory: %w", err)
			}
		} else {
			if err := os.Remove(dstAbs); err != nil {
				return fmt.Errorf("failed to remove target file: %w", err)
			}
		}
	}

	err = os.Rename(srcAbs, dstAbs)
	if err == nil {
		return nil
	}

	var linkErr *os.LinkError
	if errors.As(err, &linkErr) && errors.Is(linkErr.Err, syscall.EXDEV) {
		if copyErr := copyDirectoryOrFile(srcAbs, dstAbs); copyErr != nil {
			return fmt.Errorf("failed to move across filesystems: %w", copyErr)
		}
		if removeErr := os.RemoveAll(srcAbs); removeErr != nil {
			return fmt.Errorf("moved but failed to remove source: %w", removeErr)
		}
		return nil
	}

	return fmt.Errorf("failed to move path: %w", err)
}

type TreeInput struct {
	Path       string
	Depth      int
	ShowHidden bool
	Limit      int
}

func (r *FileRunner) Tree(ctx context.Context, input TreeInput) (string, error) {
	absPath, err := filepath.Abs(input.Path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", absPath)
	}

	var builder strings.Builder
	builder.WriteString(absPath)
	builder.WriteString("\n")

	limit := input.Limit
	var count int

	var walk func(string, string, int) error

	walk = func(current string, prefix string, depth int) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if input.Depth > 0 && depth > input.Depth {
			return nil
		}

		entries, err := os.ReadDir(current)
		if err != nil {
			return nil
		}

		filtered := make([]os.DirEntry, 0, len(entries))
		for _, entry := range entries {
			name := entry.Name()
			if !input.ShowHidden && strings.HasPrefix(name, ".") {
				continue
			}
			filtered = append(filtered, entry)
		}

		for i, entry := range filtered {
			if limit > 0 && count >= limit {
				return context.Canceled
			}

			connector := "|-- "
			childPrefix := prefix + "|   "
			if i == len(filtered)-1 {
				connector = "\\-- "
				childPrefix = prefix + "    "
			}

			line := prefix + connector + entry.Name()
			if entry.IsDir() {
				line += "/"
			}
			builder.WriteString(line)
			builder.WriteString("\n")
			count++

			if entry.IsDir() {
				nextPath := filepath.Join(current, entry.Name())
				if err := walk(nextPath, childPrefix, depth+1); err != nil {
					if errors.Is(err, context.Canceled) {
						return err
					}
				}
			}
		}

		return nil
	}

	err = walk(absPath, "", 1)
	if err != nil && !errors.Is(err, context.Canceled) {
		return "", err
	}

	return strings.TrimRight(builder.String(), "\n"), nil
}

func copyDirectoryOrFile(src, dst string) error {
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return copyDirectory(src, dst)
	}

	return copyFile(src, dst, srcInfo.Mode())
}

func copyDirectory(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if d.IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			return nil
		}

		srcInfo, err := os.Lstat(path)
		if err != nil {
			return err
		}

		return copyFile(path, target, srcInfo.Mode())
	})
}

func copyFile(src, dst string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
