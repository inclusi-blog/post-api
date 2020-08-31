package init

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/dbhelper"
)

func Db() *sqlx.DB {
	_ = logging.GetLogger(context.Background())
	_ = dbhelper.BuildConnectionString()
	return &sqlx.DB{}
}
