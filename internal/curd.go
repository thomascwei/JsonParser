package internal

import (
	"context"
	"database/sql"
	"fmt"
	"json/pkg/db"
	"time"
)

// Store defines all functions to execute db queries and transactions
type Store struct {
	*db.Queries
	TransDB *sql.DB
}

// NewStore creates a new store
func NewStore(sdb *sql.DB) *Store {
	return &Store{
		Queries: db.New(sdb),
		TransDB: sdb,
	}
}

// ExecTx executes a function within a database transaction
func (s *Store) execTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := s.TransDB.BeginTx(ctx, &TxOption)
	if err != nil {
		return err
	}
	q := db.New(tx)
	err = fn(q)
	if err != nil {
		sugarLogger.Errorf("Rollback")
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func CreateTemplateRawTX(rto RequestTemplateOne) (LID int64, err error) {
	template := rto.CreateTemplateParams
	tx, err := MyDB.BeginTx(ctx, &TxOption)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		tx.Rollback()
		return
	}
	// 連續性操作實作於此
	result, err := tx.Exec(db.CreateTemplate, template.TemplateName, template.Description, template.ParseType,
		template.AutoGenObjectID, template.OriginObjectID, template.ValueExtract, template.GoStruct)
	LID, err = result.LastInsertId()
	if err != nil {
		sugarLogger.Errorf(err.Error())
		tx.Rollback()
		return
	}
	if rto.AutoGenObjectID == true {
		return
	}
	// 建object_id
	objParam := db.CreateObjectIDsParams{
		TemplateID: int32(LID),
		Serial:     0,
		ObjectID:   "",
	}
	for no, id := range rto.ObjectIDList {
		objParam.Serial = int32(no + 1)
		objParam.ObjectID = id
		_, err = tx.Exec(db.CreateObjectIDs, objParam.TemplateID, objParam.Serial, objParam.ObjectID)
		if err != nil {
			sugarLogger.Errorf(err.Error())
			tx.Rollback()
			return
		}
	}

	return LID, tx.Commit()
}

// CreateTemplateTX all related data within astructtabase transaction
func (s *Store) CreateTemplateTX(rto RequestTemplateOne) (int64, error) {
	var LID int64
	template := rto.CreateTemplateParams
	err := s.execTx(ctx, func(q *db.Queries) error {
		// 連續性操作實作於此
		result, err := queries.CreateTemplate(ctx, template)
		if err != nil {
			sugarLogger.Errorf(err.Error())
			return err
		}
		LID, err = result.LastInsertId()
		if err != nil {
			sugarLogger.Errorf(err.Error())
			return err
		}
		// AutoGenObjectID若為真不建object_id
		if rto.AutoGenObjectID == true {
			return err
		}
		// 建object_id
		objParam := db.CreateObjectIDsParams{
			TemplateID: int32(LID),
			Serial:     0,
			ObjectID:   "",
		}
		for no, id := range rto.ObjectIDList {
			objParam.Serial = int32(no + 1)
			objParam.ObjectID = id
			_, err = queries.CreateObjectIDs(ctx, objParam)
			if err != nil {
				sugarLogger.Errorf(err.Error())
				return err
			}
		}
		return err
	})
	if err != nil {
		return -1, err
	}
	return LID, err
}

func DeleteTemplateRawTx(id int32) (err error) {
	Trace.Println(id)
	//tx, err := MyDB.BeginTx(ctx, &TxOption)
	tx, err := MyDB.BeginTx(ctx, &TxOption)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		return
	}
	// 連續性操作實作於此
	_, err = tx.ExecContext(ctx, db.DeleteObjectIDs, id)
	//err = queries.DeleteObjectIDs(ctx, id)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		tx.Rollback()
		return
	}
	Trace.Println("")
	//err = queries.DeleteTemplate(ctx, id)
	_, err = tx.ExecContext(ctx, db.DeleteTemplate, id)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}

// DeleteTemplateTx delete all related data within a database transaction
func (s *Store) DeleteTemplateTx(id int32) error {
	err := s.execTx(ctx, func(q *db.Queries) error {
		// 連續性操作實作於此
		err := queries.DeleteObjectIDs(ctx, id)
		if err != nil {
			sugarLogger.Errorf(err.Error())
			return err
		}
		err = queries.DeleteTemplate(ctx, id)
		if err != nil {
			sugarLogger.Errorf(err.Error())
			return err
		}
		return nil
	})
	return err
}

// 先更新template然後刪除再新增object_id
func (s *Store) UpdateTemplateTx(rto RequestTemplateOne) error {
	err := s.execTx(ctx, func(q *db.Queries) error {
		// 連續性操作實作於此
		err := queries.UpdateTemplate(ctx, db.UpdateTemplateParams{
			TemplateName:    rto.TemplateName,
			Description:     rto.Description,
			ParseType:       rto.ParseType,
			AutoGenObjectID: rto.AutoGenObjectID,
			OriginObjectID:  rto.OriginObjectID,
			ValueExtract:    rto.ValueExtract,
			GoStruct:        rto.GoStruct,
			UpdatedAt:       time.Now(),
			ID:              rto.TemplateID,
		})
		if err != nil {
			sugarLogger.Errorf(err.Error())
			return err
		}
		err = queries.DeleteObjectIDs(ctx, rto.TemplateID)
		if err != nil {
			sugarLogger.Errorf(err.Error())
			return err
		}
		// 建object_id
		objParam := db.CreateObjectIDsParams{
			TemplateID: rto.TemplateID,
			Serial:     0,
			ObjectID:   "",
		}
		for no, id := range rto.ObjectIDList {
			objParam.Serial = int32(no + 1)
			objParam.ObjectID = id
			_, err = queries.CreateObjectIDs(ctx, objParam)
			if err != nil {
				sugarLogger.Errorf(err.Error())
				return err
			}
		}
		return nil
	})
	return err
}

// 查詢template清單, 只返回編號名稱敘述
func ListTemplates() (result []db.ListTemplatesRow, err error) {
	result, err = queries.ListTemplates(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 查詢單一template, objecid固定返回array
func (s *Store) GetSingleTemplateTX(RequestId int32) (result RequestTemplateOne, err error) {
	err = s.execTx(ctx, func(q *db.Queries) error {
		// 連續性操作實作於此
		template, err := queries.GetTemplate(ctx, RequestId)
		if err != nil {
			return err
		}
		ids, err := queries.GetObjectIDs(ctx, RequestId)
		if err != nil {
			return err
		}
		ObjectIDList := make([]string, 0, 0)
		for _, id := range ids {
			ObjectIDList = append(ObjectIDList, id.ObjectID)
		}
		result.ObjectIDList = ObjectIDList
		result.TemplateName = template.TemplateName
		result.Description = template.Description
		result.ParseType = template.ParseType
		result.AutoGenObjectID = template.AutoGenObjectID
		result.OriginObjectID = template.OriginObjectID
		result.ValueExtract = template.ValueExtract
		result.GoStruct = template.GoStruct
		return nil
	})
	if err != nil {
		return RequestTemplateOne{}, err
	}
	return result, nil

}
