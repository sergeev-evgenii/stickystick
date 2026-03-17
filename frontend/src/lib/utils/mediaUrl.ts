/**
 * Преобразует media_url из БД в URL для отображения.
 * В БД хранится относительный путь (videos/xxx.mp4), nginx добавляет /uploads/.
 */
export function getMediaSrc(url: string | undefined | null): string {
  if (!url) return ''
  // Если уже полный путь с /uploads/ — используем как есть
  if (url.startsWith('/')) return url
  return `/uploads/${url}`
}
