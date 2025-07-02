/**
 * @author Andrey Mikhalevich <katrenplus@mail.ru>, 2017
 
 * @class
 * @classdesc
	
 * @param {string} id view identifier
 * @param {namespace} options
 * @param {namespace} options.models All data models
 * @param {namespace} options.variantStorage {name,model}
 */	
function Login_View(id,options){	

	Login_View.superclass.constructor.call(this,id,options);
	
	var self = this;
	
	var err_ctrl = new ErrorControl(id+":error");
	this.addElement(err_ctrl);
	
	var check_for_enter = function(e){
		e = EventHelper.fixKeyEvent(e);
		if (e.keyCode==13){
			self.login();
		}
	};
					
	this.addElement(new EditString(id+":user",{				
		"html":"<input/>",
		"focus":true,
		"maxlength":"100",
		"cmdClear":false,
		"required":true,
		"errorControl":err_ctrl,//new ErrorControl(id+":user:error"),
		"events":{"keydown":check_for_enter}
	}));	
	
	this.addElement(new EditPassword(id+":pwd",{
		"html":"<input/>",
		"maxlength":"50",
		"required":true,
		"errorControl":err_ctrl,//new ErrorControl(id+":pwd:error"),
		"events":{"keydown":check_for_enter}
	}));	
/*
	this.addElement(new EditCheckBox(id+":rememberMe",{
		"html":"<input/>"
	}));	
*/
	this.addElement(new Button(id+":submit_login",{
		"onClick":function(){
			self.login();
		}
	}));
	
	//Commands
	var contr = new User_Controller();
	var pm = contr.getPublicMethod("login");
	pm.setFieldValue("width_type",window.getWidthType());
	
	this.addCommand(new Command("login",{
		"publicMethod":pm,
		"control":this.getElement("submit_login"),
		"async":false,
		"bindings":[
			new DataBinding({"field":pm.getField("name"),"control":this.getElement("user")})
			,new DataBinding({"field":pm.getField("pwd"),"control":this.getElement("pwd")})
			//new DataBinding({"field":pm.getField("rememberMe"),"control":this.getElement("rememberMe")})
		]		
	}));

}
extend(Login_View,ViewAjx);

Login_View.prototype.setError = function(s){
	this.getElement("error").setValue(s);
}

Login_View.prototype.login = function(){
	var self = this;
	this.execCommand("login",
		function(){
			var REDIR_PAR = "?redir=";
			var p = window.location.href.indexOf(REDIR_PAR);
			var redir;
			if(p>=0){
				redir = "?"+CommonHelper.unserialize(decodeURI(window.location.href.substr(p+REDIR_PAR.length))).ref;
			}
			else{
				redir = window.location.href;
			}
			document.location.href = redir;//window.getApp().getHost();
		},
		function(resp,errCode,errStr){
			if (errCode==1000){
				self.setError(self.ER_LOGIN);
			}
			else{
				self.setError((errCode!=undefined)? (errCode+" "+errStr):errStr);
			}
		}
	);
}
