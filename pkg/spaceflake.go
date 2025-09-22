package pkg

import (
	"github.com/kkrypt0nn/spaceflake"
)

func GenerateID(nodeID uint64, workerID uint64, seq uint64) (int64, error) {
	settings := spaceflake.NewGeneratorSettings()
	settings.BaseEpoch = 1640995200000 // Saturday, January 1, 2022 12:00:00 AM GMT
	settings.NodeID = nodeID
	settings.WorkerID = workerID
	settings.Sequence = seq
	sf, err := spaceflake.Generate(settings)
	if err != nil {
		return 0, err

	}

	return int64(sf.ID()), nil
}

func GenerateIDs(nodeID uint64, amount int) ([]int64, error) {
	settings := spaceflake.NewGeneratorSettings()
	settings.BaseEpoch = 1640995200000 // Saturday, January 1, 2022 12:00:00 AM GMT
	node := spaceflake.NewNode(nodeID)
	sf, err := node.BulkGenerateSpaceflakes(amount)

	if err != nil {
		return []int64{}, err

	}
	id := make([]int64, amount)
	for i, obj := range sf {
		id[i] = int64(obj.ID())
	}

	return id, nil
}

func Decompose(id int64) map[string]uint64 {
	return spaceflake.Decompose(uint64(id), 1640995200000)
}
