/** Copyright (c) 2016,2019
 *	Andrey Mikhalevich, Katren ltd.
 */
 
function TimeZoneLocale_View(id,options){	

	TimeZoneLocale_View.superclass.constructor.call(this,id,options);
	
	var model = options.models.TimeZoneLocale_Model;
	var contr = new TimeZoneLocale_Controller();
	
	var is_employee = (window.getApp().getServVar("role_id")=="employee");
	var popup_menu = new PopUpMenu();
	
	this.addElement(new GridAjx(id+":grid",{
		"keyIds":["id"],
		"model":model,
		"controller":contr,
		"editInline":true,
		"editWinClass":null,
		"cmdInsert":false,
		"commands":new GridCmdContainerAjx(id+":grid:cmd",{
			"cmdEdit":!is_employee,
			"cmdDelete":!is_employee,
			"cmdInsert":false,
			"cmdSearch":false,
			"cmdExport":false
		}),		
		"popUpMenu":popup_menu,
		"head":new GridHead(id+"-grid:head",{
			"elements":[
				new GridRow(id+":grid:head:row0",{
					"elements":[
						new GridCellHead(id+":grid:head:name",{
							"value":"Наименование",
							"columns":[
								new GridColumn({"field":model.getField("name")})
							],
							"sortable":true,
							"sort":"asc"							
						}),
						new GridCellHead(id+":grid:head:descr",{
							"value":"Описание",
							"columns":[
								new GridColumn({"field":model.getField("descr")})
							],
							"sortable":true
						}),											
						new GridCellHead(id+":grid:head:hour_dif",{
							"value":"Разница GMT в часах",
							"colAttrs":{"align":"right"},
							"columns":[
								new GridColumn({"field":model.getField("hour_dif")})
							],
							"sortable":true
						})						
					]
				})
			]
		}),
		"autoRefresh":false,
		"focus":true,
		"rowSelect":true
	}));	
	


}
extend(TimeZoneLocale_View,ViewAjx);
