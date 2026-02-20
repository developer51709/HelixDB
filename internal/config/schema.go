package config

type Config struct {
        Server   ServerConfig   `json:"server"`
        Storage  StorageConfig  `json:"storage"`
        Backup   BackupConfig   `json:"backup"`
        Recovery RecoveryConfig `json:"recovery"`
        Logging  LoggingConfig  `json:"logging"`
        Security SecurityConfig `json:"security"`
}

type ServerConfig struct {
        Port int    `json:"port"`
        Host string `json:"host"`
}

type StorageConfig struct {
        DataFile           string `json:"dataFile"`
        WALDirectory       string `json:"walDirectory"`
        AutoCompact        bool   `json:"autoCompact"`
        CompactThresholdMB int    `json:"compactThresholdMB"`
}

type BackupConfig struct {
        Enabled         bool   `json:"enabled"`
        Mode            string `json:"mode"`
        Directory       string `json:"directory"`
        IntervalMinutes int    `json:"intervalMinutes"`
}

type RecoveryConfig struct {
        AutoRecover     bool `json:"autoRecover"`
        VerifyChecksums bool `json:"verifyChecksums"`
}

type LoggingConfig struct {
        Level string `json:"level"`
        File  string `json:"file"`
}

type SecurityConfig struct {
        RequireAuth bool   `json:"requireAuth"`
        Token       string `json:"token"`
}

func DefaultConfig() Config {
        return Config{
                Server: ServerConfig{
                        Port: 5000,
                        Host: "0.0.0.0",
                },
                Storage: StorageConfig{
                        DataFile:           "./data/helix.db",
                        WALDirectory:       "./data/wal",
                        AutoCompact:        true,
                        CompactThresholdMB: 128,
                },
                Backup: BackupConfig{
                        Enabled:         false,
                        Mode:            "incremental",
                        Directory:       "./backups",
                        IntervalMinutes: 30,
                },
                Recovery: RecoveryConfig{
                        AutoRecover:     true,
                        VerifyChecksums: true,
                },
                Logging: LoggingConfig{
                        Level: "info",
                        File:  "",
                },
                Security: SecurityConfig{
                        RequireAuth: false,
                        Token:       "",
                },
        }
}
