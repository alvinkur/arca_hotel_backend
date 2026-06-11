package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"arca-hotel/services/ai-service/clients"

	"github.com/gin-gonic/gin"
)

type aiChatRequest struct {
	Model    string      `json:"model"`
	Messages []aiMessage `json:"messages"`
}

type aiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type aiChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type recommendRequest struct {
	Message string `json:"message" binding:"required"`
}

var roomClient *clients.RoomClient

func initClient() {
	if roomClient == nil {
		base := os.Getenv("ROOM_SERVICE_URL")
		if base == "" {
			base = "http://localhost:8002"
		}
		roomClient = clients.NewRoomClient(base)
	}
}

func RecommendRoom(c *gin.Context) {
	initClient()

	var req recommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reply": "Tolong tuliskan preferensi kamar yang Anda inginkan."})
		return
	}

	roomTypes, err := roomClient.GetRoomTypes()
	if err != nil || len(roomTypes) == 0 {
		c.JSON(http.StatusOK, gin.H{"reply": "Maaf, saat ini belum ada tipe kamar yang tersedia."})
		return
	}

	apiKey := os.Getenv("AI_API_KEY")
	if apiKey != "" {
		reply, err := callAI(apiKey, roomTypes, req.Message)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"reply": reply})
			return
		}
	}

	reply := keywordMatch(roomTypes, req.Message)
	c.JSON(http.StatusOK, gin.H{"reply": reply})
}

func callAI(apiKey string, roomTypes []clients.RoomTypeDTO, userMessage string) (string, error) {
	baseURL := os.Getenv("AI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	var sb strings.Builder
	for _, rt := range roomTypes {
		fmt.Fprintf(&sb, "- %s: Rp%.0f/malam. %s\n", rt.Name, rt.Price, rt.Description)
	}

	systemPrompt := fmt.Sprintf(
		"Kamu adalah asisten Hotel Arca yang membantu tamu memilih kamar. "+
			"Berikut daftar tipe kamar yang tersedia:\n\n%s\n"+
			"Bantu tamu memilih kamar berdasarkan preferensi mereka. "+
			"Berikan rekomendasi singkat (2-3 kalimat) dalam Bahasa Indonesia, sebutkan nama kamar dan harganya.",
		sb.String(),
	)

	body := aiChatRequest{
		Model: model,
		Messages: []aiMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
	}

	jsonBody, _ := json.Marshal(body)
	httpReq, _ := http.NewRequest("POST",
		strings.TrimRight(baseURL, "/")+"/chat/completions",
		bytes.NewReader(jsonBody),
	)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ai api returned %d", resp.StatusCode)
	}

	var aiResp aiChatResponse
	if err := json.Unmarshal(respBytes, &aiResp); err != nil {
		return "", err
	}
	if len(aiResp.Choices) == 0 {
		return "", fmt.Errorf("empty ai response")
	}

	return aiResp.Choices[0].Message.Content, nil
}

var stopWords = map[string]bool{
	"saya": true, "ingin": true, "mau": true, "cari": true, "butuh": true,
	"yang": true, "dan": true, "atau": true, "di": true, "ke": true,
	"buat": true, "untuk": true, "dengan": true, "bisa": true, "ada": true,
	"kamar": true, "hotel": true, "arca": true, "tidak": true, "juga": true,
	"tolong": true, "si": true, "nih": true, "dong": true, "ya": true,
	"gak": true, "nggak": true, "aja": true, "deh": true, "kok": true,
	"itu": true, "ini": true, "aku": true, "gua": true, "gue": true,
}

var wordWeights = map[string]int{
	"murah": 5, "ekonomis": 5, "hemat": 5, "terjangkau": 5, "budget": 5,
	"mewah": 5, "luxury": 5, "premium": 5, "eksklusif": 5, "vip": 5,
	"mahal": 5, "high": 5,
	"kolam": 4, "renang": 4, "bathtub": 3, "balkon": 3, "pemandangan": 3,
	"luas": 3, "lapang": 3, "besar": 3,
	"nyaman": 2, "bersih": 2, "ac": 2, "tv": 1,
	"tenang": 2, "sunyi": 2, "privacy": 2, "pribadi": 2,
}

var synonyms = map[string][]string{
	"murah":      {"murah", "ekonomis", "hemat", "terjangkau", "budget", "low cost", "miring"},
	"mewah":      {"mewah", "luxury", "premium", "eksklusif", "vip", "elite"},
	"luas":       {"luas", "lapang", "besar", "spacious", "legawa"},
	"nyaman":     {"nyaman", "comfortable", "enak", "cozy", "hangat"},
	"tenang":     {"tenang", "sunyi", "quiet", "damai", "peaceful"},
	"keluarga":   {"keluarga", "family", "anak", "ramai", "rombongan"},
	"pasangan":   {"pasangan", "honeymoon", "romantis", "couple", "bulan madu"},
	"bisnis":     {"bisnis", "business", "kerja", "meeting", "rapat"},
	"sendiri":    {"sendiri", "solo", "single", "personal"},
	"menginap":   {"menginap", "inap", "tidur", "stay", "istirahat"},
	"liburan":    {"liburan", "vacation", "holiday", "healing", "staycation", "piknik"},
	"pemandangan": {"pemandangan", "view", "panorama", "scenery", "lanskap"},
	"pelayanan":  {"pelayanan", "service", "layanan", "fasilitas"},
}

var intentTemplates = map[string]map[string]int{
	"honeymoon": {"Suite": 8, "Deluxe": 6},
	"bulan madu": {"Suite": 8, "Deluxe": 6},
	"romantis":  {"Suite": 7, "Deluxe": 5},
	"keluarga":  {"Standard": 5, "Deluxe": 4},
	"family":    {"Standard": 5, "Deluxe": 4},
	"bisnis":    {"Deluxe": 5, "Standard": 3},
	"business":  {"Deluxe": 5, "Standard": 3},
	"mewah":     {"Suite": 8, "Deluxe": 4},
	"luxury":    {"Suite": 8, "Deluxe": 4},
	"murah":     {"Standard": 6},
	"budget":    {"Standard": 6},
	"ekonomis":  {"Standard": 6},
	"hemat":     {"Standard": 6},
	"kolam":     {"Suite": 6},
	"renang":    {"Suite": 6},
}

func keywordMatch(roomTypes []clients.RoomTypeDTO, message string) string {
	lower := strings.ToLower(strings.TrimSpace(message))

	templateBonuses := make(map[string]int)
	for phrase, bonuses := range intentTemplates {
		if strings.Contains(lower, phrase) {
			for roomName, bonus := range bonuses {
				templateBonuses[roomName] += bonus
			}
		}
	}

	rawWords := strings.Fields(lower)
	var expanded []string
	seen := map[string]bool{}
	for _, w := range rawWords {
		if stopWords[w] {
			continue
		}
		syns, ok := synonyms[w]
		if !ok {
			if seen[w] {
				continue
			}
			seen[w] = true
			expanded = append(expanded, w)
			continue
		}
		for _, s := range syns {
			if seen[s] {
				continue
			}
			seen[s] = true
			expanded = append(expanded, s)
		}
	}

	if len(expanded) == 0 {
		return "Silakan ceritakan preferensi kamar Anda. Contoh: 'Saya ingin kamar murah dekat kolam renang'."
	}

	type scored struct {
		roomType clients.RoomTypeDTO
		score    float64
	}
	var results []scored

	for _, rt := range roomTypes {
		rtText := strings.ToLower(rt.Name + " " + rt.Description)
		score := 0.0

		for _, kw := range expanded {
			if strings.Contains(rtText, kw) {
				w := wordWeights[kw]
				if w == 0 {
					w = 1
				}
				score += float64(w)
			}
		}

		if bonus, ok := templateBonuses[rt.Name]; ok {
			score += float64(bonus)
		}
		if bonus, ok := templateBonuses[strings.ToLower(rt.Name)]; ok {
			score += float64(bonus)
		}

		results = append(results, scored{rt, score})
	}

	priceFactor := 0
	if containsAny(lower, []string{"murah", "hemat", "ekonomis", "terjangkau", "budget", "low cost"}) {
		priceFactor = -1
	} else if containsAny(lower, []string{"mewah", "luxury", "mahal", "premium", "eksklusif", "vip"}) {
		priceFactor = 1
	}

	if priceFactor != 0 {
		prices := make([]float64, len(roomTypes))
		for i, rt := range roomTypes {
			prices[i] = rt.Price
		}
		minP, maxP := minMax(prices)

		for i := range results {
			if maxP > minP {
				norm := (results[i].roomType.Price - minP) / (maxP - minP)
				results[i].score += float64(priceFactor) * norm * 4
			}
		}
	}

	best := results[0]
	for _, r := range results[1:] {
		if r.score > best.score {
			best = r
		}
	}

	if best.score == 0 {
		bestP := roomTypes[0]
		for _, rt := range roomTypes {
			if rt.Price < bestP.Price {
				bestP = rt
			}
		}
		best = scored{bestP, 1}
	}

	return fmt.Sprintf("Berdasarkan preferensi Anda, kami merekomendasikan kamar **%s** dengan harga Rp%.0f/malam. %s",
		best.roomType.Name, best.roomType.Price, best.roomType.Description)
}

func containsAny(s string, terms []string) bool {
	for _, t := range terms {
		if strings.Contains(s, t) {
			return true
		}
	}
	return false
}

func minMax(vals []float64) (float64, float64) {
	min, max := vals[0], vals[0]
	for _, v := range vals[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}
