<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:import href="ViewBase.html.xsl"/>

<xsl:template match="/document">
<html>
	<head>
		<xsl:call-template name="initHead"/>		
		
		<script>		
			function pageLoad(){
			
				<xsl:call-template name="initApp"/>
				
				var view = new PasswordRecovery_View("PasswordRecovery");
				view.toDOM();
			}
		</script>		
		
	</head>
	<body onload="pageLoad();" class="login-container">

	<xsl:call-template name="page_header"/>
	
	<!-- Page container -->
	<div class="page-container">

		<!-- Page content -->
		<div class="page-content">

			<!-- Main content -->
			<div class="content-wrapper">

				<!-- Content area -->
				<div class="content">

					<!-- Password recovery -->
					<form id="PasswordRecovery">
						<div class="panel panel-body login-form">
							<div class="text-center">
								<div class="icon-object border-{$COLOR_PALETTE} text-{$COLOR_PALETTE}"><i class="icon-spinner11"></i></div>
								<h5 class="content-group">Восстановление пароля <small class="display-block">Мы направим Вам инструкцию по электронной почте</small></h5>
								<div id="PasswordRecovery:error"></div>
							</div>

							<div class="form-group has-feedback">
								<input id="PasswordRecovery:email" type="text" class="form-control" placeholder="Ваш адрес электронной почты"/>
								<div class="form-control-feedback">
									<i class="icon-mail5 text-muted"></i>
								</div>
								<div id="PasswordRecovery:email:error"></div>
							</div>

							<div id="PasswordRecovery:captcha"/>

							<div id="PasswordRecovery:submit" class="btn bg-{$COLOR_PALETTE} btn-block">Восстановить пароль <i class="icon-arrow-left13 position-right"></i></div>
						</div>
					</form>
					<!-- /password recovery -->


					<div id="PasswordRecovered" class="hidden">
						<div class="panel panel-body login-form">
							<div class="text-center">
								<div class="icon-object border-{$COLOR_PALETTE} text-{$COLOR_PALETTE}"><i class="icon-spinner11"></i></div>
								<h4 id="PasswordRecovery:error" class="help-block text-info"><i class="icon-alert position-left">
									</i>Новый пароль направлен на Вашу электронную почту.
								</h4>								
							</div>
						</div>
					</div>

					<!-- Footer -->
					<div class="footer text-muted text-center">
						2022. <a href="#"></a>
					</div>
					<!-- /footer -->

				</div>
				<!-- /content area -->

			</div>
			<!-- /main content -->

		</div>
		<!-- /page content -->

	</div>
	<!-- /page container -->
		
	<script>
			var n = document.getElementById("PasswordRecovery:email");
			if (n){
				if (document.activeElement &amp;&amp; document.activeElement.id!="PasswordRecovery:email"){
					n.focus();
				}
			}
	</script>
		
		<xsl:call-template name="initJS"/>
	</body>
</html>		
</xsl:template>

</xsl:stylesheet>
