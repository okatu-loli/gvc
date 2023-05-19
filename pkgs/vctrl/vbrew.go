package vctrl

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/query"
	"github.com/moqsien/gvc/pkgs/utils"
	"github.com/moqsien/gvc/pkgs/utils/tui"
	"github.com/pterm/pterm"
)

type Homebrew struct {
	Conf    *config.GVConfig
	envs    *utils.EnvsHandler
	fetcher *query.Fetcher
}

func NewHomebrew() (hb *Homebrew) {
	hb = &Homebrew{
		Conf:    config.New(),
		fetcher: query.NewFetcher(),
		envs:    utils.NewEnvsHandler(),
	}
	hb.envs.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *Homebrew) getShellScript() string {
	fPath := filepath.Join(config.HomebrewFileDir, "homebrew.sh")
	if ok, _ := utils.PathIsExist(fPath); !ok {
		that.fetcher.Url = that.Conf.Homebrew.ShellScriptUrl
		if size := that.fetcher.GetAndSaveFile(fPath); size > 0 {
			return fPath
		}
		return ""
	}
	return fPath
}

func (that *Homebrew) SetEnv() {
	mirror := ""
	tui.PrintInfo("Choose a Mirror Site in China: ")
	pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
		{Level: 0, Text: "TsingHua.", TextStyle: pterm.NewStyle(pterm.FgCyan), Bullet: "1)", BulletStyle: pterm.NewStyle(pterm.FgYellow)},
		{Level: 0, Text: "USTC.", TextStyle: pterm.NewStyle(pterm.FgCyan), Bullet: "2)", BulletStyle: pterm.NewStyle(pterm.FgYellow)},
	}).Render()
	fmt.Print(pterm.Cyan("Input>> "))
	fmt.Scan(&mirror)
	mirror = strings.TrimSpace(mirror)
	switch mirror {
	case "1":
		envMap := that.Conf.Homebrew.TsingHua
		envars := fmt.Sprintf(utils.HOMEbrewEnv,
			envMap["HOMEBREW_API_DOMAIN"],
			envMap["HOMEBREW_BOTTLE_DOMAIN"],
			envMap["HOMEBREW_BREW_GIT_REMOTE"],
			envMap["HOMEBREW_CORE_GIT_REMOTE"],
			envMap["HOMEBREW_PIP_INDEX_URL"])
		that.envs.UpdateSub(utils.SUB_BREW, envars)
	case "2":
		envMap := that.Conf.Homebrew.USTC
		envars := fmt.Sprintf(utils.HOMEbrewEnv,
			envMap["HOMEBREW_API_DOMAIN"],
			envMap["HOMEBREW_BOTTLE_DOMAIN"],
			envMap["HOMEBREW_BREW_GIT_REMOTE"],
			envMap["HOMEBREW_CORE_GIT_REMOTE"],
			envMap["HOMEBREW_PIP_INDEX_URL"])
		that.envs.UpdateSub(utils.SUB_BREW, envars)
	default:
		tui.PrintError("Unknown Mirror Choice.")
	}
}

func (that *Homebrew) Install() {
	if runtime.GOOS != utils.Windows {
		script := that.getShellScript()
		if _, err := utils.ExecuteSysCommand(false, "sh", script); err != nil {
			tui.PrintError(err)
			return
		}
		that.SetEnv()
	} else {
		tui.PrintError("Homebrew does not support Windows.")
	}
}
