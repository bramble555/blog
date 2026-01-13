package snow

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Snowflake struct {
	sync.Mutex         // 锁
	timestamp    int64 // 时间戳 ，毫秒
	workerid     int64 // 工作节点
	datacenterid int64 // 数据中心机房id
	sequence     int64 // 序列号
}

const (
	epoch             = int64(1704067200000)                           // 设置起始时间(时间戳/毫秒)：2024-01-01 00:00:00，有效期69年
	timestampBits     = int64(41)                                      // 时间戳占用位数
	datacenteridBits  = int64(2)                                       // 数据中心id所占位数
	workeridBits      = int64(7)                                       // 机器id所占位数
	sequenceBits      = int64(12)                                      // 序列所占的位数
	timestampMax      = int64(-1 ^ (-1 << timestampBits))              // 时间戳最大值
	datacenteridMax   = int64(-1 ^ (-1 << datacenteridBits))           // 支持的最大数据中心id数量
	workeridMax       = int64(-1 ^ (-1 << workeridBits))               // 支持的最大机器id数量
	sequenceMask      = int64(-1 ^ (-1 << sequenceBits))               // 支持的最大序列id数量
	workeridShift     = sequenceBits                                   // 机器id左移位数
	datacenteridShift = sequenceBits + workeridBits                    // 数据中心id左移位数
	timestampShift    = sequenceBits + workeridBits + datacenteridBits // 时间戳左移位数
)

func Init() *Snowflake {
	// 在init函数中初始化snflk
	snflk := &Snowflake{
		timestamp:    epoch,
		workerid:     1, // 示例工作节点ID
		datacenterid: 1, // 示例数据中心ID
		sequence:     0, // 序列号通常从0开始
	}
	return snflk
}
func (s *Snowflake) GetID() int64 {
	s.Lock()
	now := time.Now().UnixNano() / 1000000 // 转毫秒
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		err := fmt.Errorf("epoch must be between 0 and %d", timestampMax-1)
		log.Println(err)
		return 0
	}
	s.timestamp = now
	id := int64((t)<<timestampShift | (s.datacenterid << datacenteridShift) | (s.workerid << workeridShift) | (s.sequence))
	s.Unlock()
	return id
}
