package internal

import (
	"encoding/csv"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"json/pkg/db"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
)

func ReadTestTemplateFromCSV(path string) (result []RequestTemplateOne, err error) {
	// read csv
	f, err := os.Open(path)
	if err != nil {
		Trace.Println(err)
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		Trace.Println(err)
		return nil, err
	}
	records = records[1:]
	for _, row := range records {
		ParseType, err := strconv.Atoi(row[2])
		if err != nil {
			Trace.Println(err)
			return nil, err
		}
		AutoGenObjectID, err := strconv.ParseBool(row[3])
		if err != nil {
			Trace.Println(err)
			return nil, err
		}
		result = append(result, RequestTemplateOne{
			CreateTemplateParams: db.CreateTemplateParams{
				TemplateName:    row[0],
				Description:     row[1],
				ParseType:       int32(ParseType),
				AutoGenObjectID: AutoGenObjectID,
				OriginObjectID:  row[4],
				ValueExtract:    row[6],
				GoStruct:        row[7],
			},
			ObjectIDList: strings.Split(row[5], ","),
		})
	}
	return
}

func clearTemplateDB(t *testing.T) {
	_, err := MyDB.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	require.NoError(t, err)
	_, err = MyDB.Exec("TRUNCATE template;")
	require.NoError(t, err)
	_, err = MyDB.Exec("TRUNCATE objectids")
	require.NoError(t, err)
	_, err = MyDB.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	require.NoError(t, err)
}

func TestCreateScheduleTX(t *testing.T) {
	// 先清空DB
	clearTemplateDB(t)
	templates, err := ReadTestTemplateFromCSV(path.Join(rootPath, "test_data", "test_create.csv"))
	require.NoError(t, err)

	store := NewStore(MyDB)
	t.Run("Create type0 單值", func(t *testing.T) {
		template := templates[0]
		iid, err := store.CreateTemplateTX(template)
		require.NoError(t, err)
		// 檢查當前最大的template Id是否就是這次建立產生的
		var want int64
		row := MyDB.QueryRow("SELECT Max(id) FROM template")
		err = row.Scan(&want)
		require.NoError(t, err)
		require.Equal(t, want, iid)
	})
	t.Run("Create type1 陣列值", func(t *testing.T) {
		template := templates[1]
		iid, err := store.CreateTemplateTX(template)
		require.NoError(t, err)
		// 檢查當前最大的template Id是否就是這次建立產生的
		var want int64
		row := MyDB.QueryRow("SELECT Max(id) FROM template")
		err = row.Scan(&want)
		require.NoError(t, err)
		require.Equal(t, want, iid)
	})
	t.Run("type0 自動產生object_id", func(t *testing.T) {
		template := templates[2]
		iid, err := store.CreateTemplateTX(template)
		require.NoError(t, err)
		// 檢查當前最大的template Id是否就是這次建立產生的
		var want int64
		row := MyDB.QueryRow("SELECT Max(id) FROM template")
		err = row.Scan(&want)
		require.NoError(t, err)
		require.Equal(t, want, iid)
	})
	t.Run("Create type1 刪除使用", func(t *testing.T) {
		template := templates[3]
		iid, err := store.CreateTemplateTX(template)
		require.NoError(t, err)
		// 檢查當前最大的template Id是否就是這次建立產生的
		var want int64
		row := MyDB.QueryRow("SELECT Max(id) FROM template")
		err = row.Scan(&want)
		require.NoError(t, err)
		require.Equal(t, want, iid)
	})
	t.Run("真實數據結構", func(t *testing.T) {
		template := templates[4]
		iid, err := store.CreateTemplateTX(template)
		require.NoError(t, err)
		// 檢查當前最大的template Id是否就是這次建立產生的
		var want int64
		row := MyDB.QueryRow("SELECT Max(id) FROM template")
		err = row.Scan(&want)
		require.NoError(t, err)
		require.Equal(t, want, iid)
	})
	t.Run("真實數據結構", func(t *testing.T) {
		template := templates[5]
		iid, err := store.CreateTemplateTX(template)
		require.NoError(t, err)
		// 檢查當前最大的template Id是否就是這次建立產生的
		var want int64
		row := MyDB.QueryRow("SELECT Max(id) FROM template")
		err = row.Scan(&want)
		require.NoError(t, err)
		require.Equal(t, want, iid)
	})
	t.Run("真實數據結構", func(t *testing.T) {
		template := templates[6]
		iid, err := store.CreateTemplateTX(template)
		require.NoError(t, err)
		// 檢查當前最大的template Id是否就是這次建立產生的
		var want int64
		row := MyDB.QueryRow("SELECT Max(id) FROM template")
		err = row.Scan(&want)
		require.NoError(t, err)
		require.Equal(t, want, iid)
	})
	t.Run("真實數據結構", func(t *testing.T) {
		template := templates[7]
		iid, err := store.CreateTemplateTX(template)
		require.NoError(t, err)
		// 檢查當前最大的template Id是否就是這次建立產生的
		var want int64
		row := MyDB.QueryRow("SELECT Max(id) FROM template")
		err = row.Scan(&want)
		require.NoError(t, err)
		require.Equal(t, want, iid)
	})
	t.Run("真實複雜數據結構", func(t *testing.T) {
		template := templates[8]
		iid, err := store.CreateTemplateTX(template)
		require.NoError(t, err)
		// 檢查當前最大的template Id是否就是這次建立產生的
		var want int64
		row := MyDB.QueryRow("SELECT Max(id) FROM template")
		err = row.Scan(&want)
		require.NoError(t, err)
		require.Equal(t, want, iid)
	})
}

func TestDeleteTemplateTx(t *testing.T) {
	store := NewStore(MyDB)
	err := store.DeleteTemplateTx(4)
	require.NoError(t, err)
	// 檢查當前最大的template Id
	var got int64
	row := MyDB.QueryRow("SELECT Max(id) FROM template")
	err = row.Scan(&got)
	require.NoError(t, err)
	//require.Equal(t, 8, int(got))

}

func TestUpdateTemplateTx(t *testing.T) {
	store := NewStore(MyDB)
	templates, err := ReadTestTemplateFromCSV(path.Join(rootPath, "test_data", "test_update.csv"))
	require.NoError(t, err)
	template := RequestTemplateOne{
		CreateTemplateParams: db.CreateTemplateParams{
			TemplateName:    templates[0].TemplateName,
			Description:     templates[0].Description,
			ParseType:       templates[0].ParseType,
			AutoGenObjectID: templates[0].AutoGenObjectID,
			OriginObjectID:  templates[0].OriginObjectID,
			ValueExtract:    templates[0].ValueExtract,
			GoStruct:        templates[0].GoStruct,
		},
		ObjectIDList: templates[0].ObjectIDList,
		TemplateID:   int32(1),
	}
	err = store.UpdateTemplateTx(template)
	require.NoError(t, err)
	// 檢查當前最大的template Id
	var got string
	row := MyDB.QueryRow("SELECT Description FROM template where id=1")
	err = row.Scan(&got)
	require.NoError(t, err)
	require.Equal(t, "更新專用", got)
}

func TestListTemplates(t *testing.T) {
	results, err := ListTemplates()
	require.NoError(t, err)
	var want int
	row := MyDB.QueryRow("SELECT count(1) FROM template")
	err = row.Scan(&want)
	require.NoError(t, err)
	require.Equal(t, want, len(results))
}

func TestGetSingleTemplateTX(t *testing.T) {
	TemplateId := 1
	store := NewStore(MyDB)
	result, err := store.GetSingleTemplateTX(int32(TemplateId))
	require.NoError(t, err)
	var want int
	row := MyDB.QueryRow("SELECT parse_type FROM template where id=?", TemplateId)
	err = row.Scan(&want)
	require.Equal(t, want, int(result.ParseType))
	row = MyDB.QueryRow("SELECT count(1) FROM objectids where template_id=?", TemplateId)
	err = row.Scan(&want)
	require.Equal(t, want, len(result.ObjectIDList))
}
