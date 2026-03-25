# Jaeseong-916.github.io

Jekyll + [Minimal Mistakes](https://mmistakes.github.io/minimal-mistakes/) 테마 기반 개인 기술 블로그.

---


## 프로젝트 구조

```
.
├── _config.yml              # 사이트 설정 (제목, 저자, 플러그인 등)
├── _data/
│   └── navigation.yml       # 상단 네비게이션 메뉴
├── _pages/
│   ├── about.md             # About 페이지
│   └── categories/
│       ├── python.md        # /categories/python/
│       ├── go.md            # /categories/go/
│       ├── kubernetes.md    # /categories/kubernetes/
│       ├── openstack.md     # /categories/openstack/
│       └── system.md        # /categories/system/
├── _posts/                  # 블로그 글
├── assets/
│   └── css/main.scss        # Sass 엔트리포인트
├── index.html               # 홈 페이지
└── .github/
    └── workflows/
        └── pages.yml        # GitHub Actions 배포 워크플로
```

---

## 글 추가하는 방법

### 1. 파일 생성

`_posts/` 디렉토리에 아래 형식으로 파일을 생성합니다.

```
_posts/YYYY-MM-DD-제목.md
```

예시: `_posts/2026-03-24-python-basic.md`

### 2. Front Matter 작성

파일 최상단에 아래 형식으로 메타데이터를 작성합니다.

```yaml
---
title: "글 제목"
date: 2026-03-24
categories:
  - python        # python / go / kubernetes / openstack / system 중 하나
tags:
  - 태그1
  - 태그2
---

본문 내용...
```

### 3. 카테고리 목록

| 카테고리 | URL |
|----------|-----|
| `python` | /categories/python/ |
| `go` | /categories/go/ |
| `kubernetes` | /categories/kubernetes/ |
| `openstack` | /categories/openstack/ |
| `system` | /categories/system/ |

---

## 로컬 실행

```bash
export PATH="/opt/homebrew/opt/ruby/bin:$PATH"
export GEM_HOME="$HOME/.gem/ruby/4.0.0"
export PATH="$GEM_HOME/bin:$PATH"
bundle exec jekyll serve
```

브라우저에서 `http://localhost:4000` 접속.
