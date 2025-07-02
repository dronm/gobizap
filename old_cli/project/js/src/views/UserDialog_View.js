/* Copyright (c) 2016
 *	Andrey Mikhalevich, Katren ltd.
 */
function UserDialog_View(id,options){	

	options = options || {};
	options.HEAD_TITLE = "Карточка пользователя";
	options.cmdSave = true;
	options.controller = new User_Controller();
	options.model = (options.models&&options.models.UserDialog_Model)? options.models.UserDialog_Model : new UserDialog_Model();

	var self = this;
	options.addElement = function(){
		this.addElement(new EditString(id+":name", {
			"labelCaption": "Логин пользователя",
			"length":"100",
			"required":true
		}));

		this.addElement(new EditString(id+":name_full", {
			"labelCaption": "ФИО пользователя",
			"length":"500",
			"required":true
		}));

		this.addElement(new Enum_role_types(id+":role",{
			"labelCaption":"Роль:",
			"required":true
		}));	

		this.addElement(new EditCheckBox(id+":banned",{
			"labelCaption":"Доступ запрещен:",
			"title":"При установленной галочке пользователь не сможет заходить в программу"
		}));

		this.addElement(new UserPhotoEdit(id+":photo",{
			"view": this
		}));		

		this.addElement(new ButtonCmd(id+":cmdResetPwd",{
			"onClick":function(){
				self.resetPwd();
			}
		}));								
		//mac grid
		//this.addElement(new UserMacAddressList_View(id+":mac_list"));

		//Login grid
		this.addElement(new LoginList_View(id+":login_list",{"detail":true}));

		//Login device grid
		this.addElement(new LoginDeviceList_View(id+":login_device_list",{
			"detail":true
			,"onBanSession":function(){
				//refresh session list
				self.getElement("login_list").getElement("grid").onRefresh();
			}
		}));

		this.addElement(new EntityContactList_View(id+":contacts_list",{
			"detail":true
		}));		

		var ctrl = new SalaryCalcEdit(id+":salary_calcs_ref",{
			"buttonSelect": new ButtonEditPeriodValues(id+":btn_perVals",{
				"editControl":ctrl,
				"periodOptions":{
					"formCaption": "История знчений",
					"valueType":"salary_calc",
					"gridClass":GridColumnRef,
					"gridOptions":{
						"ctrlClass":SalaryCalcEdit,
						"ctrlOptions":{
							"labelCaption":""
						}
					},
					
					"key":function(){
						return self.getElement("id").getValue();
					}			
				}
				
			})		
		});	
		this.addElement(ctrl);

	}	
	
	UserDialog_View.superclass.constructor.call(this,id,options);
	
	//****************************************************	
	
	//read
	var r_bd = [
		new DataBinding({"control":this.getElement("name")})
		,new DataBinding({"control":this.getElement("name_full")})
		,new DataBinding({"control":this.getElement("role"), "field":this.m_model.getField("role_id")})
		,new DataBinding({"control":this.getElement("photo")})
		,new DataBinding({"control":this.getElement("salary_calcs_ref")})
	];
	//if(options.templateOptions.BAN_ALLOWED){
		r_bd.push(new DataBinding({"control":this.getElement("banned")}));
	//}
	this.setDataBindings(r_bd);
	
	//write
	var wr_b = [
		new CommandBinding({"control":this.getElement("name")})
		,new CommandBinding({"control":this.getElement("name_full")})
		,new CommandBinding({"control":this.getElement("role"),"fieldId":"role_id"})
		,new CommandBinding({"control":this.getElement("photo")})
	];
	//if(options.templateOptions.BAN_ALLOWED){
		wr_b.push(new CommandBinding({"control":this.getElement("banned")}));		
	//}
	this.setWriteBindings(wr_b);
	
	this.addDetailDataSet({
		"control":this.getElement("login_list").getElement("grid"),
		"controlFieldId":"user_id",
		"field":this.m_model.getField("id")
	});

	this.addDetailDataSet({
		"control":this.getElement("login_device_list").getElement("grid"),
		"controlFieldId":"user_id",
		"field":this.m_model.getField("id")
	});
	
	this.addDetailDataSet({
		"control":this.getElement("contacts_list").getElement("grid"),
		"controlFieldId": ["entity_type", "entity_id"],
		"value": ["users", function(){
			return self.m_model.getFieldValue("id");
		}]
	});		
	
}
extend(UserDialog_View,ViewObjectAjx);

UserDialog_View.prototype.resetPwdCont = function(){
	var pm = this.getController().getPublicMethod("reset_pwd");
	pm.setFieldValue("user_id",this.getElement("id").getValue());
	var self = this;
	this.getElement("cmdResetPwd").setEnabled(false);		
	pm.run({
		"ok":function(resp){
			window.showTempNote("Пароль сброшен. Пользователю отправлено письмо с новым паролем.",null,5000);
			//self.close({"updated":true});
		}
		,"all":function(){
			self.getElement("cmdResetPwd").setEnabled(true);		
		}
	});
}

UserDialog_View.prototype.saveAndCall = function(callBack){
	if (this.getCmd() != "insert" && !this.getModified()){
		callBack();
	}
	else{	
		var self = this;
		this.getControlOK().setEnabled(false);		
		this.onSave(
			function(){
				callBack();
			},
			null,
			function(){
				self.getControlOK().setEnabled(true);		
			}
		);			
	}
}

UserDialog_View.prototype.resetPwd = function(){
	var self = this;
	this.saveAndCall(function(){
		self.resetPwdCont();
	});
	/*
	if (this.getCmd() != "insert" && !this.getModified()){
		this.resetPwdCont();
	}
	else{	
		var self = this;
		this.getControlOK().setEnabled(false);		
		this.onSave(
			function(){
				self.resetPwdCont();
			},
			null,
			function(){
				self.getControlOK().setEnabled(true);		
			}
		);			
	}
	*/
}

