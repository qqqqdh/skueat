package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	InitDB() // db.go에 있는 함수 호출

	godotenv.Load() // .env 파일 로드
	r := gin.Default()

	// 템플릿 파일 로드 설정 (index.html 내의 {{.ApiKey}} 치환을 위해 필요)
	r.LoadHTMLGlob("index.html")

	// CORS 설정
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// 루트 접속 시 HTML 렌더링 (API 키 주입)
	r.GET("/", func(c *gin.Context) {
		apiKey := os.Getenv("KAKAO_API_KEY")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ApiKey": apiKey,
		})
	})

	// 맛집 리스트 API (카테고리 및 검색 필터링 추가)
	r.GET("/api/restaurants", func(c *gin.Context) {
		category := c.Query("category")
		search := c.Query("search")

		var list []Restaurant
		query := DB

		// 카테고리 필터링: 'all'이 아니면 Food 필드에서 검색
		if category != "" && category != "all" {
			query = query.Where("food LIKE ?", "%"+category+"%")
		}

		// 검색어 필터링: 제목(Title)이나 주소(Addr)에 포함된 경우
		if search != "" {
			query = query.Where("title LIKE ? OR addr LIKE ?", "%"+search+"%", "%"+search+"%")
		}

		query.Find(&list)
		c.JSON(http.StatusOK, list)
	})

	// 무작위 추천 API
	r.GET("/api/restaurants/random", func(c *gin.Context) {
		var pick Restaurant
		DB.Order("RANDOM()").First(&pick)
		c.JSON(http.StatusOK, pick)
	})

	r.Run(":8080")
}
