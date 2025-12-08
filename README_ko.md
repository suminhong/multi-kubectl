# Multi-Kubectl (mk)

[English](README.md) | [한국어](README_ko.md)

![GitHub all releases](https://img.shields.io/github/downloads/suminhong/multi-kubectl/total)

`mk`는 여러 Kubernetes 컨텍스트에서 `kubectl` 명령어를 실행하고 결과를 취합하여 보여주는 CLI 도구입니다.

## 사전 요구사항 (Prerequisites)
- `kubectl`이 설치되어 있고 설정되어 있어야 합니다.
```bash
brew install kubectl
```

## 설치 (Installation)

### Homebrew
커스텀 탭을 사용하여 `mk`를 설치할 수 있습니다:

```bash
brew install suminhong/tap/mk
```

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

### 컨텍스트 그룹 (Context Groups)
`~/.kube/mk_config` 파일에 컨텍스트 그룹을 정의할 수 있습니다:

```yaml
- name: dev
  contexts: dev-cluster, alpha-cluster
- name: aws
  contexts: dev-eks, prod-eks
```

그룹으로 실행:

```bash
mk get pods -g dev
```

여러 그룹으로 실행:
```bash
mk get pods -g dev,aws
```

> [!NOTE]
> `--context`와 `-g` 옵션은 함께 사용할 수 없습니다. 두 옵션이 모두 지정된 경우 `-g` 옵션이 우선 적용됩니다.

![example](images/example.png)

## 배포 프로세스 (Release Process)

GitHub Actions와 Goreleaser를 사용하여 배포가 자동화되어 있습니다. `main` 브랜치에 푸시하면 새로운 릴리스가 트리거됩니다.
