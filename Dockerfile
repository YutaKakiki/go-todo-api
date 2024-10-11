# ビルドステージ
FROM golang:1.23.1-bullseye AS deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
# go.modファイルに記載されたすべてのモジュールをダウンロード
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags "-w -s" -o -app

#----------------------------------------------------------------------

# デプロイステージ (軽量なベースイメージ)
FROM debian:bullseye-slim AS deploy

# 必要なパッケージをインストール
RUN apt-get update

# ビルドステージにおける/app/appファイル（ビルドファイル）をルートディレクトリにコピー
COPY --from=deploy-builder /app/app .

CMD [ "./app" ]

#----------------------------------------------------------------------

# ローカル環境で利用するホットリロード環境


FROM golang:1.23.1 as dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

CMD ["air"]