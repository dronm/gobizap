/** Copyright (c) 2020
 *	Andrey Mikhalevich, Katren ltd.
 */
function LoginDeviceList_View(id,options){	
	options = options || {};
	
	this.m_onBanSession = options.onBanSession;
	
	LoginDeviceList_View.superclass.constructor.call(this,id,options);

	var model = (options.models&&options.models.LoginDeviceList_Model)? options.models.LoginDeviceList_Model:new LoginDeviceList_Model();
	var contr = new LoginDevice_Controller();
	
	var constants = {"doc_per_page_count":null,"grid_refresh_interval":null};
	window.getApp().getConstantManager().get(constants);
	
	var self = this;
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
			"cmdEdit":false,
			"exportFileName" :"УстройстваВхода"			
		}),
		"popUpMenu":popup_menu,
		"head":new GridHead(id+"-grid:head",{
			"elements":[
				new GridRow(id+":grid:head:row0",{
					"elements":[
						!options.detail? new GridCellHead(id+":grid:head:user_descr",{
							"value":"Пользователь",
							"columns":[
								new GridColumn({
									"field":model.getField("user_descr")
								})
							],
							"sortable":true
						}):null
					
						,new GridCellHead(id+":grid:head:date_time_in",{
							"value":"Последняя дата входа",
							"columns":[
								new GridColumnDate({
									"field":model.getField("date_time_in"),
									"dateFormat":"d/m/y H:i"
								})
							],
							"sortable":true,
							"sort":"desc"
						})
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
							"title":"Запрет входа с устройства",
							"columns":[
								new GridColumn({
									"field":model.getField("banned")
									,"cellOptions":function(column,row){
										var res = {};
										var m = self.getElement("grid").getModel();
										var ctrl = new EditCheckBox(null,{
											"value":m.getFieldValue("banned")
											,"events":{
												"change":(function(banHash,userId){
													return function(e){
														self.switchBan(
															banHash
															,userId
															,this
														);
													}
												})(m.getFieldValue("ban_hash")
													,m.getFieldValue("user_id")
												)
											}
										});
										res.elements = [ctrl];
										
										return res;
									}
								})
							]
						})
						
					]
				})
			]
		}),
		"pagination":null,				
		"autoRefresh":false,
		"refreshInterval":constants.grid_refresh_interval.getValue()*1000,
		"rowSelect":false,
		"focus":true
	}));	
	


}
extend(LoginDeviceList_View,ViewAjxList);

LoginDeviceList_View.prototype.switchBan = function(banHash,userId,ctrl){
	var banned = ctrl.getValue();
	
	var pm = (new LoginDevice_Controller()).getPublicMethod("switch_banned");
	pm.setFieldValue("banned",banned);
	pm.setFieldValue("user_id",userId);
	pm.setFieldValue("hash",banHash);
	ctrl.setEnabled(false);
	var self = this;
	pm.run({
		"all":function(){
			ctrl.setEnabled(true);
		}
		,"ok":function(resp){
			var msg;
			if(banned){
				msg = "Закрыт доступ для пользователя с устройства";
			}
			else{
				msg = "Отменен запрет доступа для пользователя с устройства";
			}
			window.showTempNote(msg,null,5000);
			
			if(banned && self.m_onBanSession){
				self.m_onBanSession();
			}
		}
	})
}



