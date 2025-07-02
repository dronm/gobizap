<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:output method="html" indent="yes"
			doctype-public="-//W3C//DTD XHTML 1.0 Strict//EN" 
			doctype-system="http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"/>
			
<xsl:variable name="BASE_PATH" select="/document/model[@id='ServerVars']/row[1]/basePath"/>
<xsl:variable name="VERSION" select="/document/model[@id='ServerVars']/row[1]/scriptId"/>
<xsl:variable name="TITLE" select="/document/model[@id='ServerVars']/row[1]/title"/>
						
<xsl:template match="/">
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=UTF-8"/>
		<xsl:apply-templates select="/document/model[@id='ServerVars']/row"/>
		<xsl:apply-templates select="/document/model[@id='ModelStyleSheet']/row"/>		
		<title><xsl:value-of select="$TITLE"/></title>
	</head>
	<body>
	</body>
</html>		
</xsl:template>

<xsl:template match="model[@id='ServerVars']/row">
	<xsl:if test="author">
		<meta name="Author" content="{author}"></meta>
	</xsl:if>
	<xsl:if test="keywords">
		<meta name="Keywords" content="{keywords}"></meta>
	</xsl:if>
	<xsl:if test="description">
		<meta name="Description" content="{description}"></meta>
	</xsl:if>
	
</xsl:template>

<xsl:template match="model[@id='ModelStyleSheet']/row">
	<link rel="stylesheet" href="{concat($BASE_PATH,href,'?',$VERSION)}" type="text/css"/>
</xsl:template>

<xsl:template match="model[@id='ModelJavaScript']/row">
	<script src="{concat($BASE_PATH,href,'?',$VERSION)}"></script>
</xsl:template>

<xsl:template match="model[@id='ModelServResponse']/row">
	<xsl:if test="result/node()='1'">
	<div class="error"><xsl:value-of select="descr"/></div>
	</xsl:if>
</xsl:template>

</xsl:stylesheet>
