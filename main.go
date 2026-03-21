package main

import (
	"fmt"
)

func main() {
	opts := ParseArgs()

	// === config 写入模式 ===
	if opts.ConfigMode {
		if opts.From == "" || opts.Password == "" || opts.To == "" {
			fmt.Println("Usage: --config --from xxx --password xxx --to xxx")
			return
		}

		cfg := &Config{
			From:     opts.From,
			Password: opts.Password,
			To:       opts.To,
		}

		if err := SaveConfig(cfg); err != nil {
			fmt.Println("Save config failed:", err)
			return
		}

		fmt.Println("Config saved to", GetConfigPath())
		return
	}

	if len(opts.Command) == 0 {
		fmt.Println("Usage: mydaemon <command> [args...] / --config --from xxx --password xxx --to xxx")
		return
	}

	// 读取配置
	cfg, _ := LoadConfig()
	if cfg == nil {
		cfg = &Config{}
	}

	// CLI 覆盖
	if opts.From != "" {
		cfg.From = opts.From
	}
	if opts.Password != "" {
		cfg.Password = opts.Password
	}
	if opts.To != "" {
		cfg.To = opts.To
	}

	logFile := CreateLogFile()
	defer logFile.Close()

	duration, exitCode, tail := Run(opts.Command, logFile)
	result := FormatResult(opts.Command, duration, exitCode, tail)

	fmt.Println(result)

	if cfg.From != "" && cfg.Password != "" && cfg.To != "" {
		err := SendEmail(result, cfg.From, cfg.Password, cfg.To)
		if err != nil {
			fmt.Println("Send email failed:", err)
		} else {
			fmt.Println("Email sent successfully")
		}
	}
}
