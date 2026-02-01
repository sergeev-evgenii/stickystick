package service

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"image/jpeg"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

type MediaService interface {
	SaveFile(file *multipart.FileHeader, mediaType string) (string, error)
	GetFileURL(filename string) string
	DeleteFile(filename string) error
	ValidateFileType(file *multipart.FileHeader, allowedTypes []string) error
	GetMediaType(filename string) string
}

type mediaService struct {
	uploadDir string
	baseURL   string
}

func NewMediaService(uploadDir, baseURL string) MediaService {
	// Создаем директорию для загрузок, если её нет
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	// Создаем поддиректории для разных типов медиа
	os.MkdirAll(filepath.Join(uploadDir, "videos"), 0755)
	os.MkdirAll(filepath.Join(uploadDir, "photos"), 0755)
	os.MkdirAll(filepath.Join(uploadDir, "gifs"), 0755)

	return &mediaService{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

func (s *mediaService) GetMediaType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".mp4", ".mov", ".avi", ".webm", ".mkv":
		return "video"
	case ".gif":
		return "gif"
	case ".jpg", ".jpeg", ".png", ".webp":
		return "photo"
	default:
		return "unknown"
	}
}

func (s *mediaService) ValidateFileType(file *multipart.FileHeader, allowedTypes []string) error {
	mediaType := s.GetMediaType(file.Filename)
	
	for _, allowed := range allowedTypes {
		if mediaType == allowed {
			return nil
		}
	}
	
	return fmt.Errorf("file type not allowed. Allowed types: %v", allowedTypes)
}

func (s *mediaService) SaveFile(file *multipart.FileHeader, mediaType string) (string, error) {
	// Определяем директорию и расширение для сохранения
	var subDir string
	var finalExt string
	switch mediaType {
	case "video":
		subDir = "videos"
		finalExt = ".mp4" // Все видео конвертируем в MP4
	case "photo":
		subDir = "photos"
		finalExt = ".jpg" // Все фото конвертируем в JPEG
	case "gif":
		subDir = "gifs"
		finalExt = ".gif" // GIF оставляем как есть
	default:
		subDir = "misc"
		finalExt = filepath.Ext(file.Filename)
	}

	// Генерируем уникальное имя файла
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d%s", timestamp, finalExt)
	
	// Полный путь для сохранения
	fullPath := filepath.Join(s.uploadDir, subDir, filename)

	// Обрабатываем в зависимости от типа
	switch mediaType {
	case "photo":
		return s.saveAndCompressImage(file, fullPath, subDir, filename)
	case "video":
		return s.saveAndCompressVideo(file, fullPath, subDir, filename)
	case "gif":
		return s.saveAndCompressGif(file, fullPath, subDir, filename)
	default:
		return s.saveFileDirect(file, fullPath, subDir, filename)
	}
}

// saveFileDirect просто сохраняет файл без обработки
func (s *mediaService) saveFileDirect(file *multipart.FileHeader, fullPath, subDir, filename string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filepath.Join(subDir, filename), nil
}

// saveAndCompressImage сохраняет и сжимает изображение
func (s *mediaService) saveAndCompressImage(file *multipart.FileHeader, fullPath, subDir, filename string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Пробуем сжать изображение
	if err := s.compressImage(src, fullPath); err != nil {
		// Если сжатие не удалось, сохраняем оригинал
		src2, err2 := file.Open()
		if err2 != nil {
			return "", fmt.Errorf("failed to open file for fallback: %w", err2)
		}
		defer src2.Close()

		dst, err2 := os.Create(fullPath)
		if err2 != nil {
			return "", fmt.Errorf("failed to create file: %w", err2)
		}
		defer dst.Close()

		if _, err2 := io.Copy(dst, src2); err2 != nil {
			return "", fmt.Errorf("failed to save file: %w", err2)
		}
	}

	return filepath.Join(subDir, filename), nil
}

// saveAndCompressVideo сохраняет и сжимает видео
func (s *mediaService) saveAndCompressVideo(file *multipart.FileHeader, fullPath, subDir, filename string) (string, error) {
	// Сохраняем оригинал во временный файл
	timestamp := time.Now().UnixNano()
	tempPath := filepath.Join(s.uploadDir, subDir, fmt.Sprintf("temp_%d%s", timestamp, filepath.Ext(file.Filename)))

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := io.Copy(tempFile, src); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return "", fmt.Errorf("failed to save temp file: %w", err)
	}
	tempFile.Close()

	// Пробуем сжать видео
	if err := s.compressVideo(tempPath, fullPath); err != nil {
		// Если сжатие не удалось, копируем оригинал
		if err := s.copyFile(tempPath, fullPath); err != nil {
			os.Remove(tempPath)
			return "", fmt.Errorf("failed to copy file: %w", err)
		}
	}

	// Удаляем временный файл
	os.Remove(tempPath)

	return filepath.Join(subDir, filename), nil
}

// saveAndCompressGif сохраняет и сжимает GIF
func (s *mediaService) saveAndCompressGif(file *multipart.FileHeader, fullPath, subDir, filename string) (string, error) {
	// Сохраняем оригинал во временный файл
	timestamp := time.Now().UnixNano()
	tempPath := filepath.Join(s.uploadDir, subDir, fmt.Sprintf("temp_%d%s", timestamp, filepath.Ext(file.Filename)))

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := io.Copy(tempFile, src); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return "", fmt.Errorf("failed to save temp file: %w", err)
	}
	tempFile.Close()

	// Пробуем сжать GIF
	if err := s.compressGif(tempPath, fullPath); err != nil {
		// Если сжатие не удалось, копируем оригинал
		if err := s.copyFile(tempPath, fullPath); err != nil {
			os.Remove(tempPath)
			return "", fmt.Errorf("failed to copy file: %w", err)
		}
	}

	// Удаляем временный файл
	os.Remove(tempPath)

	return filepath.Join(subDir, filename), nil
}

// compressImage сжимает изображение
func (s *mediaService) compressImage(src io.Reader, outputPath string) error {
	// Декодируем изображение
	img, _, err := image.Decode(src)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Максимальные размеры (1920x1080)
	maxWidth := 1920
	maxHeight := 1080

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Уменьшаем размер, если нужно
	var resizedImg image.Image = img
	if width > maxWidth || height > maxHeight {
		resizedImg = imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
	} else {
		resizedImg = imaging.Clone(img)
	}

	// Сохраняем как JPEG с качеством 85%
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 85})
}

// compressVideo сжимает видео с помощью ffmpeg
func (s *mediaService) compressVideo(inputPath, outputPath string) error {
	// Проверяем наличие ffmpeg
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found")
	}

	// Команда ffmpeg для сжатия видео
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c:v", "libx264",
		"-preset", "medium",
		"-crf", "23",
		"-maxrate", "2M",
		"-bufsize", "4M",
		"-c:a", "aac",
		"-b:a", "128k",
		"-movflags", "+faststart",
		"-y",
		outputPath,
	)

	// Захватываем stderr для логирования ошибок
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg compression failed: %w", err)
	}

	return nil
}

// compressGif сжимает GIF с помощью gifsicle или ffmpeg
func (s *mediaService) compressGif(inputPath, outputPath string) error {
	// Сначала пробуем gifsicle (лучше для GIF)
	if _, err := exec.LookPath("gifsicle"); err == nil {
		cmd := exec.Command("gifsicle",
			"--optimize=3",
			"--colors", "256",
			"-o", outputPath,
			inputPath,
		)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Если gifsicle не доступен, пробуем ffmpeg
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		cmd := exec.Command("ffmpeg",
			"-i", inputPath,
			"-vf", "fps=10,scale=800:-1:flags=lanczos",
			"-y",
			outputPath,
		)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no compression tool available")
}

// copyFile копирует файл
func (s *mediaService) copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (s *mediaService) GetFileURL(filename string) string {
	if filename == "" {
		return ""
	}
	// Заменяем обратные слеши на прямые для URL
	urlPath := strings.ReplaceAll(filename, "\\", "/")
	return fmt.Sprintf("%s/uploads/%s", s.baseURL, urlPath)
}

func (s *mediaService) DeleteFile(filename string) error {
	if filename == "" {
		return nil
	}
	
	fullPath := filepath.Join(s.uploadDir, filename)
	return os.Remove(fullPath)
}
