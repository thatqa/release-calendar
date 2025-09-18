package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"release-calendar/backend/internal/models"
	"release-calendar/backend/internal/repostiry"

	"gorm.io/gorm"
)

type ReleaseSummarizer struct {
	releaseRepo repostiry.ReleaseRepository
	commentRepo repostiry.CommentRepository
	httpClient  *http.Client
}

func NewReleaseSummarizer(db *gorm.DB) ReleaseSummarizer {
	return ReleaseSummarizer{
		releaseRepo: repostiry.NewReleaseRepository(db),
		commentRepo: repostiry.NewCommentRepository(db),
		httpClient:  &http.Client{Timeout: 20 * time.Second},
	}
}

type summaryResponse struct {
	Summary  string `json:"summary"`
	Provider string `json:"provider"`
}

func (u *ReleaseSummarizer) Summarize(releaseID uint) (summaryResponse, error) {
	rel, err := u.releaseRepo.GetById(releaseID)
	if err != nil {
		return summaryResponse{}, fmt.Errorf("release not found")
	}
	comments, err := u.commentRepo.GetByReleaseID(releaseID)
	if err != nil {
		return summaryResponse{}, fmt.Errorf("failed to load comments")
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Title: %s\nDate: %s\nStatus: %s\n", rel.Title, rel.Date.Format(time.RFC3339), rel.Status)
	if strings.TrimSpace(rel.Notes) != "" {
		fmt.Fprintf(&b, "\nNotes:\n%s\n", rel.Notes)
	}
	if len(rel.Links) > 0 {
		fmt.Fprintf(&b, "\nLinks:\n")
		for _, l := range rel.Links {
			fmt.Fprintf(&b, "- %s: %s\n", l.Name, l.URL)
		}
	}
	if len(comments) > 0 {
		fmt.Fprintf(&b, "\nComments (%d):\n", len(comments))
		for _, c := range comments {
			msg := c.Message
			if len(msg) > 500 {
				msg = msg[:500] + "…"
			}
			fmt.Fprintf(&b, "- %s: %s\n", c.Author, msg)
		}
	}
	raw := b.String()

	if apiKey, ok := os.LookupEnv("AI_API_KEY"); ok && apiKey != "" {
		text, err := u.callOpenAI(apiKey, raw)
		if err == nil && strings.TrimSpace(text) != "" {
			return summaryResponse{Summary: text, Provider: "openai"}, nil
		}
	}

	return summaryResponse{Summary: fallbackSummarize(rel, comments), Provider: "local"}, nil
}

func (u *ReleaseSummarizer) callOpenAI(apiKey, raw string) (string, error) {
	temperature, _ := strconv.ParseFloat(os.Getenv("AI_TEMPERATURE"), 64)
	maxTokens, _ := strconv.ParseInt(os.Getenv("AI_MAX_TOKENS"), 10, 64)
	payload := map[string]any{
		"model":       os.Getenv("AI_MODEL"),
		"temperature": temperature,
		"max_tokens":  maxTokens,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful release manager assistant. Summarize briefly and actionably."},
			{"role": "user", "content": "Create a concise, helpful summary for a release based on notes and comments. Include: 1) current status and date/time, 2) main issues/risks, 3) decisions/actions (bulleted), 4) next steps. Keep it under 8 bullet points. Text:\n\n" + raw},
		},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", os.Getenv("AI_URL"), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		log.Printf(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	var r struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Printf(err.Error())
		return "", err
	}

	log.Printf(string(len(r.Choices)))

	if len(r.Choices) == 0 {
		return "", fmt.Errorf("no choices")
	}
	return strings.TrimSpace(r.Choices[0].Message.Content), nil
}

func fallbackSummarize(r *models.Release, cs []*models.Comment) string {
	var items []string

	items = append(items, fmt.Sprintf("**Status**: %s · %s", r.Status, r.Date.Format("02 Jan 2006 15:04")))

	var txt strings.Builder
	if s := strings.TrimSpace(r.Notes); s != "" {
		txt.WriteString(s)
		txt.WriteString("\n")
	}
	for _, c := range cs {
		txt.WriteString(c.Message)
		txt.WriteString("\n")
	}
	sentences := splitSentences(txt.String())

	keywords := []string{"fail", "error", "block", "risk", "issue", "bug", "fixed", "success", "ready", "deploy", "rollback", "mitigate", "investigate", "monitor"}
	scores := make([]struct {
		s string
		w int
	}, 0, len(sentences))

	for _, s := range sentences {
		ss := strings.ToLower(s)
		w := 0
		for _, kw := range keywords {
			if strings.Contains(ss, kw) {
				w += 3
			}
		}
		if strings.Contains(ss, "!") {
			w += 1
		}
		if l := len(s); l >= 40 && l <= 220 {
			w += 1
		}
		if w > 0 {
			scores = append(scores, struct {
				s string
				w int
			}{s: strings.TrimSpace(s), w: w})
		}
	}
	sort.Slice(scores, func(i, j int) bool { return scores[i].w > scores[j].w })

	max := 6
	for i := 0; i < len(scores) && i < max; i++ {
		items = append(items, "• "+scores[i].s)
	}
	if len(scores) == 0 {
		items = append(items, "• No critical issues mentioned.")
	}

	//hasRisk := false
	//for _, s := range sentences {
	//	ss := strings.ToLower(s)
	//	if strings.Contains(ss, "risk") || strings.Contains(ss, "fail") || strings.Contains(ss, "issue") || strings.Contains(ss, "error") || strings.Contains(ss, "bug") {
	//		hasRisk = true
	//		break
	//	}
	//}
	return strings.Join(items, "\n")
}

func splitSentences(text string) []string {
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.Join(strings.Fields(text), " ")

	var out []string
	seen := make(map[string]bool)

	var curr []rune
	for _, r := range text {
		// пропускаем лидирующие пробелы у нового предложения
		if len(curr) == 0 && (r == ' ' || r == '\t') {
			continue
		}
		curr = append(curr, r)

		// конец предложения
		if r == '.' || r == '!' || r == '?' {
			s := strings.TrimSpace(string(curr))
			if len(s) >= 8 {
				key := strings.ToLower(s)
				if !seen[key] {
					out = append(out, s)
					seen[key] = true
				}
			}
			curr = curr[:0]
		}
	}

	if len(curr) > 0 {
		s := strings.TrimSpace(string(curr))
		if len(s) >= 8 {
			key := strings.ToLower(s)
			if !seen[key] {
				out = append(out, s)
				seen[key] = true
			}
		}
	}

	return out
}
