package sync

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/agentwiki/agentwiki/internal/network/protocol"
	"github.com/agentwiki/agentwiki/internal/storage"
	"github.com/agentwiki/agentwiki/internal/storage/model"
)

type SyncState string

const (
	SyncStateIdle      SyncState = "idle"
	SyncStateSyncing   SyncState = "syncing"
	SyncStateError     SyncState = "error"
	SyncStateComplete  SyncState = "complete"
)

type SyncConfig struct {
	AutoSync         bool
	IntervalSeconds  int
	MirrorCategories []string
	MaxLocalSizeMB   int
	BatchSize        int
}

type VersionVector map[string]int64

func (vv VersionVector) Increment(entryID string) int64 {
	vv[entryID]++
	return vv[entryID]
}

func (vv VersionVector) Get(entryID string) int64 {
	if v, ok := vv[entryID]; ok {
		return v
	}
	return 0
}

func (vv VersionVector) Merge(other VersionVector) VersionVector {
	result := make(VersionVector)
	for k, v := range vv {
		result[k] = v
	}
	for k, v := range other {
		if v > result[k] {
			result[k] = v
		}
	}
	return result
}

func (vv VersionVector) Diff(other VersionVector) []string {
	var needed []string
	for id, theirV := range other {
		if theirV > vv.Get(id) {
			needed = append(needed, id)
		}
	}
	return needed
}

func (vv VersionVector) ToProto() map[string]int64 {
	result := make(map[string]int64)
	for k, v := range vv {
		result[k] = v
	}
	return result
}

func VersionVectorFromProto(m map[string]int64) VersionVector {
	vv := make(VersionVector)
	for k, v := range m {
		vv[k] = v
	}
	return vv
}

type SyncEngine struct {
	store       *storage.Store
	config      *SyncConfig
	state       SyncState
	versionVec  VersionVector
	lastSync    int64
	mu          sync.RWMutex
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

func NewSyncEngine(store *storage.Store, cfg *SyncConfig) *SyncEngine {
	return &SyncEngine{
		store:      store,
		config:     cfg,
		state:      SyncStateIdle,
		versionVec: make(VersionVector),
	}
}

func (se *SyncEngine) Start(ctx context.Context) error {
	se.mu.Lock()
	defer se.mu.Unlock()

	if !se.config.AutoSync {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	se.cancel = cancel

	se.wg.Add(1)
	go se.syncLoop(ctx)

	return nil
}

func (se *SyncEngine) Stop() error {
	se.mu.Lock()
	defer se.mu.Unlock()

	if se.cancel != nil {
		se.cancel()
	}
	se.wg.Wait()
	return nil
}

func (se *SyncEngine) syncLoop(ctx context.Context) {
	defer se.wg.Done()

	ticker := time.NewTicker(time.Duration(se.config.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			se.IncrementalSync(ctx)
		}
	}
}

func (se *SyncEngine) IncrementalSync(ctx context.Context) error {
	se.mu.Lock()
	if se.state == SyncStateSyncing {
		se.mu.Unlock()
		return nil
	}
	se.state = SyncStateSyncing
	se.mu.Unlock()

	defer func() {
		se.mu.Lock()
		se.state = SyncStateIdle
		se.lastSync = time.Now().UnixMilli()
		se.mu.Unlock()
	}()

	// 这里应该实现与其他种子节点的同步逻辑
	// 1. 发现其他种子节点
	// 2. 与每个种子节点交换版本向量
	// 3. 拉取需要更新的数据
	// 4. 合并数据并解决冲突

	return nil
}

// SeedSync 执行种子节点间的同步
func (se *SyncEngine) SeedSync(ctx context.Context, peerIDs []string) error {
	// 实现种子节点间的同步逻辑
	// 1. 与每个种子节点建立连接
	// 2. 交换版本向量
	// 3. 拉取差异数据
	// 4. 合并数据并解决冲突
	// 5. 推送本地更新

	return nil
}

// ResolveConflict 解决数据冲突
func (se *SyncEngine) ResolveConflict(local, remote *model.KnowledgeEntry) *model.KnowledgeEntry {
	// 冲突解决策略：最后写入胜出
	if remote.UpdatedAt > local.UpdatedAt {
		return remote
	}
	return local
}

// GetSyncStats 获取同步统计信息
func (se *SyncEngine) GetSyncStats() map[string]interface{} {
	se.mu.RLock()
	defer se.mu.RUnlock()

	return map[string]interface{}{
		"state":     string(se.state),
		"last_sync": se.lastSync,
		"entry_count": len(se.versionVec),
	}
}

func (se *SyncEngine) GetState() SyncState {
	se.mu.RLock()
	defer se.mu.RUnlock()
	return se.state
}

func (se *SyncEngine) GetVersionVector() VersionVector {
	se.mu.RLock()
	defer se.mu.RUnlock()

	vv := make(VersionVector)
	for k, v := range se.versionVec {
		vv[k] = v
	}
	return vv
}

func (se *SyncEngine) MergeEntries(ctx context.Context, entries []*model.KnowledgeEntry) error {
	for _, entry := range entries {
		localVersion := se.versionVec.Get(entry.ID)
		if entry.Version > localVersion {
			if _, err := se.store.Entry.Create(ctx, entry); err != nil {
				continue
			}
			se.versionVec[entry.ID] = entry.Version
		}
	}
	return nil
}

func (se *SyncEngine) HandleSyncRequest(ctx context.Context, req *protocol.SyncRequest) (*protocol.SyncResponse, error) {
	se.mu.RLock()
	defer se.mu.RUnlock()

	clientVV := VersionVectorFromProto(req.VersionVector)

	resp := &protocol.SyncResponse{
		RequestID:           req.RequestID,
		NewEntries:          [][]byte{},
		UpdatedEntries:      [][]byte{},
		DeletedEntryIDs:     []string{},
		NewRatings:          [][]byte{},
		ServerVersionVector: se.versionVec.ToProto(),
		ServerTimestamp:     time.Now().UnixMilli(),
	}

	entries, _, err := se.store.Entry.List(ctx, storage.EntryFilter{Limit: 1000})
	if err != nil {
		return nil, fmt.Errorf("list entries: %w", err)
	}

	for _, entry := range entries {
		if entry.UpdatedAt > req.LastSyncTimestamp {
			if clientVV.Get(entry.ID) < entry.Version {
				data, err := entry.ToJSON()
				if err != nil {
					continue
				}
				if entry.Version == 1 {
					resp.NewEntries = append(resp.NewEntries, data)
				} else {
					resp.UpdatedEntries = append(resp.UpdatedEntries, data)
				}
			}
		}
	}

	return resp, nil
}

func (se *SyncEngine) HandleMirrorRequest(ctx context.Context, req *protocol.MirrorRequest) (<-chan *protocol.MirrorData, error) {
	dataCh := make(chan *protocol.MirrorData, 10)

	go func() {
		defer close(dataCh)

		// 构建分类过滤条件
		filter := storage.EntryFilter{
			Limit: 10000,
		}

		// 如果请求中指定了分类，则只返回对应分类的条目
		if len(req.Categories) > 0 {
			filter.Categories = req.Categories
		} else if len(se.config.MirrorCategories) > 0 {
			// 如果配置中指定了镜像分类，则使用配置的分类
			filter.Categories = se.config.MirrorCategories
		}

		entries, _, err := se.store.Entry.List(ctx, filter)
		if err != nil {
			return
		}

		batchSize := int(req.BatchSize)
		if batchSize <= 0 {
			batchSize = se.config.BatchSize
		}
		if batchSize <= 0 {
			batchSize = 100
		}

		totalBatches := (len(entries) + batchSize - 1) / batchSize

		for i := 0; i < totalBatches; i++ {
			start := i * batchSize
			end := start + batchSize
			if end > len(entries) {
				end = len(entries)
			}

			batch := entries[start:end]
			var entryData [][]byte
			for _, e := range batch {
				data, err := e.ToJSON()
				if err != nil {
					continue
				}
				entryData = append(entryData, data)
			}

			dataCh <- &protocol.MirrorData{
				RequestID:    req.RequestID,
				BatchIndex:   int32(i),
				TotalBatches: int32(totalBatches),
				Entries:      entryData,
			}
		}
	}()

	return dataCh, nil
}

// StartMirrorSync 开始镜像同步
func (se *SyncEngine) StartMirrorSync(ctx context.Context, categories []string) error {
	// 实现镜像同步逻辑
	// 1. 连接种子节点
	// 2. 发送镜像请求
	// 3. 接收并处理镜像数据
	// 4. 存储到本地

	return nil
}

// GetMirrorStatus 获取镜像状态
func (se *SyncEngine) GetMirrorStatus() map[string]interface{} {
	se.mu.RLock()
	defer se.mu.RUnlock()

	return map[string]interface{}{
		"mirror_categories": se.config.MirrorCategories,
		"max_local_size_mb": se.config.MaxLocalSizeMB,
		"last_sync":         se.lastSync,
	}
}
