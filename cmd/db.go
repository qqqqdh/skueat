package main

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Restaurant struct {
	gorm.Model
	Title       string  `json:"title"`
	Addr        string  `json:"addr"`
	Food        string  `json:"food"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	URL         string  `json:"url"`
	AvgRating   float64 `json:"avg_rating" gorm:"default:0"`   // 평균 별점
	RatingCount int     `json:"rating_count" gorm:"default:0"` // 참여 인원
}

// 별점 기록 테이블
type Rating struct {
	gorm.Model
	RestaurantID uint   `json:"restaurant_id"`
	UserID       string `json:"user_id"` // 카카오 고유 ID
	Score        int    `json:"score"`
}

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("restaurants.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("DB 연결 실패:", err)
	}
	// 두 테이블 모두 마이그레이션
	DB.AutoMigrate(&Restaurant{}, &Rating{})

	var count int64
	DB.Model(&Restaurant{}).Count(&count)
	if count == 0 {
		seedData()
	}
}
func seedData() {
	samples := []Restaurant{
		{Title: "부산가야밀면 안양본점", Addr: "경기도 안양시 만안구 문예로36번길 15", Food: "국수", X: 126.932263875909, Y: 37.3848854642594, URL: "https://place.map.kakao.com/13092162"},
		{Title: "지호한방삼계탕 만안구청점", Addr: "경기 안양시 만안구 안양로 115", Food: "닭요리, 고기", X: 126.932522205927, Y: 37.3849110208067, URL: "https://place.map.kakao.com/17978026"},
		{Title: "미소푸드", Addr: "경기 안양시 만안구 안양로 119", Food: "한식뷔페", X: 126.932155188339, Y: 37.385330147164, URL: "https://place.map.kakao.com/888574466"},
		{Title: "중찬미식", Addr: "경기 안양시 만안구 냉천로 4", Food: "중식", X: 126.929747804889, Y: 37.3822933729029, URL: "https://place.map.kakao.com/1001883450"},
		{Title: "와우리순대국", Addr: "경기 안양시 만안구 성결대학로 34-1", Food: "순대", X: 126.930395802091, Y: 37.3824244045313, URL: "https://place.map.kakao.com/16462275"},
		{Title: "엄마돈국수", Addr: "경기 안양시 만안구 냉천로 11", Food: "국수", X: 126.929052462381, Y: 37.3826812991012, URL: "https://place.map.kakao.com/118520162"},
		{Title: "토박이감자탕", Addr: "경기 안양시 만안구 안양로111번길 29", Food: "감자탕, 탕", X: 126.931543619665, Y: 37.3839742047466, URL: "https://place.map.kakao.com/16273180"},
		{Title: "참맛남원추어탕", Addr: "경기 안양시 만안구 문예로36번길 11", Food: "추어탕, 탕", X: 126.932179039671, Y: 37.3850507534564, URL: "https://place.map.kakao.com/12557676"},
		{Title: "오늘은수제돈까스", Addr: "경기 안양시 만안구 안양로111번길 33", Food: "돈까스, 우동", X: 126.931419468162, Y: 37.3839136741253, URL: "https://place.map.kakao.com/1229249213"},
		{Title: "부엉이샤브샤브스키야키", Addr: "경기 안양시 만안구 문예로36번길 15", Food: "샤브샤브", X: 126.932263875909, Y: 37.3848854642594, URL: "https://place.map.kakao.com/2105635163"},
		{Title: "돈컵치컵", Addr: "경기 안양시 만안구 성결대학로 47", Food: "분식", X: 126.929342303114, Y: 37.3816061005363, URL: "https://place.map.kakao.com/1864598783"},
		{Title: "소림마라 안양만안점", Addr: "경기 안양시 만안구 성결대학로 38", Food: "중식", X: 126.930151143008, Y: 37.3822779339963, URL: "https://place.map.kakao.com/1179574561"},
		{Title: "소선", Addr: "경기 안양시 만안구 안양로112번길 13", Food: "중식", X: 126.933877079736, Y: 37.3851813723658, URL: "https://place.map.kakao.com/2065526380"},
		{Title: "아리산", Addr: "경기 안양시 만안구 문예로 59", Food: "중식", X: 126.933853870283, Y: 37.3867938308462, URL: "https://place.map.kakao.com/98003334"},
		{Title: "몽샹", Addr: "경기 안양시 만안구 만안로 35", Food: "중식", X: 126.933874259452, Y: 37.386335581797, URL: "https://place.map.kakao.com/1380419551"},
		{Title: "풍미양꼬치", Addr: "경기 안양시 만안구 만안로 35", Food: "양꼬치", X: 126.933874259452, Y: 37.386335581797, URL: "https://place.map.kakao.com/2092991987"},
		{Title: "소문난김밥처럼", Addr: "경기 안양시 만안구 성결대학로 20", Food: "분식", X: 126.931872646792, Y: 37.383108871331, URL: "https://place.map.kakao.com/1726415782"},
		{Title: "신전떡볶이 안양성결대점", Addr: "경기 안양시 만안구 성결대학로 36", Food: "떡볶이", X: 126.930292097675, Y: 37.3823569468604, URL: "https://place.map.kakao.com/20493874"},
		{Title: "소림김밥 안양만안점", Addr: "경기 안양시 만안구 성결대학로 38", Food: "분식", X: 126.930151143008, Y: 37.3822779339963, URL: "https://place.map.kakao.com/7207904"},
		{Title: "남촌김밥 본점", Addr: "경기 안양시 만안구 안양로 110", Food: "분식", X: 126.933807735556, Y: 37.384811012897, URL: "https://place.map.kakao.com/832238107"},
		{Title: "남촌김밥 별관", Addr: "경기 안양시 만안구 안양로 102", Food: "분식", X: 126.934115411124, Y: 37.3841903794355, URL: "https://place.map.kakao.com/1297051684"},
		{Title: "할머니가래떡볶이 안양점", Addr: "경기 안양시 만안구 안양로 96", Food: "떡볶이", X: 126.9343209889, Y: 37.3839787532217, URL: "https://place.map.kakao.com/1846388455"},
		{Title: "일대김밥", Addr: "경기 안양시 만안구 문예로52번길 18", Food: "분식", X: 126.933622449079, Y: 37.3855690302257, URL: "https://place.map.kakao.com/932346628"},
		{Title: "우리분식", Addr: "경기 안양시 만안구 안양로 112", Food: "분식", X: 126.933373536382, Y: 37.3849789901676, URL: "https://place.map.kakao.com/16946757"},
		{Title: "맘스", Addr: "경기 안양시 만안구 만안로 11", Food: "분식", X: 126.934877389049, Y: 37.3844231754901, URL: "https://place.map.kakao.com/1063347983"},
		{Title: "호치킨 안양성결대점", Addr: "경기 안양시 만안구 성결대학로 28", Food: "치킨", X: 126.931089001347, Y: 37.3827644063939, URL: "https://place.map.kakao.com/2126044626"},
		{Title: "45정닭도리탕 본점", Addr: "경기 안양시 만안구 안양로111번길 35", Food: "닭요리", X: 126.931264031048, Y: 37.3838623156351, URL: "https://place.map.kakao.com/1521907226"},
		{Title: "청국닭", Addr: "경기 안양시 만안구 냉천로 14", Food: "닭요리", X: 126.929123499945, Y: 37.3830250822937, URL: "https://place.map.kakao.com/224678890"},
		{Title: "반주", Addr: "경기 안양시 만안구 안양로111번길 37", Food: "술집", X: 126.931151294564, Y: 37.3837906186715, URL: "https://place.map.kakao.com/616190514"},
		{Title: "작은울타리", Addr: "경기 안양시 만안구 문예로36번길 11", Food: "호프, 요리주점", X: 126.932179039671, Y: 37.3850507534564, URL: "https://place.map.kakao.com/24612700"},
		{Title: "별밤지기", Addr: "경기 안양시 만안구 냉천로 11", Food: "술집", X: 126.929052462381, Y: 37.3826812991012, URL: "https://place.map.kakao.com/15498052"},
		{Title: "세븐마일 비어앤굿즈", Addr: "경기 안양시 만안구 만안로 11", Food: "술집", X: 126.934877389049, Y: 37.3844231754901, URL: "https://place.map.kakao.com/1392808843"},
		{Title: "바지", Addr: "경기 안양시 만안구 안양로 96", Food: "칵테일바", X: 126.9343209889, Y: 37.3839787532217, URL: "https://place.map.kakao.com/556107500"},
		{Title: "이모네", Addr: "경기 안양시 만안구 만안로 21", Food: "호프, 요리주점", X: 126.93458324149, Y: 37.3853352094007, URL: "https://place.map.kakao.com/1383867494"},
		{Title: "논산훈련소포차", Addr: "경기 안양시 만안구 문예로52번길 14", Food: "실내포장마차, 호프", X: 126.933505809097, Y: 37.3856925851157, URL: "https://place.map.kakao.com/18251761"},
		{Title: "동막골", Addr: "경기 안양시 만안구 안양로112번길 13", Food: "호프, 요리주점", X: 126.933877079736, Y: 37.3851813723658, URL: "https://place.map.kakao.com/18824487"},
		{Title: "비어캐빈 명학점", Addr: "경기 안양시 만안구 문예로52번길 19", Food: "호프, 요리주점", X: 126.934059264774, Y: 37.3855110684841, URL: "https://place.map.kakao.com/16169962"},
		{Title: "명학맥주커피", Addr: "경기 안양시 만안구 안양로112번길 13", Food: "호프, 요리주점", X: 126.933877079736, Y: 37.3851813723658, URL: "https://place.map.kakao.com/202656010"},
		{Title: "맥주톡", Addr: "경기 안양시 만안구 만안로 19", Food: "술집, 호프", X: 126.934642908972, Y: 37.3851507127771, URL: "https://place.map.kakao.com/979421874"},
		{Title: "밤이술이", Addr: "경기 안양시 만안구 만안로 11", Food: "호프, 요리주점", X: 126.934877389049, Y: 37.3844231754901, URL: "https://place.map.kakao.com/1516437783"},
		{Title: "7MILE", Addr: "경기 안양시 만안구 만안로 11", Food: "호프, 요리주점", X: 126.934877389049, Y: 37.3844231754901, URL: "https://place.map.kakao.com/175222422"},
		{Title: "복돼지숯불갈비", Addr: "경기 안양시 만안구 냉천로 12", Food: "갈비, 육류", X: 126.929347655111, Y: 37.3827614867437, URL: "https://place.map.kakao.com/1633437082"},
		{Title: "고기마을", Addr: "경기 안양시 만안구 안양로111번길 10", Food: "육류, 고기", X: 126.932379841442, Y: 37.3847593873264, URL: "https://place.map.kakao.com/9683431"},
		{Title: "갈빗 안양점", Addr: "경기 안양시 만안구 문예로18번길 31", Food: "육류, 고기", X: 126.931201820704, Y: 37.3834852010292, URL: "https://place.map.kakao.com/281588556"},
		{Title: "연신내 생제육볶음전문점", Addr: "경기 안양시 만안구 냉천로 6-1", Food: "육류, 고기", X: 126.929587763284, Y: 37.3824616790948, URL: "https://place.map.kakao.com/1372621907"},
		{Title: "베스트생갈비찜 찜닭", Addr: "경기 안양시 만안구 냉천로 12-1", Food: "갈비, 육류", X: 126.929272940761, Y: 37.3828452373025, URL: "https://place.map.kakao.com/384164778"},
		{Title: "더두툼삼겹식당", Addr: "경기 안양시 만안구 냉천로 6-1", Food: "삼겹살, 육류", X: 126.929587763284, Y: 37.3824616790948, URL: "http://place.map.kakao.com/1175395613"},
		{Title: "마구아 만안구청점", Addr: "경기 안양시 만안구 문예로 35", Food: "육류, 고기", X: 126.931582288348, Y: 37.3857814100562, URL: "https://place.map.kakao.com/733705500"},
		{Title: "불꽃 안양본점", Addr: "경기 안양시 만안구 문예로52번길 15", Food: "육류, 고기", X: 126.933795288338, Y: 37.3857400512593, URL: "https://place.map.kakao.com/60886083"},
		{Title: "푸짐한마을", Addr: "경기 안양시 만안구 안양로111번길 25", Food: "찌개, 전골", X: 126.93172600054, Y: 37.3840734227839, URL: "https://place.map.kakao.com/18521540"},
		{Title: "배스킨라빈스 안양만안구청점", Addr: "경기 안양시 만안구 안양로 119", Food: "아이스크림", X: 126.932155188339, Y: 37.385330147164, URL: "https://place.map.kakao.com/1157830979"},
		{Title: "맘스터치 성결대점", Addr: "경기 안양시 만안구 성결대학로 38", Food: "패스트푸드", X: 126.930151143008, Y: 37.3822779339963, URL: "https://place.map.kakao.com/24768588"},
		{Title: "떡궁", Addr: "경기 안양시 만안구 성결대학로 22", Food: "떡, 한과", X: 126.931519503345, Y: 37.3830528038902, URL: "https://place.map.kakao.com/16918055"},
		{Title: "장터보쌈", Addr: "경기 안양시 만안구 문예로18번길 25", Food: "족발, 보쌈", X: 126.931006837114, Y: 37.3837103433724, URL: "https://place.map.kakao.com/169986040"},
		{Title: "맵당 안양점", Addr: "경기 안양시 만안구 냉천로 2", Food: "갈비, 육류", X: 126.929919541365, Y: 37.3821758911713, URL: "https://place.map.kakao.com/926912566"},
		{Title: "부산회집", Addr: "경기 안양시 만안구 안양로111번길 21", Food: "회, 생선", X: 126.931893907867, Y: 37.3841956984097, URL: "https://place.map.kakao.com/19305451"},
		{Title: "또와유명태간장조림 안양점", Addr: "경기 안양시 만안구 안양로111번길 33", Food: "해물, 생선", X: 126.931419468162, Y: 37.3839136741253, URL: "https://place.map.kakao.com/2061821812"},
		{Title: "육꼬", Addr: "경기 안양시 만안구 만안로 11", Food: "술집, 고기", X: 126.931419468162, Y: 37.3839136741253, URL: "https://place.map.kakao.com/1509168042"},
		{Title: "메가MGC커피 안양성결대점", Addr: "경기 안양시 만안구 성결대학로 34", Food: "카페", X: 126.930431, Y: 37.382421, URL: "https://place.map.kakao.com/1393693253"},
		{Title: "컴포즈커피 안양성결대점", Addr: "경기 안양시 만안구 성결대학로 28", Food: "카페", X: 126.931089, Y: 37.382764, URL: "https://place.map.kakao.com/1381368940"},
		{Title: "에이바우트커피 성결대점", Addr: "경기 안양시 만안구 성결대학로 38", Food: "카페", X: 126.930151, Y: 37.382278, URL: "https://place.map.kakao.com/1330962388"},
		{Title: "더카페", Addr: "경기 안양시 만안구 성결대학로 48-1 1층", Food: "카페", X: 126.92905100691, Y: 37.3818291998909, URL: "https://place.map.kakao.com/624783644"},
		{Title: "하이포커스", Addr: "경기 안양시 만안구 성결대학로 31", Food: "카페", X: 126.928976, Y: 37.381696, URL: "https://place.map.kakao.com/624783644"},
		{Title: "동대문 엽기떡볶이", Addr: "경기 안양시 만안구 성결대학로 28 1층", Food: "분식", X: 126.931089001347, Y: 37.3827644063939, URL: "https://place.map.kakao.com/16170476"},
		{Title: "신전떡볶이", Addr: "경기 안양시 만안구 성결대학로 36 1층", Food: "분식", X: 126.930292097675, Y: 37.3823569468604, URL: "https://place.map.kakao.com/16170476"},
		{Title: "힐링돈가스", Addr: "경기 안양시 만안구 성결대학로 47 1층", Food: "고기", X: 126.929342303114, Y: 37.3816061005363, URL: "https://place.map.kakao.com/279095344"},
		{Title: "가마치통닭", Addr: "경기 안양시 만안구 성결대학로 30 1층 101호", Food: "치킨", X: 126.930879977321, Y: 37.3826648113753, URL: "https://place.map.kakao.com/1051546409"},
	}
	DB.Create(&samples)
}
