<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:decimal-format name="num_money" decimal-separator="," grouping-separator=" "/>

<xsl:template name="format_num">
	<xsl:param name="val"/>
	<xsl:choose>
		<xsl:when test="$val='0' or string(number($val))='NaN'">
			<xsl:text>&#160;</xsl:text>
		</xsl:when>
		<xsl:otherwise>
			<xsl:value-of select="$val"/>
		</xsl:otherwise>		
	</xsl:choose>
</xsl:template>

<xsl:template name="format_money">
	<xsl:param name="val"/>
	<xsl:choose>
		<xsl:when test="$val='0' or string(number($val))='NaN'">
			<xsl:text>&#160;</xsl:text>
		</xsl:when>
		<xsl:otherwise>
			<!-- format-number($val,'### ###,00','num_money') -->
			<xsl:value-of select="format-number($val,'# ##0,00')"/>
		</xsl:otherwise>		
	</xsl:choose>
</xsl:template>

<xsl:template name="string-replace-all">
  <xsl:param name="text" />
  <xsl:param name="replace" />
  <xsl:param name="by" />
  <xsl:choose>
    <xsl:when test="contains($text, $replace)">
      <xsl:value-of select="substring-before($text,$replace)" />
      <xsl:value-of select="$by" />
      <xsl:call-template name="string-replace-all">
        <xsl:with-param name="text"
        select="substring-after($text,$replace)" />
        <xsl:with-param name="replace" select="$replace" />
        <xsl:with-param name="by" select="$by" />
      </xsl:call-template>
    </xsl:when>
    <xsl:otherwise>
      <xsl:value-of select="$text" />
    </xsl:otherwise>
  </xsl:choose>
</xsl:template>

<xsl:template name="format_quant">
	<xsl:param name="val"/>
	<xsl:choose>
		<xsl:when test="$val='0' or string(number($val))='NaN'">
			<xsl:text>&#160;</xsl:text>
		</xsl:when>
		<xsl:otherwise>
			<xsl:value-of select="format-number( round(10000*$val) div 10000 ,'##0.0000')"/>
		</xsl:otherwise>		
	</xsl:choose>
</xsl:template>

<xsl:template name="format_date">
	<xsl:param name="val"/>
	<xsl:param name="formatStr"/>
	<xsl:choose>
		<xsl:when test="$formatStr='d.m.y'">
			<xsl:variable name="val_year" select="substring-before($val,'-')"/>
			<xsl:variable name="part_month" select="substring-after($val,'-')"/>
			<xsl:variable name="val_month" select="substring-before($part_month,'-')"/>
			<xsl:variable name="part_date" select="substring-after($part_month,'-')"/>
			<xsl:variable name="val_date" select="substring-before($part_date,'T')"/>
			<xsl:value-of select="concat($val_date, '.', $val_month, '.', $val_year)" />
		</xsl:when>
		<xsl:otherwise>
			<xsl:value-of select="$val" />
		</xsl:otherwise>
	</xsl:choose>
</xsl:template>


</xsl:stylesheet>
