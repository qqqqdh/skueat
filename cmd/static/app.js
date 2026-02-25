let map, currentMarker = null, currentOverlay = null;

// 안양 8동 경계 좌표
const boundaryCoords = [[126.936009, 37.38326], [126.936242, 37.382663], [126.936444, 37.382211], [126.936508, 37.382068], [126.936646, 37.381803], [126.936743, 37.381642], [126.936765, 37.381605], [126.93691, 37.381366], [126.937138, 37.381042], [126.937483, 37.380612], [126.937944, 37.380137], [126.938376, 37.379733], [126.938679, 37.379465], [126.938863, 37.379302], [126.939684, 37.378558], [126.940026, 37.37823], [126.940195, 37.378053], [126.940214, 37.378032], [126.940296, 37.377945], [126.940511, 37.377704], [126.940739, 37.377422], [126.940785, 37.37736], [126.940877, 37.377235], [126.940956, 37.377128], [126.941185, 37.376762], [126.940555, 37.376242], [126.940511, 37.376205], [126.940318, 37.376048], [126.940076, 37.375851], [126.939165, 37.376011], [126.937373, 37.376555], [126.934893, 37.376996], [126.932765, 37.37747], [126.930979, 37.377528], [126.926947, 37.377638], [126.924137, 37.377312], [126.923923, 37.377256], [126.921417, 37.376518], [126.919899, 37.376564], [126.919178, 37.376933], [126.91921, 37.377178], [126.919285, 37.377544], [126.919353, 37.377714], [126.919431, 37.377818], [126.920582, 37.379335], [126.921934, 37.380899], [126.924169, 37.383026], [126.925436, 37.384031], [126.925688, 37.384199], [126.925741, 37.384224], [126.92579, 37.384248], [126.932237, 37.385795], [126.93377, 37.383967], [126.933825, 37.383896], [126.936009, 37.383261]];

function initMap() {
    const container = document.getElementById('map');
    const options = { center: new kakao.maps.LatLng(37.382, 126.931), level: 3 };
    map = new kakao.maps.Map(container, options);

    const path = boundaryCoords.map(c => new kakao.maps.LatLng(c[1], c[0]));
    new kakao.maps.Polygon({
        path: path, strokeWeight: 3, strokeColor: '#FF0000', strokeOpacity: 0.6,
        fillColor: '#FF0000', fillOpacity: 0.05
    }).setMap(map);
}

async function fetchData(category = 'all', search = '') {
    const url = `/api/restaurants?category=${category}&search=${search}`;
    try {
        const response = await fetch(url);
        const data = await response.json();
        renderList(data);
    } catch (e) { console.error("로드 실패", e); }
}

function renderList(data) {
    const container = document.getElementById('res-list');
    document.getElementById('list-count').textContent = `주변 맛집 ${data.length}곳`;
    container.innerHTML = '';

    data.forEach(item => {
        const card = document.createElement('div');
        card.className = 'res-card';
        
        const avg = item.avg_rating || 0;
        const count = item.rating_count || 0;

        card.innerHTML = `
            <a class="res-title" href="${item.url}" target="_blank" onclick="event.stopPropagation()">${item.title}</a>
            <div class="rating-box">
                <span class="star-display filled">★</span>
                <span class="rating-score">${avg.toFixed(1)}</span>
                <span class="rating-count">(${count})</span>
            </div>
            <div class="res-addr">${item.addr}</div>
            <div class="rate-btn-group">
                <div class="rate-dot" onclick="submitRate(${item.ID}, 1, event)">1</div>
                <div class="rate-dot" onclick="submitRate(${item.ID}, 2, event)">2</div>
                <div class="rate-dot" onclick="submitRate(${item.ID}, 3, event)">3</div>
                <div class="rate-dot" onclick="submitRate(${item.ID}, 4, event)">4</div>
                <div class="rate-dot" onclick="submitRate(${item.ID}, 5, event)">5</div>
            </div>
        `;
        card.onclick = () => focusOn(item, card);
        container.appendChild(card);
    });
}

async function submitRate(resId, score, event) {
    event.stopPropagation();
    // 전역 변수 IS_LOGGED_IN 사용 (index.html에서 넘어옴)
    if (typeof IS_LOGGED_IN === 'undefined' || !IS_LOGGED_IN) {
        alert("별점을 남기려면 카카오 로그인이 필요합니다! 🔒");
        return;
    }

    const formData = new URLSearchParams();
    formData.append('restaurant_id', resId);
    formData.append('score', score);

    try {
        const response = await fetch('/api/rate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: formData
        });
        const result = await response.json();
        if (response.ok) {
            alert(`평가 완료! 현재 평점: ${result.new_avg.toFixed(1)}`);
            applySearch(); 
        } else {
            alert(result.error);
        }
    } catch (e) { alert("평가 전송 중 오류가 발생했습니다."); }
}

function focusOn(item, cardElement) {
    document.querySelectorAll('.res-card').forEach(c => c.classList.remove('active'));
    if(cardElement) cardElement.classList.add('active');

    if (currentMarker) currentMarker.setMap(null);
    if (currentOverlay) currentOverlay.setMap(null);

    const pos = new kakao.maps.LatLng(item.y, item.x);
    currentMarker = new kakao.maps.Marker({ position: pos, map: map });
    currentOverlay = new kakao.maps.CustomOverlay({
        position: pos, content: `<div class="custom-overlay">${item.title}</div>`, yAnchor: 2.2
    });
    currentOverlay.setMap(map);
    map.panTo(pos);

    // 모바일(화면 너비 768px 이하)일 때 마커 클릭 시 바텀시트 자동으로 내리기
    if (window.innerWidth <= 768) {
        const sheet = document.getElementById('bottom-sheet');
        if (sheet) {
            sheet.style.transition = 'height 0.3s ease-out';
            sheet.style.height = '15%'; // vh가 아닌 % 사용
        }
    }
}

function applyFilter(category, btn) {
    document.querySelectorAll('.cat-btn').forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
    fetchData(category, document.getElementById('search-input').value);
}

function applySearch() {
    const activeBtn = document.querySelector('.cat-btn.active');
    const category = activeBtn ? activeBtn.dataset.category : 'all';
    fetchData(category === 'all' ? 'all' : category, document.getElementById('search-input').value);
}

async function pickRandom() {
    const res = await fetch('/api/restaurants/random');
    const pick = await res.json();
    const cards = document.querySelectorAll('.res-card');
    cards.forEach(card => {
        if(card.querySelector('.res-title').textContent === pick.title) {
            card.scrollIntoView({ behavior: 'smooth', block: 'center' });
            focusOn(pick, card);
        }
    });
}

// 💡 새롭게 고쳐진 무적의 바텀시트 로직
function initBottomSheet() {
    const sheet = document.getElementById('bottom-sheet');
    const handle = document.getElementById('sheet-handle');
    const mainContent = document.querySelector('.main-content'); // 지도+시트 영역 (헤더 제외)
    
    if (!sheet || !handle || !mainContent) return;

    let isDragging = false;
    let startY, startHeight;

    handle.addEventListener('touchstart', (e) => {
        isDragging = true;
        startY = e.touches[0].clientY; // pageY 대신 클라이언트Y 사용
        startHeight = sheet.getBoundingClientRect().height;
        sheet.style.transition = 'none'; 
    }, { passive: true });

    document.addEventListener('touchmove', (e) => {
        if (!isDragging) return;
        const deltaY = startY - e.touches[0].clientY;
        let newHeight = startHeight + deltaY;

        // 화면 전체가 아닌 main-content의 높이를 기준으로 제한
        const mainHeight = mainContent.getBoundingClientRect().height;
        const minHeight = mainHeight * 0.15; // 최소 15%
        const maxHeight = mainHeight * 0.95; // 최대 95%
        
        // 범위를 벗어나면 최소/최대값으로 강제 고정하여 이벤트 끊김 방지
        if (newHeight < minHeight) newHeight = minHeight;
        if (newHeight > maxHeight) newHeight = maxHeight;
        
        sheet.style.height = `${newHeight}px`;
    }, { passive: true });

    document.addEventListener('touchend', () => {
        if (!isDragging) return;
        isDragging = false;
        sheet.style.transition = 'height 0.3s ease-out'; 

        const currentHeight = sheet.getBoundingClientRect().height;
        const mainHeight = mainContent.getBoundingClientRect().height;

        // 터치를 뗐을 때 위치에 따라 스냅(자동 위치 조정) 작동 (% 기준)
        if (currentHeight > mainHeight * 0.6) {
            sheet.style.height = '95%'; // 최대로 올렸을 때
        } else if (currentHeight < mainHeight * 0.3) {
            sheet.style.height = '15%'; // 최소화
        } else {
            sheet.style.height = '45%'; // 중간 지점
        }
    });
}

// 초기화 및 이벤트 리스너 등록
document.addEventListener('DOMContentLoaded', () => {
    initMap();
    fetchData();
    initBottomSheet();

    // 검색어 입력
    document.getElementById('search-input').addEventListener('keyup', (e) => {
        if (e.key === 'Enter') applySearch();
    });

    // 랜덤 버튼
    document.getElementById('btn-random').addEventListener('click', pickRandom);

    // 카테고리 필터
    document.getElementById('category-nav').addEventListener('click', (e) => {
        if (e.target.classList.contains('cat-btn')) {
            applyFilter(e.target.dataset.category, e.target);
        }
    });
});