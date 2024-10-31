# SyncTimeStamp

SyncTimeStamp is a command-line tool for Windows that synchronizes the timestamps (creation, modification, and access times) of target files or directories with those of reference files or directories. It supports an optional time shift for timezone adjustments and can operate in test mode to preview changes without modifying any files.

[↓日本語](#日本語)

## Features

- **File and Directory Synchronization**: Synchronize timestamps between individual files or entire directories.
- **Time Shift Adjustment**: Apply a time shift in hours to the reference timestamps for timezone adjustments.
- **Test Mode**: Preview the changes without actually modifying the target files.
- **Interactive Mode**: If command-line arguments are omitted, the tool will prompt for inputs interactively.

## Download

Pre-built binaries for Windows (`SyncTimeStamp.exe`) are available on the [GitHub Releases](https://github.com/motoacs/sync-time-stamp/releases) page.

## Usage

```bash
SyncTimeStamp.exe [options]
```

### Options

- `-t <path>`: **Target file or directory path**. The file or directory whose timestamps you want to modify.
- `-r <path>`: **Reference file or directory path**. The file or directory whose timestamps will be used.
- `-shift <hours>`: **Time shift in hours** (from -24 to 24). Adjusts the reference timestamps by the specified number of hours before applying them to the target files.
- `-test`: **Test mode**. Runs the tool in test mode, displaying the changes without modifying any files.

### Interactive Mode

If you omit any of the required options, the tool will prompt you to enter them interactively.

## Examples

### Synchronize a Single File

Synchronize the timestamps of `target.mp4` with those of `reference.mp4`:

```bash
SyncTimeStamp.exe -t "C:\path\to\target.mp4" -r "C:\path\to\reference.mp4"
```

### Synchronize Directories

Synchronize the timestamps of all files in `target_folder` with matching files in `reference_folder`:

```bash
SyncTimeStamp.exe -t "C:\path\to\target_folder" -r "C:\path\to\reference_folder"
```

**Note**: Files are matched based on their names without extensions. A target file is matched with a reference file if the target filename contains the reference filename.

### Apply Time Shift

Apply a time shift of +9 hours (e.g., adjusting from UTC to JST) when synchronizing:

```bash
SyncTimeStamp.exe -t "C:\path\to\target_folder" -r "C:\path\to\reference_folder" -shift 9
```

### Run in Test Mode

Preview the changes without modifying any files:

```bash
SyncTimeStamp.exe -t "C:\path\to\target_folder" -r "C:\path\to\reference_folder" -test
```

## Build from Source

To build `SyncTimeStamp.exe` from the source code, ensure you have [Go](https://golang.org/dl/) installed.

### Build Command

```bash
go build -o "SyncTimeStamp.exe"
```

Alternatively, you can use the provided `build.cmd` script:

```bash
build.cmd
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## Disclaimer

Use this tool at your own risk. Always ensure you have backups of your files before modifying them.


---


以下は、インタラクティブモードでの使用方法を詳しく説明するように編集したREADME.mdです。

---

# 日本語

SyncTimeStampは、Windows用のコマンドラインツールで、ターゲットのファイルやディレクトリのタイムスタンプ（作成日時、更新日時、アクセス日時）を参照ファイルやディレクトリのタイムスタンプと同期します。タイムゾーン調整のためにオプションで時間のずれを適用でき、テストモードで変更のプレビューも可能です。

## 機能

- **ファイルおよびディレクトリの同期**: 個別のファイルまたはディレクトリ全体のタイムスタンプを同期。
- **時間シフト調整**: タイムゾーン調整用に、参照タイムスタンプに時間のずれ（±時間）を適用可能。
- **テストモード**: ターゲットファイルを実際に変更せずに、変更内容をプレビュー。
- **インタラクティブモード**: コマンドライン引数が指定されていない場合、インタラクティブに入力を促します。

## ダウンロード

Windows用の実行ファイル（`SyncTimeStamp.exe`）は[GitHub Releases](https://github.com/motoacs/sync-time-stamp/releases)ページからダウンロードできます。

## 使用方法

```bash
SyncTimeStamp.exe [options]
```

### オプション

- `-t <path>`: **ターゲットファイルまたはディレクトリのパス**。タイムスタンプを変更するファイルまたはディレクトリを指定。
- `-r <path>`: **参照ファイルまたはディレクトリのパス**。参照タイムスタンプとして使用するファイルまたはディレクトリを指定。
- `-shift <hours>`: **時間のシフト量（-24〜+24時間）**。参照タイムスタンプに適用する時間のずれを設定。
- `-test`: **テストモード**。ターゲットファイルを変更せずに、プレビューのみを表示します。

### インタラクティブモード

コマンドライン引数を指定せずに実行した場合、ツールはインタラクティブモードで起動します。必要な情報が順に求められますので、プロンプトに従って入力してください。

#### インタラクティブモードの手順

1. **ターゲットパスの入力**:
   ```
   ターゲットディレクトリまたはファイルのパスを入力 (-t):
   ```
   タイムスタンプを変更したいファイルまたはディレクトリのパスを入力します。

2. **参照パスの入力**:
   ```
   参照ディレクトリまたはファイルのパスを入力 (-r):
   ```
   参照とするタイムスタンプを持つファイルまたはディレクトリのパスを入力します。

3. **時間シフトの入力**（オプション）:
   ```
   タイムゾーン調整のための時間シフト (-24から24) [デフォルト 0]:
   ```
   タイムゾーン調整が必要な場合、シフトしたい時間（-24から24の整数）を入力します。入力を省略するとデフォルト値の`0`が適用されます。

4. **テストモードの確認**:
   ツールは最初にテストモードで変更内容をプレビューします。
   ```
   テストモードで実行中...
   ```
   変更内容が表示された後、実際に変更を適用するかを確認されます。
   ```
   実際の処理を実行しますか？ (y/n):
   ```
   `y`または`Y`を入力すると、変更が適用されます。それ以外を入力すると、処理がキャンセルされます。

### インタラクティブモードの例

#### 例1: ファイルの同期

```bash
SyncTimeStamp.exe
```

**出力例**:
```
ターゲットディレクトリまたはファイルのパスを入力 (-t): C:\path\to\target.mp4
参照ディレクトリまたはファイルのパスを入力 (-r): C:\path\to\reference.mp4
タイムゾーン調整のための時間シフト (-24から24) [デフォルト 0]: 0
テストモードで実行中...
ファイル: C:\path\to\target.mp4
  作成日時 (変更前): 2021-09-01 12:00:00
  作成日時 (変更後): 2021-09-01 10:00:00
  変更日時 (変更前): 2021-09-01 12:00:00
  変更日時 (変更後): 2021-09-01 10:00:00
  アクセス日時 (変更前): 2021-09-01 12:00:00
  アクセス日時 (変更後): 2021-09-01 10:00:00
実際の処理を実行しますか？ (y/n): y
処理が正常に完了しました。
```

#### 例2: ディレクトリの同期と時間シフト

```bash
SyncTimeStamp.exe
```

**出力例**:
```
ターゲットディレクトリまたはファイルのパスを入力 (-t): C:\path\to\target_folder
参照ディレクトリまたはファイルのパスを入力 (-r): C:\path\to\reference_folder
タイムゾーン調整のための時間シフト (-24から24) [デフォルト 0]: 9
テストモードで実行中...
ファイル: C:\path\to\target_folder\file1.mp4
  作成日時 (変更前): 2021-09-01 12:00:00
  作成日時 (変更後): 2021-09-01 21:00:00
  変更日時 (変更前): 2021-09-01 12:00:00
  変更日時 (変更後): 2021-09-01 21:00:00
  アクセス日時 (変更前): 2021-09-01 12:00:00
  アクセス日時 (変更後): 2021-09-01 21:00:00
...
実際の処理を実行しますか？ (y/n): n
処理がキャンセルされました。
```

## 使用例

### 単一ファイルの同期

`target.mp4`のタイムスタンプを`reference.mp4`に同期する場合:

```bash
SyncTimeStamp.exe -t "C:\path\to\target.mp4" -r "C:\path\to\reference.mp4"
```

### ディレクトリの同期

`target_folder`内のすべてのファイルのタイムスタンプを、`reference_folder`内の対応するファイルに同期する場合:

```bash
SyncTimeStamp.exe -t "C:\path\to\target_folder" -r "C:\path\to\reference_folder"
```

**注**: ファイルは拡張子を無視し、ファイル名に基づいて一致します。ターゲットファイルのファイル名が参照ファイル名を含んでいれば、対応するファイルと見なされます。

### 時間シフトの適用

+9時間の時間シフトを適用して同期（例：UTCからJSTに調整）する場合:

```bash
SyncTimeStamp.exe -t "C:\path\to\target_folder" -r "C:\path\to\reference_folder" -shift 9
```

### テストモードでの実行

ファイルを変更せずにプレビューのみ表示する場合:

```bash
SyncTimeStamp.exe -t "C:\path\to\target_folder" -r "C:\path\to\reference_folder" -test
```

## ソースコードからのビルド

`SyncTimeStamp.exe`をソースコードからビルドするには、[Go](https://golang.org/dl/)がインストールされている必要があります。

### ビルドコマンド

```bash
go build -o "SyncTimeStamp.exe"
```

または、提供されている`build.cmd`スクリプトを使用できます。

```bash
build.cmd
```

## ライセンス

このプロジェクトはMITライセンスの下で提供されています。詳細は[LICENSE](LICENSE)ファイルをご覧ください。

## コントリビューション

コントリビューションは歓迎です！GitHubでIssueの作成やPull Requestの送信をお待ちしています。

## 免責事項

このツールの使用は自己責任で行ってください。ファイルを変更する前に必ずバックアップを作成してください。

---