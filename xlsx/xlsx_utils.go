/**
 * @Author: upingli
 * @Description:
 * @File:  xlsx_utils.go
 * @Version: 1.0.0
 * @Date: 2020/12/26 10:59 上午
 */
package xlsx

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"path/filepath"
)

func ParsingXlsx(path string) ([]map[string]string, error){
	file ,err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Println("当前文件打开失败", err)
		return nil, err
	}
	rowData := []map[string]string{}
	for _, sheet := range file.Sheets {
		firstRow := true
		firstRowData := make([]string, 0)
		err := sheet.ForEachRow(func(r *xlsx.Row) error {
			cellData := make(map[string]string, 0)
			count := 0
			err := r.ForEachCell(func(c *xlsx.Cell) error {
				if firstRow {
					firstRowData = append(firstRowData, c.Value)
				} else {
					index := count % (len(firstRowData))
					key := firstRowData[index]
					cellData[key] = c.Value
					count += 1
				}
				return nil
			})
			if err != nil {
				fmt.Println("读取当前cell失败", err)
				return err
			}
			if !firstRow {
				rowData = append(rowData, cellData)
			} else {
				firstRow = false
			}
			return nil
		})
		if err != nil {
			fmt.Println("读取当前sheet失败", err)
			return nil, err
		}
		fmt.Println("当前sheet的数据", rowData)
	}
	return rowData, nil
}

func CreateXlsxFile(data []map[string]string, fileName, path string) error{
	xlsxFile := xlsx.NewFile()
	keys := getMapKeys(data[0])
	sheet, err := xlsxFile.AddSheet("Shee1")
	if err != nil {
		fmt.Println("创建xlsx文件失败", err)
		return err
	}
	defer sheet.Close()
	// 创建第一行
	firstRow := sheet.AddRow()
	for _, v := range keys {
		cell := firstRow.AddCell()
		cell.Value = v
	}
	for _, d := range data {
		newRow := sheet.AddRow()
		for _, key := range keys {
			newCell := newRow.AddCell()
			newCell.Value = d[key]
			fmt.Println("当前数据", d[key])
		}
	}
	err = xlsxFile.Save(filepath.Join(path, fileName))
	fmt.Println("文件保存路径", filepath.Join(path, fileName))
	if err != nil {
		fmt.Println("文件保存失败")
		return err
	}
	return nil
}

func getMapKeys(data map[string]string) []string{
	keys := make([]string, 0)
	for k, _ := range data {
		keys = append(keys, k)
	}
	return keys
}