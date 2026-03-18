package handler

import "strings"

const projectURL = "https://stickystick.ru"
const telegramChannelURL = "https://t.me/uncensored_mems"
const vkGroupURL = "https://vk.com/club236352692"

func normalizeLines(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

// ensureLinksFirst гарантирует, что указанные ссылки будут первыми строками (в заданном порядке),
// без дублей. Остальной текст идёт ниже.
func ensureLinksFirst(text string, links []string) string {
	existing := normalizeLines(text)
	seen := make(map[string]bool, len(existing))
	for _, l := range existing {
		seen[l] = true
	}

	var out []string
	for _, link := range links {
		link = strings.TrimSpace(link)
		if link == "" {
			continue
		}
		out = append(out, link)
		seen[link] = true
	}

	for _, l := range existing {
		if seen[l] {
			// если строка уже добавлена в префикс — пропускаем
			continue
		}
		out = append(out, l)
	}

	return strings.Join(out, "\n")
}

