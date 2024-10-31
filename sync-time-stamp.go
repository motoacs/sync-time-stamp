package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

func stripQuotes(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"")
	return s
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type FileTimes struct {
	creationTime     time.Time
	modificationTime time.Time
	accessTime       time.Time
}

func getFileTimes(filePath string) (FileTimes, error) {
	var times FileTimes

	pathp, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		return times, err
	}

	handle, err := windows.CreateFile(
		pathp,
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0)
	if err != nil {
		return times, err
	}
	defer windows.CloseHandle(handle)

	var ctime, atime, mtime windows.Filetime

	err = windows.GetFileTime(handle, &ctime, &atime, &mtime)
	if err != nil {
		return times, err
	}

	times.creationTime = time.Unix(0, ctime.Nanoseconds())
	times.accessTime = time.Unix(0, atime.Nanoseconds())
	times.modificationTime = time.Unix(0, mtime.Nanoseconds())

	return times, nil
}

func shiftFileTimes(times FileTimes, shiftHours int) FileTimes {
	shiftedTimes := FileTimes{
		creationTime:     times.creationTime.Add(time.Duration(shiftHours) * time.Hour),
		modificationTime: times.modificationTime.Add(time.Duration(shiftHours) * time.Hour),
		accessTime:       times.accessTime.Add(time.Duration(shiftHours) * time.Hour),
	}
	return shiftedTimes
}

func setFileTimes(filePath string, times FileTimes) error {
	pathp, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		return err
	}

	handle, err := windows.CreateFile(
		pathp,
		windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)

	ctime := windows.NsecToFiletime(times.creationTime.UnixNano())
	atime := windows.NsecToFiletime(times.accessTime.UnixNano())
	mtime := windows.NsecToFiletime(times.modificationTime.UnixNano())

	err = windows.SetFileTime(handle, &ctime, &atime, &mtime)
	if err != nil {
		return err
	}

	return nil
}

func processFile(targetFilePath string, referenceFilePath string, shiftHours int, testMode bool) error {
	// 参照ファイルのタイムスタンプを取得
	refTimes, err := getFileTimes(referenceFilePath)
	if err != nil {
		return fmt.Errorf("参照ファイルのタイムスタンプ取得エラー: %v", err)
	}
	targetTimes, err := getFileTimes(targetFilePath)
	if err != nil {
		return fmt.Errorf("ターゲットファイルのタイムスタンプ取得エラー: %v", err)
	}

	shiftedTimes := shiftFileTimes(refTimes, shiftHours)

	// 変更内容を表示
	fmt.Printf("ファイル: %s\n", targetFilePath)
	fmt.Printf("  作成日時 (変更前): %v\n", targetTimes.creationTime)
	fmt.Printf("  作成日時 (変更後): %v\n", shiftedTimes.creationTime)
	fmt.Printf("  変更日時 (変更前): %v\n", targetTimes.modificationTime)
	fmt.Printf("  変更日時 (変更後): %v\n", shiftedTimes.modificationTime)
	fmt.Printf("  アクセス日時 (変更前): %v\n", targetTimes.accessTime)
	fmt.Printf("  アクセス日時 (変更後): %v\n", shiftedTimes.accessTime)

	if !testMode {
		// ターゲットファイルのタイムスタンプを更新
		err = setFileTimes(targetFilePath, shiftedTimes)
		if err != nil {
			return fmt.Errorf("ファイルタイムスタンプ設定エラー: %v", err)
		}
	}

	return nil
}

func processDirectories(targetDir string, referenceDir string, shiftHours int, testMode bool) error {
	// ターゲットディレクトリのファイルを取得
	targetFiles, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("ターゲットディレクトリ読み込みエラー: %v", err)
	}

	// 参照ディレクトリのファイルを取得
	referenceFiles, err := ioutil.ReadDir(referenceDir)
	if err != nil {
		return fmt.Errorf("参照ディレクトリ読み込みエラー: %v", err)
	}

	// 参照ファイルのマップを作成（拡張子を除くファイル名をキーとする）
	refFileMap := make(map[string]string)
	for _, refFile := range referenceFiles {
		if !refFile.IsDir() {
			name := strings.TrimSuffix(refFile.Name(), filepath.Ext(refFile.Name()))
			refFileMap[name] = filepath.Join(referenceDir, refFile.Name())
		}
	}

	// ターゲットファイルごとに処理
	for _, targetFile := range targetFiles {
		if !targetFile.IsDir() {
			targetName := strings.TrimSuffix(targetFile.Name(), filepath.Ext(targetFile.Name()))
			matched := false
			for refName, refPath := range refFileMap {
				if strings.Contains(targetName, refName) {
					targetFilePath := filepath.Join(targetDir, targetFile.Name())
					err := processFile(targetFilePath, refPath, shiftHours, testMode)
					if err != nil {
						fmt.Println("ファイル処理エラー:", err)
					}
					matched = true
					break
				}
			}
			if !matched {
				fmt.Printf("対応する参照ファイルが見つかりません: %s\n", targetFile.Name())
			}
		}
	}

	return nil
}

func processFiles(targetPath string, referencePath string, shiftHours int, testMode bool) error {
	targetInfo, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("ターゲットパスアクセスエラー: %v", err)
	}
	referenceInfo, err := os.Stat(referencePath)
	if err != nil {
		return fmt.Errorf("参照パスアクセスエラー: %v", err)
	}

	if targetInfo.IsDir() != referenceInfo.IsDir() {
		return fmt.Errorf("両方のパスは同じタイプである必要があります（ファイル同士またはディレクトリ同士）")
	}

	if targetInfo.IsDir() {
		return processDirectories(targetPath, referencePath, shiftHours, testMode)
	} else {
		return processFile(targetPath, referencePath, shiftHours, testMode)
	}
}

func main() {
	var (
		targetPathFlag    = flag.String("t", "", "ターゲットディレクトリまたはファイルのパス")
		referencePathFlag = flag.String("r", "", "参照ディレクトリまたはファイルのパス")
		shiftHoursFlag    = flag.Int("shift", 0, "タイムゾーン調整のための時間シフト")
		testMode          = flag.Bool("test", false, "テストモードで実行")
	)

	flag.Parse()

	targetPath := stripQuotes(*targetPathFlag)
	referencePath := stripQuotes(*referencePathFlag)
	shiftHours := *shiftHoursFlag

	reader := bufio.NewReader(os.Stdin)

	if targetPath == "" {
		for {
			fmt.Print("ターゲットディレクトリまたはファイルのパスを入力 (-t): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			input = stripQuotes(input)
			if input != "" {
				targetPath = input
				break
			}
			fmt.Println("エラー: ターゲットパスは空にできません。")
		}
	}

	if referencePath == "" {
		for {
			fmt.Print("参照ディレクトリまたはファイルのパスを入力 (-r): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			input = stripQuotes(input)
			if input != "" {
				referencePath = input
				break
			}
			fmt.Println("エラー: 参照パスは空にできません。")
		}
	}

	// シフト値が提供されていない場合、デフォルトを0として入力を促す
	shiftHoursFlagProvided := false
	flag.CommandLine.Visit(func(f *flag.Flag) {
		if f.Name == "shift" {
			shiftHoursFlagProvided = true
		}
	})

	if !shiftHoursFlagProvided {
		for {
			fmt.Print("タイムゾーン調整のための時間シフト (-24から24) [デフォルト 0]: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "" {
				shiftHours = 0
				break
			} else {
				s, err := strconv.Atoi(input)
				if err != nil || s < -24 || s > 24 {
					fmt.Println("エラー: シフト時間は-24から24の数値でなければなりません。")
				} else {
					shiftHours = s
					break
				}
			}
		}
	} else {
		if shiftHours < -24 || shiftHours > 24 {
			fmt.Println("エラー: --shift の値は -24 から 24 の間でなければなりません")
			os.Exit(1)
		}
	}

	// パスの存在確認
	targetExists, err := pathExists(targetPath)
	if err != nil {
		fmt.Println("ターゲットパスアクセスエラー:", err)
		os.Exit(1)
	}
	if !targetExists {
		fmt.Println("エラー: ターゲットパスが存在しません。")
		os.Exit(1)
	}

	referenceExists, err := pathExists(referencePath)
	if err != nil {
		fmt.Println("参照パスアクセスエラー:", err)
		os.Exit(1)
	}
	if !referenceExists {
		fmt.Println("エラー: 参照パスが存在しません。")
		os.Exit(1)
	}

	if *testMode {
		fmt.Println("[[[ テストモード ]]]")
		err = processFiles(targetPath, referencePath, shiftHours, true)
		if err != nil {
			fmt.Println("テストモード処理中のエラー:", err)
			os.Exit(1)
		}
		fmt.Println("テストモードが完了しました。")
	} else {
		fmt.Println("テストモードで実行中...")
		err = processFiles(targetPath, referencePath, shiftHours, true)
		if err != nil {
			fmt.Println("テストモード処理中のエラー:", err)
			os.Exit(1)
		}

		fmt.Print("実際の処理を実行しますか？ (y/n): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "y" || input == "Y" {
			err = processFiles(targetPath, referencePath, shiftHours, false)
			if err != nil {
				fmt.Println("処理中のエラー:", err)
				os.Exit(1)
			}
			fmt.Println("処理が正常に完了しました。")
		} else {
			fmt.Println("処理がキャンセルされました。")
		}
	}
}
