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
	depsBinDir := filepath.Join(depsDir, "bin")
	if err := copyDir(depsBinDir, installDir); err != nil {
		fmt.Printf("Failed to copy deps binary directory: %v\n", err)
		os.Exit(1)
	}

	excludeLibs := []string{
		"MaaDbgControlUnit",
		"MaaWin32ControlUnit",
	}
	var libPrefix, libSuffix string
	piCli := "MaaPiCli"
	switch runtime.GOOS {
	case "windows":
		libSuffix = ".dll"
		piCli += ".exe"
	case "darwin":
		libPrefix = "lib"
		libSuffix = ".dylib"
		excludeLibs = []string{"MaaDbgControlUnit"}
	default:
		libPrefix = "lib"
		libSuffix = ".so"
		excludeLibs = []string{"MaaDbgControlUnit"}
	}

	for _, lib := range excludeLibs {
		libPath := filepath.Join(installDir, libPrefix+lib+libSuffix)
		if err := os.Remove(libPath); err != nil {
			fmt.Printf("Failed to remove %s: %v\n", libPath, err)
			os.Exit(1)
		}
	}
	piCliPath := filepath.Join(installDir, piCli)
	if err := os.Remove(piCliPath); err != nil {
		fmt.Printf("Failed to remove %s: %v\n", piCliPath, err)
		os.Exit(1)
	}

	maaAgentBinaryDir := filepath.Join(depsDir, "share", "MaaAgentBinary")
	dstPath = filepath.Join(installDir, "MaaAgentBinary")
	if err := copyDir(maaAgentBinaryDir, dstPath); err != nil {
		fmt.Printf("Failed to copy MaaAgentBinary directory %s: %v\n", maaAgentBinaryDir, err)
		os.Exit(1)
	}

	assetsDir := "assets"
	resDir := filepath.Join(assetsDir, "resource")
	dtsPath := filepath.Join(installDir, "resource")
	if err := copyDir(resDir, dtsPath); err != nil {
		fmt.Printf("Failed to copy resource directory %s: %v\n", resDir, err)
		os.Exit(1)
	}

	installConfigDir := filepath.Join(installDir, "config")
	if err := os.MkdirAll(installConfigDir, 0755); err != nil {
		fmt.Printf("Failed to create install config directory: %v\n", err)
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
