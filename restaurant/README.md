## 동작 기능

### 새로운 가게 등록
1. 새로운 가게를 등록합니다.
2. 기본적인 등록에 필요한 정보를 입력합니다.
    - 가게명, 가게 주소, 사업자 번호, 사업자 대표자명을 입력합니다.
    - 사업자 번화 사업자 대표자명, 가게명 등이 올바른지 확인합니다.
4. 위 절차가 끝나면 성공적으로 끝났습니다.

### 가게 관리자 로그인 
1. 가게의 사업자 번호를 입력합니다.
2. 회원가입 시 입력한 비밀번호를 입력합니다.
    - 성공시 JWT 발행
    - 실패시 오류 메세지으로 응답

## 기능
- [X] 새로운 가게 가입 기능
    - 가게명, 가게 주소, 가게 전화번호, 사업자 번호, 사업자 대표자명
- [X] 가게 관리용 페이지 로그인 - 2022.07.26
    - [X] JWT 발행 및 쿠키 설정 완료
- [ ] 사업자 번호가 올바른 확인
    - [X] 사업자 등록 상태조회 기능  - 2022.07.17
        - [국세청 사업자등록정보 진위확인 및 상태조회 서비스](https://www.data.go.kr/tcs/dss/selectApiDataDetailView.do?publicDataPk=15081808) 사용
    - [ ] 사업자 등록 정보 진위 확인 기능
        - [국세청 사업자등록정보 진위확인 및 상태조회 서비스](https://www.data.go.kr/tcs/dss/selectApiDataDetailView.do?publicDataPk=15081808) 사용

## Docs
### `/api/restaurant/new` - 회원가입
#### Request
```json
{
    "restaurant_name": "팀그릿",
    "restaurant_location": "경기 성남시 수정구 대왕판교로 815",
    "restaurant_owner_name": "김기령",
    "restaurant_business_number": "1858800876",
    "restaurant_password": "parkhs0625!"
}
```

#### Response
```json
{
    "message": "성공적으로 가게를 등록했어요. 사업자 번호와 패스워드를 통해서 실시간으로 현황을 파악할 수 있어요.",
    "status": 200,
    "time": "2022-07-25T15:03:12.618421+09:00"
}
```

### `/api/restaurant/login` - 로그인
#### Request
```json
{
    "restaurant_business_number": "1858800876",
    "restaurant_password": "parkhs0625!"
}
```

#### Response
```
Redirct to Logined Page
```