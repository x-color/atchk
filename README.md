# Atcoder Code Checker

Atcoderの問題に対する解答コードのテストツール。

問題についている入出力例を自動で読み込み、作成した解答コードが正しいかテストしてくれる。
現在、ABC,AGC問題に対応している。

## Usage

```
# Go言語のコードをテスト
# main.go(解答)をABC001 A問題の入出力例でテストする。
$ atchk abc001a go run main.go
Check ABC001A

- Execute Sample 1 ... SUCCESS
- Execute Sample 2 ... SUCCESS
- Execute Sample 3 ... SUCCESS

Your code is OK!!!

# Pythonのコードをテスト
# main.py(解答)をABC110 C問題の入出力例でテストする。
$ atchk abc110c python main.py
Check ABC110C

- Execute Sample 1 ... SUCCESS
- Execute Sample 2 ... SUCCESS
- Execute Sample 3 ... FAILURE
    expected:Yes
    actual  :No

Your code is NG...
```

## Install

必要となるツール

- Git
- Go1.11 以上

インストール手順

```
$ git clone https://github.com/x-color/atchk.git
$ cd atchk
$ GO111MODULE=on go install
```

## License

MIT License