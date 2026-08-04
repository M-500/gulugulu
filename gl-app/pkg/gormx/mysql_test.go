package gormx

import "testing"

func TestValidateMySQLConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  MySQLConfig
		wantErr bool
	}{
		{name: "合法配置", config: MySQLConfig{DataSource: "user:pass@tcp(localhost:3306)/db", MaxIdleConns: 10, MaxOpenConns: 100, ConnMaxLife: 3600}},
		{name: "缺少DSN", config: MySQLConfig{MaxIdleConns: 10, MaxOpenConns: 100, ConnMaxLife: 3600}, wantErr: true},
		{name: "空闲连接超过最大连接", config: MySQLConfig{DataSource: "dsn", MaxIdleConns: 20, MaxOpenConns: 10, ConnMaxLife: 3600}, wantErr: true},
		{name: "连接生命周期不合法", config: MySQLConfig{DataSource: "dsn", MaxIdleConns: 1, MaxOpenConns: 10}, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateMySQLConfig(test.config)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateMySQLConfig() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
