/** Copyright (c) 2019
 * 	Andrey Mikhalevich, Katren ltd.
 */
function LoginList_View(id,options){	

	options = options || {};

	if(!options.detail){
		options.HEAD_TITLE = "Сеансы";
	}

	LoginList_View.superclass.constructor.call(this,id,options);
	
	var model = (options.models&&options.models.LoginList_Model)? options.models.LoginList_Model:new LoginList_Model();
	var contr = new Login_Controller();
	
	var constants = {"doc_per_page_count":null,"grid_refresh_interval":null};
	window.getApp().getConstantManager().get(constants);
	
	var popup_menu = new PopUpMenu();
	var pagClass = window.getApp().getPaginationClass();
	this.addElement(new GridAjx(id+":grid",{
		"model":model,
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
						,new GridCellHead(id+":grid:head:users_ref",{
							"value":"Пользоватль",
							"columns":[
								new GridColumnRef({
									"field":model.getField("users_ref"),
									"form":User_Form,
									"ctrlClass":UserEditRef,
									"searchOptions":{
										"field":new FieldInt("user_id"),
										"searchType":"on_match",
										"typeChange":false
									}									
								})
							],
							"sortable":true
						})
					
						,new GridCellHead(id+":grid:head:date_time_in",{
							"value":"Время входа",
							"columns":[
								new GridColumnDate({
									"field":model.getField("date_time_in"),
									"dateFormat":"d/m/y H:i"
								})
							],
							"sortable":true,
							"sort":"desc"							
						})
						,new GridCellHead(id+":grid:head:date_time_out",{
							"value":"Время выхода",
							"columns":[
								new GridColumnDate({
									"field":model.getField("date_time_out"),
									"dateFormat":"d/m/y H:i"
								})
							],
							"sortable":true
						})
						,new GridCellHead(id+":grid:head:set_date_time",{
							"value":"Время обновления",
							"columns":[
								new GridColumnDate({
									"field":model.getField("set_date_time"),
									"dateFormat":"d/m/y H:i"
								})
							],
							"sortable":true
						})
						,new GridCellHead(id+":grid:head:ip",{
							"value":"IP",
							"columns":[
								new GridColumn({
									"field":model.getField("ip")
								})
							]
						})
						,new GridCellHead(id+":grid:head:user_agent_descr",{
							"value":"Устройство",
							"columns":[
								new GridColumn({
									"field":model.getField("user_agent_descr")
								})
							]
						})
						
						/*,new GridCellHead(id+":grid:head:user_agent",{
							"value":"User-Agent",
							"columns":[
								new GridColumn({
									"field":model.getField("user_agent")
								})
							]
						})
						
						,new GridCellHead(id+":grid:head:headers",{
							"value":"Заголовки",
							"columns":[
								new GridColumn({
									"field":model.getField("headers"),
									"formatFunction":function(f){
										return CommonHelper.serialize(f.getFields());
									}
								})
							]
						})
						*/
					]
				})
			]
		}),
		"pagination":new pagClass(id+"_page",
			{"countPerPage":!options.detail? constants.doc_per_page_count.getValue():this.REC_COUNT_IN_DLG_MODE}),		
		
		"autoRefresh":false,
		"refreshInterval":constants.grid_refresh_interval.getValue()*1000,
		"rowSelect":false,
		"focus":true
	}));	
	


}
extend(LoginList_View,ViewAjxList);

LoginList_View.prototype.REC_COUNT_IN_DLG_MODE = 10;



