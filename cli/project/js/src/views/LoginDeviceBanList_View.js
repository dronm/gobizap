/** Copyright (c) 2020
 *	Andrey Mikhalevich, Katren ltd.
 */
function LoginDeviceBanList_View(id,options){	
	options = options || {};
	LoginDeviceBanList_View.superclass.constructor.call(this,id,options);

	var model = (options.models&&options.models.LoginDeviceBanList)? options.models.LoginDeviceBanList_Model:new LoginDeviceBanList_Model();
	var contr = new LoginDeviceBan_Controller();
	
	var constants = {"doc_per_page_count":null,"grid_refresh_interval":null};
	window.getApp().getConstantManager().get(constants);
	
	var popup_menu = new PopUpMenu();
	var pagClass = window.getApp().getPaginationClass();
	this.addElement(new GridAjx(id+":grid",{
		"model":model,
		"keyIds":["user_id","user_agent"],
		"controller":contr,
		"editInline":null,
		"editWinClass":null,
		"commands":new GridCmdContainerAjx(id+":grid:cmd",{
			"cmdInsert":false,
			"cmdEdit":false
		}),
		"popUpMenu":popup_menu,
		"head":new GridHead(id+"-grid:head",{
			"elements":[
				new GridRow(id+":grid:head:row0",{
					"elements":[
						!options.detail? new GridCellHead(id+":grid:head:users_ref",{
							"value":"Пользователь",
							"columns":[
								new GridColumnRef({
									"field":model.getField("users_ref")
								})
							],
							"sortable":true
						}):null
					
						,new GridCellHead(id+":grid:head:user_agent",{
							"value":"Устройство",
							"columns":[
								new GridColumn({
									"field":model.getField("user_agent")
								})
							]
						})
						,new GridCellHead(id+":grid:head:banned",{
							"value":"Доступ закрыт",
							"columns":[
								new GridColumnBool({
									"field":model.getField("banned")
									,"showFalse":false
								})
							]
						})
						
					]
				})
			]
		}),
		"pagination":new pagClass(id+"_page",
			{"countPerPage":constants.doc_per_page_count.getValue()}),			
		"autoRefresh":false,
		"refreshInterval":constants.grid_refresh_interval.getValue()*1000,
		"rowSelect":false,
		"focus":true
	}));	
	


}
extend(LoginDeviceBanList_View,ViewAjxList);

