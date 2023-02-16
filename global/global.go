package global

import (
	"context"
	"github.com/SCUTKing/config"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/sony/sonyflake"
	"gorm.io/gorm"
)

var (
	CONFIG               config.System            // 绯荤粺閰嶇疆淇℃伅
	DB                   *gorm.DB                 // 鏁版嵁搴撴帴鍙
	REDIS                *redis.Client            // Redis 缂撳瓨鎺ュ彛
	FILE_TYPE_MAP        sync.Map                 // 鏂囦欢绫诲瀷鏄犲皠
	ID_GENERATOR         *sonyflake.Sonyflake     // 涓婚敭鐢熸垚鍣
	CONTEXT              = context.Background()   // 涓婁笅鏂囦俊鎭
	AUTO_CREATE_DB       = true                   // 鏄惁鑷姩鐢熸垚鏁版嵁搴
	MAX_USERNAME_LENGTH  = 32                     // 鐢ㄦ埛鍚嶆渶澶ч暱搴
	MIN_PASSWORD_PATTERN = "^[_a-zA-Z0-9]{6,32}$" // 瀵嗙爜鏍煎紡
	START_TIME           = "2022-05-21 00:00:01"  // 鍥哄畾鍚姩鏃堕棿锛屼繚璇佺敓鎴  ID 鍞竴鎬
	FEED_NUM             = 30                     // 姣忔杩斿洖瑙嗛鏁伴噺
	VIDEO_ADDR           = "./public/video/"      // 瑙嗛瀛樻斁浣嶇疆
	COVER_ADDR           = "./public/cover/"      // 灏侀潰瀛樻斁浣嶇疆
	MAX_FILE_SIZE        = int64(10 << 20)        // 涓婁紶鏂囦欢澶у皬闄愬埗涓 10MB
	MAX_TITLE_LENGTH     = 140                    // 瑙嗛鎻忚堪鏈€澶ч暱搴
	MAX_COMMENT_LENGTH   = 300                    // 璇勮鏈€澶ч暱搴
	WHITELIST_VIDEO      = map[string]bool{".mp4": true, ".avi": true, ".wmv": true, ".mpeg": true,
		".mov": true, ".flv": true, ".rmvb": true, ".3gb": true, ".vob": true, ".m4v": true}
)

// 杩囨湡鏃堕棿
var (
	FAVORITE_EXPIRE       = 10 * time.Minute
	VIDEO_COMMENTS_EXPIRE = 10 * time.Minute
	COMMENT_EXPIRE        = 10 * time.Minute
	FOLLOW_EXPIRE         = 10 * time.Minute
	USER_INFO_EXPIRE      = 10 * time.Minute
	VIDEO_EXPIRE          = 10 * time.Minute
	PUBLISH_EXPIRE        = 10 * time.Minute
	EMPTY_EXPIRE          = 10 * time.Minute
	EXPIRE_TIME_JITTER    = 10 * time.Minute
)
