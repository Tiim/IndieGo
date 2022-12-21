package manualbackup

import (
	"fmt"
	"tiim/go-comment-api/config"
	"tiim/go-comment-api/model"
	"tiim/go-comment-api/plugins/admin"
)

type ManualBackupPlugin struct{}

func init() {
	config.RegisterModule(&ManualBackupPlugin{})
}

func (p *ManualBackupPlugin) IndieGoModule() config.ModuleInfo {
	return config.ModuleInfo{
		Name: "manual-backup",
		New:  func() config.Module { return new(ManualBackupPlugin) },
		Docs: config.ConfigDocs{
			DocString: `Manual backup module. This module enables the manual backup feature section in the admin dashboard. Must be loaded after the admin module and after a store module that implements model.BackupStore.`,
		},
	}
}

func (p *ManualBackupPlugin) Load(config config.GlobalConfig, _ interface{}) (config.ModuleInstance, error) {
	adminInt, err := config.GetModule("admin")
	if err != nil {
		return nil, fmt.Errorf("admin plugin not loaded, can not register admin section: %w", err)
	}
	admin, ok := adminInt.(*admin.AdminModule)
	if !ok {
		return nil, fmt.Errorf("admin is not a of type admin.AdminModule: %T", adminInt)
	}
	storeInt, err := config.GetModule("store")
	if err != nil {
		return nil, fmt.Errorf("store plugin not loaded, can not register admin section: %w", err)
	}
	store, ok := storeInt.(model.BackupStore)
	if !ok {
		return nil, fmt.Errorf("store is not a of type model.BackupStore: %T", storeInt)
	}
	admin.RegisterSection(newAdminBackupSection(store))

	return new(interface{}), nil
}
