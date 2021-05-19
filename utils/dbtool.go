package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func PreparingQurySQL(columns []string, tableName string, offset int, limit int, order string, whereSQL string) string {
	byteBuf := bytes.NewBufferString("")
	byteBuf.WriteString(fmt.Sprintf("select %s from %s where 1=1 %s order by %s", strings.Join(columns, ","), tableName, whereSQL, order))
	if limit >= 0 {
		byteBuf.WriteString(" limit ")
		byteBuf.WriteString(strconv.Itoa(limit))
	}
	if offset >= 0 {
		byteBuf.WriteString(" offset ")
		byteBuf.WriteString(strconv.Itoa(offset))
	}
	fullSQL := byteBuf.String()
	return fullSQL
}

func GeneWhereLike(k, v string, index int32, sqlKeys []string, sqlValues []interface{}) ([]string, []interface{}, int32) {
	sqlKeys = append(sqlKeys, k+" LIKE '%"+v+"%'")
	sqlValues = append(sqlValues, strings.ToLower(strings.TrimSpace(v)))
	index++
	return sqlKeys, sqlValues, index
}

func GeneWhereEqual(k, v string, index int32, sqlKeys []string, sqlValues []interface{}) ([]string, []interface{}, int32) {
	sqlKeys = append(sqlKeys, fmt.Sprintf("%s = $%d", k, index))
	sqlValues = append(sqlValues, strings.TrimSpace(v))
	index++
	return sqlKeys, sqlValues, index
}

func GeneWhereInt(k, v string, index int32, sqlKeys []string, sqlValues []interface{}) ([]string, []interface{}, int32) {
	tempValue, err := strconv.Atoi(v)
	if nil != err {
		return sqlKeys, sqlValues, index
	}
	sqlKeys = append(sqlKeys, fmt.Sprintf("%s=$%d", k, index))
	sqlValues = append(sqlValues, tempValue)
	index++
	return sqlKeys, sqlValues, index
}

func GeneWhereBool(k, v string, index int32, sqlKeys []string, sqlValues []interface{}) ([]string, []interface{}, int32) {
	tempValue, err := strconv.ParseBool(v)
	if nil != err {
		return sqlKeys, sqlValues, index
	}
	sqlKeys = append(sqlKeys, fmt.Sprintf("%s=$%d", k, index))
	sqlValues = append(sqlValues, tempValue)
	index++
	return sqlKeys, sqlValues, index
}

func GeneWhereSQL(key []string, value []interface{}) (string, []interface{}, error) {
	if len(key) > 0 {
		return "and " + strings.Join(key, " and "), value, nil
	} else {
		return "", nil, nil
	}
}
