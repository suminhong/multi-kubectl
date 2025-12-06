# Multi-Kubectl (mk)

[English](README.md) | [한국어](README_ko.md)

`mk`는 여러 Kubernetes 컨텍스트에서 동시에 `kubectl` 명령어를 실행할 수 있게 해주는 CLI 도구입니다.

## 사용법 (Usage)

```bash
mk [kubectl-args] [--context <filter>]
```

### 예시 (Examples)

모든 컨텍스트에서 `get pods` 실행:
```bash
mk get pods -n devops
```

"prod" 문자열을 포함하는 컨텍스트에서 실행:
```bash
mk get pods --context prod
```

"prod" 또는 "stage" 문자열을 포함하는 컨텍스트에서 실행:
```bash
mk get pods --context prod,stage
```

## 설치 (Installation)

### Homebrew
커스텀 탭을 사용하여 `mk`를 설치할 수 있습니다:

```bash
brew install suminhong/tap/mk
```

### 소스코드에서 빌드 (From Source)
```bash
git clone https://github.com/suminhong/multi-kubectl.git
cd multi-kubectl
make build
sudo make install
```

## 배포 프로세스 (Release Process)

GitHub Actions와 Goreleaser를 사용하여 배포가 자동화되어 있습니다. `main` 브랜치에 푸시하면 새로운 릴리스가 트리거됩니다.
