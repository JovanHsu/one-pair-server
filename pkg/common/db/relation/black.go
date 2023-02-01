package relation

import (
	"Open_IM/pkg/common/db/table"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type BlackGorm struct {
	DB *gorm.DB
}

func NewBlackGorm(db *gorm.DB) *BlackGorm {
	var black BlackGorm
	black.DB = db
	return &black
}

func (b *BlackGorm) Create(ctx context.Context, blacks []*table.BlackModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks)
	}()
	return utils.Wrap(b.DB.Model(&table.BlackModel{}).Create(&blacks).Error, "")
}

func (b *BlackGorm) Delete(ctx context.Context, blacks []*table.BlackModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks)
	}()
	return utils.Wrap(b.DB.Model(&table.BlackModel{}).Delete(blacks).Error, "")
}

func (b *BlackGorm) UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "blockUserID", blockUserID, "args", args)
	}()
	return utils.Wrap(b.DB.Model(&table.BlackModel{}).Where("block_user_id = ? and block_user_id = ?", ownerUserID, blockUserID).Updates(args).Error, "")
}

func (b *BlackGorm) Update(ctx context.Context, blacks []*table.BlackModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks)
	}()
	return utils.Wrap(b.DB.Model(&table.BlackModel{}).Updates(&blacks).Error, "")
}

func (b *BlackGorm) Find(ctx context.Context, blacks []*table.BlackModel) (blackList []*table.BlackModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks, "blackList", blackList)
	}()
	var where [][]interface{}
	for _, black := range blacks {
		where = append(where, []interface{}{black.OwnerUserID, black.BlockUserID})
	}
	return blackList, utils.Wrap(b.DB.Model(&table.BlackModel{}).Where("(owner_user_id, block_user_id) in ?", where).Find(&blackList).Error, "")
}

func (b *BlackGorm) Take(ctx context.Context, ownerUserID, blockUserID string) (black *table.BlackModel, err error) {
	black = &table.BlackModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "blockUserID", blockUserID, "black", *black)
	}()
	return black, utils.Wrap(b.DB.Model(&table.BlackModel{}).Where("owner_user_id = ? and block_user_id = ?", ownerUserID, blockUserID).Take(black).Error, "")
}

func (b *BlackGorm) FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blacks []*table.BlackModel, total int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "blacks", blacks)
	}()
	err = b.DB.Model(&table.BlackModel{}).Model(b).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(b.DB.Model(&table.BlackModel{}).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&blacks).Error, "")
	return
}
