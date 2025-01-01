package chatsetting

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gf-chat/internal/cache"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterChatSetting(&sChatSetting{
		trait.Curd[model.CustomerChatSetting]{
			Dao: &dao.CustomerChatSettings,
		},
	})
}

const (
	nameCacheKey     = "customer:%d:setting:name"
	avatarCacheKey   = "customer:%d:setting:avatar"
	transferCacheKey = "customer:%d:setting:transfer"
	readCacheKey     = "customer:%d:setting:read"
	queueCacheKey    = "customer:%d:setting:queue"
)

type sChatSetting struct {
	trait.Curd[model.CustomerChatSetting]
}

func (s *sChatSetting) RemoveCache(ctx context.Context, customerId uint) error {
	_, err := cache.Def.Remove(ctx,
		fmt.Sprintf(nameCacheKey, customerId),
		fmt.Sprintf(avatarCacheKey, customerId),
		fmt.Sprintf(transferCacheKey, customerId),
		fmt.Sprintf(readCacheKey, customerId),
		fmt.Sprintf(queueCacheKey, customerId),
	)
	return err
}

func (s *sChatSetting) GetName(ctx context.Context, customerId uint) (name string, err error) {
	v, err := cache.Def.GetOrSetFunc(ctx, fmt.Sprintf(nameCacheKey, customerId), func(ctx context.Context) (r any, err error) {
		setting, err := s.First(ctx, do.CustomerChatSettings{
			CustomerId: customerId,
			Name:       consts.ChatSettingSystemName,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", nil
			}
		}
		return setting.Value, nil
	}, time.Minute*10)
	if err != nil {
		return
	}
	return v.String(), nil

}
func (s *sChatSetting) GetIsUserShowRead(ctx context.Context, customerId uint) (isShow bool, err error) {
	v, err := cache.Def.GetOrSetFunc(ctx, fmt.Sprintf(readCacheKey, customerId), func(ctx context.Context) (r any, err error) {
		setting, err := s.First(ctx, do.CustomerChatSettings{
			CustomerId: customerId,
			Name:       consts.ChatSettingShowRead,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}
			return
		}
		return setting.Value == "1", nil
	}, time.Minute*10)
	if err != nil {
		return
	}
	return v.Bool(), nil
}
func (s *sChatSetting) GetIsUserShowQueue(ctx context.Context, customerId uint) (isShow bool, err error) {
	v, err := cache.Def.GetOrSetFunc(ctx, fmt.Sprintf(queueCacheKey, customerId), func(ctx context.Context) (r any, err error) {
		setting, err := s.First(ctx, do.CustomerChatSettings{
			CustomerId: customerId,
			Name:       consts.ChatSettingShowQueue,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}
			return
		}
		return setting.Value == "1", nil
	}, time.Minute*10)
	if err != nil {
		return
	}
	return v.Bool(), nil
}
func (s *sChatSetting) GetAvatar(ctx context.Context, customerId uint) (name string, err error) {
	v, err := cache.Def.GetOrSetFunc(ctx, fmt.Sprintf(avatarCacheKey, customerId), func(ctx context.Context) (r any, err error) {
		setting, err := s.First(ctx, do.CustomerChatSettings{
			CustomerId: customerId,
			Name:       consts.ChatSettingSystemAvatar,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", nil
			}
			return
		}
		file, err := service.File().First(ctx, do.CustomerChatFiles{
			Id: setting.Value,
		})
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		fileUrl := service.File().ToApi(file).Url
		return fileUrl, nil
	}, time.Minute*10)
	if err != nil {
		return
	}
	return v.String(), nil
}

// GetIsAutoTransferManual 是否自动转接人工客服
func (s *sChatSetting) GetIsAutoTransferManual(ctx context.Context, customerId uint) (b bool, err error) {
	v, err := cache.Def.GetOrSetFunc(ctx, fmt.Sprintf(transferCacheKey, customerId), func(ctx context.Context) (r any, err error) {
		setting, err := s.First(ctx, do.CustomerChatSettings{
			CustomerId: customerId,
			Name:       consts.ChatSettingIsAutoTransfer,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}
			return
		}
		return gconv.Bool(gconv.Int(setting.Value)), nil
	}, time.Minute*10)
	if err != nil {
		return
	}
	return v.Bool(), nil
}
