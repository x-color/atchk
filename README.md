# Atcoder Code Checker

Atcoderの問題に対する解答コードのテストツール。

問題についている入出力例を自動で読み込み、作成した解答コードが正しいかテストしてくれる。
また、解答コードの提出機能も付いている。
現在、ABC,AGC問題に対応している。

## できること

- 入出力例を用いた解答コードの自動テスト
- 解答コードの提出

## Usage

```
$ atchk --help
Usage:
  atchk [flags]
  atchk [command]

Available Commands:
  config      edit & show config
  help        Help about any command
  init        initialize config of atchk
  login       login atcoder
  logout      logout atcoder
  submit      submit answer code
  test        test the answer code

Flags:
  -h, --help   help for atchk

Use "atchk [command] --help" for more information about a command.
```

### 使い方の流れ

```
# 設定ファイルを用意する1（初回のみ実行）
$ atchk init

# Atcoderにログイン
$ atchk login

# Go言語のコードをテスト
# main.go(解答)をABC001 A問題の入出力例でテストする。
$ atchk test abc001a "go run main.go"
- Execute Sample 1 ... SUCCESS
- Execute Sample 2 ... SUCCESS
- Execute Sample 3 ... SUCCESS

# Pythonのコードをテスト
# main.py(解答)をABC110 C問題の入出力例でテストする。
$ atchk abc110c "python main.py"
- Execute Sample 1 ... SUCCESS
- Execute Sample 2 ... SUCCESS
- Execute Sample 3 ... FAILURE

# main.go(解答)をABC001 A問題の解答として提出する。
$ atchk submit abc001a main.go
```

### 設定ファイル

以下のファイルを編集することで解答コードの言語やテストコードの実行時のコマンドを指定することができる。
`config` 内を編集可能。`cache` 内は編集しないでください。

`~/.atchk/config.json`

```json
{
  "cache": {},
  "config": {
    "cmd": "",
    "lang_id": ""
  }
}
```

- `config.cmd`: `atchk test` 実行時のコマンドを指定可能。

```bash
# 例として、"cmd": "go run" とした場合
# 以下でテスト可能となる。
$ atchk test abc001a main.go # atchk test abc001a "go run main.go" を実行していることになる
```

- `config.lang_id`: `atchk submit` 実行時の言語を指定。(**Required**)

指定方法は以下の通り。

```bash
#  Go言語を指定したい場合。
$ atchk config --lang-list
  :
- Fortran (gfortran v4.8.4) : 3012
- Go (1.6) : 3013 # Golang's ID is 3013 
- Haskell (GHC 7.10.3) : 3014
  :

$ vim ~/.atchk/config.json
# "lang_id": "3013"
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