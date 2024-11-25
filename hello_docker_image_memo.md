# hello を実行する Docker イメージを作成するメモ

1. hello のバイナリを作成する
    ```sh
    (cd hello && GOOS=linux GOARCH=arm64 go build -o ../output/hello)
    ```
1. Docker イメージをビルドする
    ```sh
    docker buildx build -t go-example .
    ```
