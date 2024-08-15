package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	var clearInstall bool
	flag.BoolVar(&clearInstall, "clear", false, "Clear install directory before installation")
	flag.Parse()

	binaryName := "eam"
	binDir := "bin"
	installDir := "install"

	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		fmt.Println("Error: bin directory does not exist")
		os.Exit(1)
	}

	if clearInstall {
		fmt.Println("Clearing install directory...")
		if err := os.RemoveAll(installDir); err != nil {
			fmt.Printf("Failed to clear install directory: %v\n", err)
			os.Exit(1)
		}
	}

	if err := os.MkdirAll(installDir, 0755); err != nil {
		fmt.Printf("Failed to create install directory: %v\n", err)
		os.Exit(1)
	}

	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	srcPath := filepath.Join(binDir, binaryName)
	dstPath := filepath.Join(installDir, binaryName)
	if err := copyFile(srcPath, dstPath); err != nil {
		fmt.Printf("Failed to copy binary file %s: %v\n", binaryName, err)
		os.Exit(1)
	}

	depsDir := "deps"
	depFiles := []string{
		"fastdeploy_ppocr_maa",
		"MaaAdbControlUnit",
		"MaaFramework",
		"MaaToolkit",
		"MaaUtils",
		"onnxruntime_maa",
		"opencv_world4_maa",
	}

	var libSuffix string
	switch runtime.GOOS {
	case "windows":
		libSuffix = ".dll"
	case "darwin":
		libSuffix = ".dylib"
	default:
		libSuffix = ".so"
	}

	depsBinDir := filepath.Join(depsDir, "bin")
	for _, file := range depFiles {
		srcPath = filepath.Join(depsBinDir, file+libSuffix)
		dstPath = filepath.Join(installDir, file+libSuffix)
		if err := copyFile(srcPath, dstPath); err != nil {
			fmt.Printf("Failed to copy dependency file %s: %v\n", file+libSuffix, err)
			os.Exit(1)
		}
	}

	MaaAgentBinaryDir := filepath.Join(depsDir, "share", "MaaAgentBinary")
	dstPath = filepath.Join(installDir, "MaaAgentBinary")
	if err := copyDir(MaaAgentBinaryDir, dstPath); err != nil {
		fmt.Printf("Failed to copy MaaAgentBinary directory %s: %v\n", MaaAgentBinaryDir, err)
		os.Exit(1)
	}

	fmt.Println("Install completed successfully.")
}

func copyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 构建目标路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			// 创建目标目录
			if err := os.MkdirAll(dstPath, info.Mode()); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err := copyFile(path, dstPath); err != nil {
				return err
			}
		}

		return nil
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}
