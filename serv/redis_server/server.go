package redis_server

import (
	"github.com/ikaven1024/redisAdmin/repository"
	"github.com/jinzhu/gorm"
)

type RedisMode int

const (
	RedisModeStandalone RedisMode = 0
	RedisModeCluster    RedisMode = 1
)

type RedisServer struct {
	repository.Model
	Name      string                 `gorm:"type:varchar(20);not null"`
	Mode      RedisMode              `gorm:"type:integer;not null"`
	Addresses repository.StringSlice `gorm:"type:varchar(500);not null"`
	Password  string                 `gorm:"type:varchar(20);not null"`
}

type Manager struct {
	repository *gorm.DB
}

func NewManager(repository *gorm.DB) *Manager {
	repository.AutoMigrate(&RedisServer{})

	return &Manager{
		repository: repository,
	}
}

func (m *Manager) Create(server *RedisServer) error {
	return m.repository.Create(server).Error
}

func (m *Manager) Update(server *RedisServer) error {
	// Save可以保存所有字段。Update只能保存非零字段。Password可能是空的，所以只能使用Save来保存
	return m.repository.Save(server).Error
}

func (m *Manager) Delete(id uint) error {
	server := RedisServer{}
	server.ID = id
	return m.repository.Delete(&server).Error
}

func (m *Manager) Get(id uint) (server RedisServer, err error) {
	err = m.repository.First(&server, id).Error
	return
}

func (m *Manager) List() (servers []RedisServer, err error) {
	err = m.repository.Find(&servers).Error
	return
}

func (s RedisServer) IsCluster() bool {
	return s.Mode == RedisModeCluster
}
