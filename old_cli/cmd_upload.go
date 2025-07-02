package main

import (
	"fmt"
	"path/filepath"
)

// Upload uploads project to production servers.
func Upload(args []string) error {
	conf := ProjectConfig{}
	if err := conf.Read(); err != nil {
		return fmt.Errorf("conf.Read() failed: %v", err)
	}

	proj_dir, err := GetProjectDir()
	if err != nil {
		return fmt.Errorf("GetProjectDir() failed: %v", err)
	}

	//compile app for production
	fmt.Println(conf.Production.Compile)
	out_text, err := RunCMD(conf.Production.Compile, true)
	if err != nil {
		return err
	}
	if len(out_text) > 0 {
		fmt.Println(out_text)
	}

	for _, host := range conf.Production.Hosts {
		//stop app on host
		bash_cmd := fmt.Sprintf("ssh -p %d %s@%s %s",
			host.Port,
			host.User,
			host.Ip,
			host.AppStop,
		)
		fmt.Println(bash_cmd)
		out_text, err = RunCMD(bash_cmd, true)
		if err != nil {
			return err
		}
		if len(out_text) > 0 {
			fmt.Println(out_text)
		}

		//copy modified files
		//TODO: send only modified files
		for _, file := range conf.Production.Files {
			//"/" in the end and APP_DIR/dirname(file)
			bash_cmd := fmt.Sprintf("rsync -az -e \"ssh\" -p %d %s %s%s:%s",
				host.Port,
				filepath.Join(proj_dir, file),
				host.User,
				host.Ip,
				filepath.Join(host.AppDir, file),
			)
			fmt.Println(bash_cmd)
			out_text, err = RunCMD(bash_cmd, true)
			if err != nil {
				return err
			}
			if len(out_text) > 0 {
				fmt.Println(out_text)
			}
		}
	}

	//apply sql
	slq_dir := filepath.Join(proj_dir, UPDATES_DIR)
	if err := applySQL(proj_dir, slq_dir); err != nil {
		return err
	}

	//start remote app
	for _, host := range conf.Production.Hosts {
		bash_cmd := fmt.Sprintf("ssh -p %d %s@%s %s",
			host.Port,
			host.User,
			host.Ip,
			host.AppStart,
		)
		fmt.Println(bash_cmd)
		out_text, err = RunCMD(bash_cmd, true)
		if err != nil {
			return err
		}
		if len(out_text) > 0 {
			fmt.Println(out_text)
		}
	}
	return nil

}
