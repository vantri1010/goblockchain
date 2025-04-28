package block

import (
	"goblockchain/utils"
	"time"
)

// SetNeighbors discovers and sets the list of neighbor nodes for the blockchain
func (bc *Blockchain) SetNeighbors() {
	// Call FindNeighbors to scan for nodes within the specified IP and port ranges
	// utils.GetHost() retrieves the current node's IP (e.g., "192.168.1.100")
	// bc.port is the current node's port (e.g., 5000)
	bc.neighbors = utils.FindNeighbors(
		utils.GetHost(), bc.port,
		NEIGHBOR_IP_RANGE_START, NEIGHBOR_IP_RANGE_END,
		BLOCKCHAIN_PORT_RANGE_START, BLOCKCHAIN_PORT_RANGE_END,
	)
	// Log the discovered neighbors for debugging (e.g., ["192.168.1.101:5000", "192.168.1.102:5001"])
	//log.Printf("%v", bc.neighbors)
}

// SyncNeighbors synchronizes the neighbor list with thread safety
func (bc *Blockchain) SyncNeighbors() {
	// Lock the mutex to prevent concurrent access to the neighbors list
	bc.muxNeighbors.Lock()
	// Ensure the mutex is unlocked after the function completes
	defer bc.muxNeighbors.Unlock()
	// Update the neighbors list by calling SetNeighbors
	bc.SetNeighbors()
}

// StartSyncNeighbors starts a periodic task to sync neighbors at regular intervals
func (bc *Blockchain) StartSyncNeighbors() {
	// Perform an immediate sync of neighbors
	bc.SyncNeighbors()
	// Schedule the next sync after BLOCKCHAIN_NEIGHBOR_SYNC_TIME_SEC (20 seconds)
	// time.AfterFunc runs StartSyncNeighbors again, creating a recursive loop
	_ = time.AfterFunc(time.Second*BLOCKCHAIN_NEIGHBOR_SYNC_TIME_SEC, bc.StartSyncNeighbors)
}
