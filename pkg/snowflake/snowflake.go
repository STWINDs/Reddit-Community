package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func InitSnowflake(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	sf.Epoch = st.UnixNano() / 1e6
	node, err = sf.NewNode(machineID)
	if err != nil {
		return err
	}
	return nil
}

func GenID() int64 {
	return node.Generate().Int64()
}
