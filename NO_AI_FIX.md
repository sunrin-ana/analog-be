# NO AI FIX란?

## Details 
NO AI FIX는 해당 레포지토리가 생성형 인공지능 모델에 의해 작성되었을때 사람이 직접 생성형 인공지능이 만든 오류(논리적 오류 포함)하여 제거한 FIX를 의미합니다.

## Roadmap
- [ ] (0001) 인증 구현 flow를 재작성

# Fix

## (0001) Auth flow

### Current Problem
기존 `/api/auth/[login/signup]/[init/callback]`은 너무 복잡하며 극단의 비효율을 초례함.

### Solution
`/api/auth/callback`으로 통합하고 login 및 signup은 서버가 판단하도록 변경.
프론트엔드는 유저 정보 요청을 제공할때에 회원가입이 되어있지 않는다면 `425 Too Early` 를 반환하도록 변경. (회원가입이 되어있다면 `200 OK` 반환)

### Affected

- [FE] An Account의 OAuth으로 리다이렉트 시켜야함. (state는 URI ENCODING된 로그인 이후 리다이렉트될 URI를 제공해야함.)
- [FE] `GET /api/me`를 이용하여 유저 정보를 확인해야함.
- [FE] `PUT /api/me`를 이용하여 유저 정보를 업데이트 해야함.
