<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:output method="html" indent="yes"
			doctype-public="-//W3C//DTD XHTML 1.0 Strict//EN" 
			doctype-system="http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"/>

<xsl:variable name="BASE_PATH" select="/document/model[@id='ModelVars']/row[1]/basePath"/>
<xsl:variable name="VERSION" select="/document/model[@id='ModelVars']/row[1]/scriptId"/>
<xsl:variable name="TITLE" select="/document/model[@id='ModelVars']/row[1]/title"/>

<!--************* Main template ******************** -->		
<xsl:template match="/document">
<html>
	<head>
		<title>CRM</title>
		<xsl:apply-templates select="model[@id='ModelStyleSheet']"/>
	</head>
	<body>
		<div id="wrapper">
			
			<!-- Page Content -->
			<div id="page-wrapper">
			    <div class="container-fluid">
			    	<h1>Фатальная ошибка</h1>
			    	<h3>
			    	<xsl:value-of select="/document/model[@id='ModelServResponse']/row[1]/result"/>:
			    	<xsl:value-of select="/document/model[@id='ModelServResponse']/row[1]/descr"/>
			    	</h3>
			    	Необходима <a href="{/document/model[@id='ModelVars']/row[1]/basePath}">авторизация</a>.
			    </div>
			    <!-- /.container-fluid -->
			</div>
			<!-- /#page-wrapper -->
		</div>
		<!-- /#wrapper -->
	    
	</body>
</html>		
</xsl:template>

<!-- CSS -->
<xsl:template match="model[@id='ModelStyleSheet']/row">
	<link rel="stylesheet" href="{concat($BASE_PATH,href,'?',$VERSION)}" type="text/css"/>
</xsl:template>


</xsl:stylesheet>
