/**
 * @author Andrey Mikhalevich <katrenplus@mail.ru>, 2016
 
 * @class
 * @classdesc
	
 * @param {string} id view identifier
 * @param {namespace} options
 * @param {namespace} options.models All data models
 * @param {namespace} options.variantStorage {name,model}
 */	
function About_View(id,options){
	options = options || {};	
	
	options.tagName = "template";

	this.model = options.models.About_Model;
	this.model.getRow(0);
	
	options.templateOptions = {
		"LB_APP_NAME":this.LB_APP_NAME,
		"LB_VERSION":this.LB_VERSION,
		"LB_AUTHOR":this.LB_AUTHOR,
		"LB_TECH_MAIL":this.LB_TECH_MAIL,
		"LB_DB_NAME":this.LB_DB_NAME,
		"LB_FW_SERVER_VERSION":this.LB_FW_SERVER_VERSION,
		"LB_FW_CLIENT_VERSION":this.LB_FW_CLIENT_VERSION,
		"VERSION":window.getApp().VERSION
	};
	var fields = this.model.getFields();
	for(var fid in fields){
		options.templateOptions[fid] = fields[fid].getValue();
	}

	About_View.superclass.constructor.call(this,id,options);
}
extend(About_View,ViewAjx);

/* Constants */


/* private members */

/* protected*/


/* public methods */

