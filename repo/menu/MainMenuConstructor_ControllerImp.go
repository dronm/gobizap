package menu

import (
	"reflect"	
	"fmt"	
	"context"
	
	"github.com/dronm/ds/pgds"
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"	
		
	"github.com/jackc/pgx/v5/pgxpool"
)

//Method implemenation delete
func (pm *MainMenuConstructor_Controller_delete) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["MainMenuConstructor"], sock.GetPresetFilter("MainMenuConstructor"))	
}

//Method implemenation get_object
func (pm *MainMenuConstructor_Controller_get_object) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["MainMenuConstructorDialog"], &MainMenuConstructorDialog{}, sock.GetPresetFilter("MainMenuConstructorDialog"))	
}

//Method implemenation get_list
func (pm *MainMenuConstructor_Controller_get_list) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["MainMenuConstructorList"], &MainMenuConstructorList{}, sock.GetPresetFilter("MainMenuConstructorList"))	
}

func (pm *MainMenuConstructor_Controller_insert) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetPrimary()
	if err != nil {
		return err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
	
	args := rfltArgs.Interface().(*MainMenuConstructor)
	new_cont, err := gen_user_menu(app, conn, args.Content.GetValue())
	if err != nil {
		return gobizap.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("MainMenuConstructor_Controller_insert gen_user_menu(): %v",err))
	}
	args.Model_content.SetValue(new_cont)

	return gobizap.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["MainMenuConstructor"], &MainMenuConstructor_keys{}, nil)
}

//Method implemenation
func (pm *MainMenuConstructor_Controller_update) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetPrimary()
	if err != nil {
		return err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
	
	args := rfltArgs.Interface().(*MainMenuConstructor_old_keys)
	if !args.Content.GetIsSet() {
		if err := conn.QueryRow(context.Background(),
			`SELECT content FROM main_menus WHERE id = $1`,
			args.Old_id.GetValue()).Scan(&args.Content); err != nil {
			return gobizap.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("MainMenuConstructor_Controller_update pgx.Conn.QueryRow(): %v",err))
		}		
	}
	new_cont, err := gen_user_menu(app, conn, args.Content.GetValue())
	if err != nil {
		return gobizap.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("MainMenuConstructor_Controller_update gen_user_menu(): %v",err))
	}
	args.Model_content.SetValue(new_cont)
	
	return gobizap.UpdateOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["MainMenuConstructor"], nil)
}

