/**
 * @Author: upingli
 * @Description:
 * @File:  zip_utils
 * @Version: 1.0.0
 * @Date: 2020/12/26 10:57 上午
 */
package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Zip(srcFileName, destFileName string) error {
	// 创建压缩文件目录
	zipFile, err := os.Create(destFileName)
	if err != nil {
		fmt.Println("创建压缩文件失败", err)
		return err
	}
	defer zipFile.Close()
	// 创建管道
	writer := zip.NewWriter(zipFile)
	defer writer.Close()
	// 遍历目录写入文件
	filepath.Walk(srcFileName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("遍历文件失败")
			return err
		}
		// zip获取文件信息
		fHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			fmt.Println("获取文件信息失败")
			return err
		}
		// 文件夹
		fHeader.Name = path
		if info.IsDir() {
			fHeader.Name += "/"
		} else { // 不是文件则
			fHeader.Method = zip.Deflate
		}
		// 创建文件信息
		w, err := writer.CreateHeader(fHeader)
		if err != nil {
			fmt.Println("文件信息写入失败", err)
			return err
		}

		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				fmt.Println("打开文件失败", err)
				return err
			}
			defer f.Close()
			_, err = io.Copy(w, f)
			if err != nil {
				fmt.Println("文件复制失败")
				return err
			}
		}
		return nil
	})
	return nil
}

// 文件解压

func UnZip(zipFile, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		fmt.Println("文件读取失败")
		return err
	}
	defer zipReader.Close()
	for _, f := range zipReader.File {
		// 文件路径
		path := filepath.Join(destDir, f.Name)
		// 是文件夹则创建文件夹
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return err
			}
			// 打开文件
			freader, err := f.Open()
			if err != nil {
				fmt.Println("文件打开失败")
				return err
			}
			defer freader.Close()
			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				fmt.Println("保存文件创建失败")
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, freader)
			if err != nil {
				fmt.Println("文件复制失败")
				return err
			}
		}
	}
	return nil
}