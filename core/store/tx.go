// @Contact:    huaxinrui
// @Time:       2019/8/8 上午10:39

package store

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	TX       *gorm.DB
	once     sync.Once
	rollback bool
	err      error
}

func (t *Transaction) Close(c *gin.Context) {
	t.once.Do(func() {
		if t.rollback {
			t.TX.Rollback()
			if c != nil {
				c.JSON(200, gin.H{"error_code": 1, "err_msg": t.err.Error(), "data": nil})
			}
		} else {
			t.TX.Commit()
			if c != nil {
				c.JSON(200, gin.H{"error_code": 0, "err_msg": "", "data": nil})
			}
		}
	})
}

func (t *Transaction) Fail(err error) {
	t.err = err
	t.rollback = true
}

func newTransaction(db *gorm.DB) *Transaction {
	return &Transaction{TX: db}
}
