package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"arca-hotel/services/ai-service/clients"

	"github.com/gin-gonic/gin"
)

type recommendRequest struct {
	Message string `json:"message" binding:"required"`
}

var roomClient *clients.RoomClient
var mlClient *clients.MLClient

func initClients() {
	if roomClient == nil {
		base := os.Getenv("ROOM_SERVICE_URL")
		if base == "" {
			base = "http://localhost:8002"
		}
		roomClient = clients.NewRoomClient(base)
	}
	if mlClient == nil {
		base := os.Getenv("ML_SERVICE_URL")
		if base == "" {
			base = "http://localhost:8010"
		}
		mlClient = clients.NewMLClient(base)
	}
}

func formatRoomNumbers(nums []string) string {
	if len(nums) == 0 {
		return ""
	}
	return " — Room " + strings.Join(nums, ", ")
}

func RecommendRoom(c *gin.Context) {
	initClients()

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

	result, err := mlClient.Recommend(req.Message)
	if err == nil {
		rooms := formatRoomNumbers(result.RoomNumbers)
		reply := fmt.Sprintf("Berdasarkan preferensi Anda, kami merekomendasikan tipe **%s**%s dengan harga Rp%.0f/malam. %s",
			result.Name, rooms, result.Price, result.Description)
		c.JSON(http.StatusOK, gin.H{"reply": reply})
		return
	}

	rooms, _ := roomClient.GetRooms()
	reply := keywordMatch(roomTypes, rooms, req.Message)
	c.JSON(http.StatusOK, gin.H{"reply": reply})
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
	"murah":       {"murah", "ekonomis", "hemat", "terjangkau", "budget", "low cost", "miring"},
	"mewah":       {"mewah", "luxury", "premium", "eksklusif", "vip", "elite"},
	"luas":        {"luas", "lapang", "besar", "spacious", "legawa"},
	"nyaman":      {"nyaman", "comfortable", "enak", "cozy", "hangat"},
	"tenang":      {"tenang", "sunyi", "quiet", "damai", "peaceful"},
	"keluarga":    {"keluarga", "family", "anak", "ramai", "rombongan"},
	"pasangan":    {"pasangan", "honeymoon", "romantis", "couple", "bulan madu"},
	"bisnis":      {"bisnis", "business", "kerja", "meeting", "rapat"},
	"sendiri":     {"sendiri", "solo", "single", "personal"},
	"menginap":    {"menginap", "inap", "tidur", "stay", "istirahat"},
	"liburan":     {"liburan", "vacation", "holiday", "healing", "staycation", "piknik"},
	"pemandangan": {"pemandangan", "view", "panorama", "scenery", "lanskap"},
	"pelayanan":   {"pelayanan", "service", "layanan", "fasilitas"},
}

var intentTemplates = map[string]map[string]int{
	"honeymoon":       {"Honeymoon Suite": 9, "Suite": 6, "Deluxe": 4},
	"bulan madu":      {"Honeymoon Suite": 9, "Suite": 6, "Deluxe": 4},
	"romantis":        {"Honeymoon Suite": 8, "Suite": 6, "Deluxe": 4},
	"pasangan":        {"Honeymoon Suite": 7, "Deluxe": 5},
	"couple":          {"Honeymoon Suite": 7, "Deluxe": 5},
	"keluarga":        {"Family Room": 9, "Standard": 4, "Deluxe": 3},
	"family":          {"Family Room": 9, "Standard": 4, "Deluxe": 3},
	"anak":            {"Family Room": 8, "Standard": 4},
	"bisnis":          {"Business Suite": 9, "Deluxe": 4, "Standard": 2},
	"business":        {"Business Suite": 9, "Deluxe": 4, "Standard": 2},
	"mewah":           {"Pool Villa": 8, "Suite": 7, "Honeymoon Suite": 5},
	"luxury":          {"Pool Villa": 9, "Suite": 7, "Honeymoon Suite": 6},
	"eksklusif":       {"Pool Villa": 8, "Honeymoon Suite": 6},
	"vip":             {"Pool Villa": 8, "Suite": 5},
	"villa":           {"Pool Villa": 10},
	"murah":           {"Economy": 8, "Standard": 5},
	"budget":          {"Economy": 8, "Standard": 5},
	"ekonomis":        {"Economy": 9, "Standard": 4},
	"hemat":           {"Economy": 8, "Standard": 5},
	"backpacker":      {"Economy": 10},
	"kolam":           {"Pool Villa": 7, "Suite": 6},
	"renang":          {"Pool Villa": 7, "Suite": 6},
	"jacuzzi":         {"Honeymoon Suite": 7},
	"pemandangan":     {"Pool Villa": 6, "Honeymoon Suite": 5, "Deluxe": 3},
	"sunset":          {"Honeymoon Suite": 7, "Pool Villa": 6},
}

func keywordMatch(roomTypes []clients.RoomTypeDTO, rooms []clients.RoomDTO, message string) string {
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

	roomMap := make(map[uint][]string)
	for _, r := range rooms {
		if r.Availability {
			roomMap[r.RoomTypeID] = append(roomMap[r.RoomTypeID], r.RoomNumber)
		}
	}

	return fmt.Sprintf("Berdasarkan preferensi Anda, kami merekomendasikan tipe **%s**%s dengan harga Rp%.0f/malam. %s",
		best.roomType.Name, formatRoomNumbers(roomMap[best.roomType.ID]), best.roomType.Price, best.roomType.Description)
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
