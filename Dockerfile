####################### Development stage #######################
FROM golang:1.20.5-alpine3.18 AS development

# 作業ディレクトリの定義
WORKDIR /twitter_clone

# go.mod と go.sum を twitter_clone ディレクトリにコピー
COPY go.mod go.sum ./

# 指定されたモジュールをダウンロード
RUN go mod download

# ルートディレクトリの中身を twitter_clone フォルダにコピー
COPY . .

# ビルドするファイルを指定：main.go
RUN go build -o main /twitter_clone/main.go

# air インストール
RUN go get -u github.com/cosmtrek/air && go build -o /go/bin/air github.com/cosmtrek/air

# .air.toml ファイルをコピー
COPY .air.toml ./

CMD ["air", "-c", ".air.toml"]

####################### Production stage #######################
FROM alpine:3.18 AS production

# 作業ディレクトリの定義
WORKDIR /twitter_clone

# Development stage からビルドされた main だけを Production stage にコピー
COPY --from=development /twitter_clone/main .

CMD ["./main"]
