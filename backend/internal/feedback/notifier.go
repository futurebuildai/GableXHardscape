package feedback

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Notifier sends feedback notifications to Google Chat via webhook.
type Notifier struct {
	webhookURL string
	httpClient *http.Client
	logger     *slog.Logger
}

// NewNotifier creates a Google Chat notifier. The webhook URL is read from
// the GOOGLE_CHAT_WEBHOOK_URL environment variable.
func NewNotifier(logger *slog.Logger) *Notifier {
	return &Notifier{
		webhookURL: os.Getenv("GOOGLE_CHAT_WEBHOOK_URL"),
		httpClient: &http.Client{Timeout: 10 * time.Second},
		logger:     logger,
	}
}

// NotifyNewFeedback sends a card message to Google Chat for a new feedback
// submission. Runs asynchronously — errors are logged, never returned.
func (n *Notifier) NotifyNewFeedback(fb *Feedback) {
	if n.webhookURL == "" {
		n.logger.Warn("Feedback notifier: GOOGLE_CHAT_WEBHOOK_URL not set, skipping notification")
		return
	}

	go func() {
		if err := n.send(fb); err != nil {
			n.logger.Error("Feedback notifier: failed to send Google Chat notification",
				"feedback_id", fb.ID,
				"error", err,
			)
		}
	}()
}

func (n *Notifier) send(fb *Feedback) error {
	// Truncate long descriptions for the card preview.
	desc := fb.Description
	if len(desc) > 300 {
		desc = desc[:297] + "..."
	}

	categoryEmoji := map[string]string{
		"Bug":             "🐛",
		"UI/UX":           "🎨",
		"Feature Request": "💡",
		"Data Issue":      "📊",
		"Question":        "❓",
		"Other":           "📝",
	}
	emoji := categoryEmoji[fb.Category]
	if emoji == "" {
		emoji = "📝"
	}

	submitter := fb.SubmittedByName
	if submitter == "" {
		submitter = fb.SubmittedByEmail
	}
	if submitter == "" {
		submitter = "Anonymous"
	}

	subtitle := fmt.Sprintf("From: %s via %s", submitter, fb.Source)

	payload := map[string]interface{}{
		"cardsV2": []map[string]interface{}{
			{
				"cardId": fb.ID.String(),
				"card": map[string]interface{}{
					"header": map[string]interface{}{
						"title":    fmt.Sprintf("%s New Feedback — %s", emoji, fb.Category),
						"subtitle": subtitle,
					},
					"sections": []map[string]interface{}{
						{
							"widgets": []map[string]interface{}{
								{
									"decoratedText": map[string]interface{}{
										"topLabel": "Title",
										"text":     fb.Title,
									},
								},
								{
									"decoratedText": map[string]interface{}{
										"topLabel": "Description",
										"text":     desc,
									},
								},
								{
									"decoratedText": map[string]interface{}{
										"topLabel": "Page",
										"text":     fb.PageURL,
									},
								},
								{
									"decoratedText": map[string]interface{}{
										"topLabel": "Priority",
										"text":     fb.Priority,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Google Chat returned status %d", resp.StatusCode)
	}

	n.logger.Info("Feedback notifier: Google Chat notification sent",
		"feedback_id", fb.ID,
		"category", fb.Category,
	)
	return nil
}
