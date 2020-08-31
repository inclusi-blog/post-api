package init

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
	"post-api/dbhelper"
)

func Db() *sqlx.DB {
	_ = logging.GetLogger(context.Background())
	_ = dbhelper.BuildConnectionString()
	return &sqlx.DB{}
}
