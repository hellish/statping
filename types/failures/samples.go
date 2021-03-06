package failures

import (
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/statping/statping/types"
	"github.com/statping/statping/utils"
	"sync"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Samples() error {
	tx := db.Begin()
	sg := new(sync.WaitGroup)

	createdAt := utils.Now().Add(-3 * types.Day)

	for i := int64(1); i <= 4; i++ {
		sg.Add(1)

		f1 := &Failure{
			Service:   i,
			Issue:     "Server failure",
			CreatedAt: utils.Now().Add(-time.Duration(3*i) * 86400),
		}
		f1.Create()

		f2 := &Failure{
			Service:   i,
			Issue:     "Server failure",
			CreatedAt: utils.Now().Add(-time.Duration(5*i) * 12400),
		}
		f2.Create()

		log.Infoln(fmt.Sprintf("Adding %v Failure records to service", 400))

		go func() {
			defer sg.Done()
			for fi := 0.; fi <= float64(400); fi++ {
				createdAt = createdAt.Add(35 * time.Minute)
				failure := &Failure{
					Service:   i,
					Issue:     "testing right here",
					CreatedAt: createdAt.UTC(),
				}

				tx = tx.Create(&failure)
			}
		}()
	}
	sg.Wait()

	if err := tx.Commit().Error(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
