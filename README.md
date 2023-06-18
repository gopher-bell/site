## Site

This project is a web application called Jong Box.  
It's a simple CRUD board with user authentications.

API Lists

| URL                 | Method | Description    |
|---------------------|--------|:---------------|
| "/"                 | GET    | 글 목록           |
| "/snippet/view/:id" | GET    | 글 보기           |
| "/user/signup"      | GET    | 회원가입 폼 호출      |
| "/user/signup"      | POST   | 회원가입 API 호출    |
| "/user/login"       | GET    | 로그인 폼 호출       |
| "/user/login"       | POST   | 로그인 API 호출     |
| "/sessions"         | GET    | 현재 세선 확인       |
| "/snippet/create"   | GET    | 글 쓰기 폼 호출      |
| "/snippet/create"   | POST   | 글 쓰기 API 호출    |
| "/user/logout"      | POST   | 로그아웃 API 호출    |