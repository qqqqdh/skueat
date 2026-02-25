package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// KakaoTokenResponse: 카카오 토큰 발급 응답 구조체
type KakaoTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// KakaoUserResponse: 카카오 사용자 정보 응답 구조체
type KakaoUserResponse struct {
	ID         int64 `json:"id"`
	Properties struct {
		Nickname string `json:"nickname"`
	} `json:"properties"`
	KakaoAccount struct {
		Profile struct {
			Nickname string `json:"nickname"`
		} `json:"profile"`
	} `json:"kakao_account"`
}

// ... existing code ...

func main() {
	// 1. 환경변수 초기화
	godotenv.Load()

	// 2. 로그 시스템 설정
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	f, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("로그 파일을 열 수 없습니다: %v", err)
		return
	}

	multiWriter := io.MultiWriter(f, os.Stdout)
	gin.DefaultWriter = multiWriter
	log.SetOutput(multiWriter)
	log.SetFlags(0)

	// 3. Gin 엔진 생성
	r := gin.New()
	r.Use(gin.Recovery())

	// 4. 커스텀 로그 포맷터 설정
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		level := "INFO "
		if param.StatusCode >= 400 && param.StatusCode < 500 {
			level = "WARN "
		} else if param.StatusCode >= 500 {
			level = "ERROR"
		}

		return fmt.Sprintf("%s [%d] %s %s (%s)\n",
			level,
			param.StatusCode,
			param.Method,
			param.Path,
			param.Latency,
		)
	}))

	// 5. DB 초기화
	InitDB()
	log.Println("INFO  app started")

	// 6. 세션 설정
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("index.html")

	// --- [라우터 설정] ---

	// 메인 페이지
	r.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		userName := session.Get("userName")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ApiKey":     os.Getenv("KAKAO_API_KEY"),
			"IsLoggedIn": userName != nil,
			"UserName":   userName,
			"AppDomain":  os.Getenv("APP_DOMAIN"),
		})
	})

	// 맛집 리스트 API
	r.GET("/api/restaurants", func(c *gin.Context) {
		category := c.Query("category")
		search := c.Query("search")

		var list []Restaurant
		query := DB.Model(&Restaurant{})

		if category != "" && category != "all" {
			query = query.Where("food LIKE ?", "%"+category+"%")
		}
		if search != "" {
			query = query.Where("title LIKE ? OR addr LIKE ?", "%"+search+"%", "%"+search+"%")
		}

		query.Find(&list)
		c.JSON(http.StatusOK, list)
	})

	// 무작위 추천 API
	r.GET("/api/restaurants/random", func(c *gin.Context) {
		var pick Restaurant
		if err := DB.Order("RANDOM()").First(&pick).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "데이터를 찾을 수 없습니다."})
			return
		}
		c.JSON(http.StatusOK, pick)
	})

	// 카카오 로그인 시작
	r.GET("/login/kakao", func(c *gin.Context) {
		clientID := os.Getenv("REST_API_KEY")
		appDomain := os.Getenv("APP_DOMAIN")
		redirectURI := appDomain + "/auth/kakao/callback"
		kakaoURL := fmt.Sprintf("https://kauth.kakao.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", clientID, url.QueryEscape(redirectURI))
		c.Redirect(http.StatusFound, kakaoURL)
	})

	// 카카오 콜백 처리
	r.GET("/auth/kakao/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.String(http.StatusBadRequest, "인가 코드가 없습니다.")
			return
		}

		appDomain := os.Getenv("APP_DOMAIN")
		tokenRes, err := getKakaoToken(code, appDomain)
		if err != nil {
			c.String(http.StatusInternalServerError, "토큰 발급 실패")
			return
		}

		userInfo, err := getKakaoUserInfo(tokenRes.AccessToken)
		if err != nil {
			c.String(http.StatusInternalServerError, "사용자 정보 조회 실패")
			return
		}

		session := sessions.Default(c)
		session.Set("userName", userInfo.Properties.Nickname)
		session.Save()

		c.Redirect(http.StatusFound, "/")
	})

	// 별점 평가 API
	r.POST("/api/rate", func(c *gin.Context) {
		session := sessions.Default(c)
		userName := session.Get("userName")
		if userName == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "로그인이 필요합니다."})
			return
		}

		resID, _ := strconv.Atoi(c.PostForm("restaurant_id"))
		score, _ := strconv.Atoi(c.PostForm("score"))

		var res Restaurant
		if err := DB.First(&res, resID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "식당을 찾을 수 없습니다."})
			return
		}

		rating := Rating{
			RestaurantID: uint(resID),
			UserID:       userName.(string),
			Score:        score,
		}
		DB.Create(&rating)

		newCount := res.RatingCount + 1
		newAvg := (res.AvgRating*float64(res.RatingCount) + float64(score)) / float64(newCount)
		DB.Model(&res).Updates(map[string]interface{}{
			"AvgRating":   newAvg,
			"RatingCount": newCount,
		})

		c.JSON(http.StatusOK, gin.H{"message": "평가가 완료되었습니다.", "new_avg": newAvg})
	})

	// 로그아웃
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/")
	})

	// 동적 포트 바인딩
	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")
	if appHost == "" {
		appHost = "0.0.0.0"
	}
	if appPort == "" {
		appPort = "8080"
	}

	addr := appHost + ":" + appPort
	log.Printf("INFO  server listening on %s\n", addr)
	r.Run(addr)
}

// --- [도움 함수] ---

func getKakaoToken(code string, appDomain string) (*KakaoTokenResponse, error) {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", os.Getenv("REST_API_KEY"))
	params.Add("redirect_uri", appDomain+"/auth/kakao/callback")
	params.Add("code", code)

	resp, err := http.PostForm("https://kauth.kakao.com/oauth/token", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenRes KakaoTokenResponse
	json.NewDecoder(resp.Body).Decode(&tokenRes)
	return &tokenRes, nil
}

func getKakaoUserInfo(token string) (*KakaoUserResponse, error) {
	req, _ := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userRes KakaoUserResponse
	json.NewDecoder(resp.Body).Decode(&userRes)
	return &userRes, nil
}
